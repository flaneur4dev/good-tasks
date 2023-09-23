package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" //nolint:revive

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	es "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/mistakes"
)

type DBStorage struct {
	db *sql.DB
}

func New(dsn string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	ds := &DBStorage{
		db: db,
	}

	if err = ds.initDB(); err != nil {
		return nil, err
	}

	return ds, nil
}

func (ds *DBStorage) Events(ctx context.Context, start, end time.Time) ([]cs.Event, error) {
	es := []cs.Event{}

	selq := `SELECT id, title, description, owner_id, start_date, finish_date, notification_time
		FROM events
		WHERE start_date BETWEEN $1 AND $2;
	`
	rows, err := ds.db.QueryContext(ctx, selq, start.Format(time.DateTime), end.Format(time.DateTime))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, title, desc, oID string
			st, ft, nt           sql.NullTime
		)

		if err := rows.Scan(&id, &title, &desc, &oID, &st, &ft, &nt); err != nil {
			return nil, err
		}

		es = append(es, cs.Event{
			ID:               id,
			Title:            title,
			Description:      desc,
			OwnerID:          oID,
			StartDate:        st.Time,
			FinishDate:       ft.Time,
			NotificationTime: nt.Time,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return es, nil
}

func (ds *DBStorage) CreateEvent(ctx context.Context, ne cs.Event) error {
	var id string
	selq := "SELECT id FROM events WHERE start_date<=$1 AND finish_date>=$1 LIMIT 1"

	err := ds.db.QueryRowContext(ctx, selq, ne.StartDate.Format(time.DateTime)).Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		break
	case id != "":
		return es.ErrDateBusy
	case err != nil:
		return err
	}

	insq := `INSERT INTO events (id, title, description, owner_id, start_date, finish_date, notification_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	res, err := ds.db.ExecContext(ctx, insq,
		ne.ID,
		ne.Title,
		ne.Description,
		ne.OwnerID,
		ne.StartDate.Format(time.DateTime),
		ne.FinishDate.Format(time.DateTime),
		ne.NotificationTime.Format(time.DateTime),
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return es.ErrCreateEvent
	}

	return nil
}

func (ds *DBStorage) UpdateEvent(ctx context.Context, id string, ne cs.Event) error {
	updq := `UPDATE events
		SET title=$1, description=$2, owner_id=$3, start_date=$4, finish_date=$5, notification_time=$6
		WHERE id=$7;
	`
	res, err := ds.db.ExecContext(ctx, updq,
		ne.Title,
		ne.Description,
		ne.OwnerID,
		ne.StartDate.Format(time.DateTime),
		ne.FinishDate.Format(time.DateTime),
		ne.NotificationTime.Format(time.DateTime),
		id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return es.ErrUpdateEvent
	}

	return nil
}

func (ds *DBStorage) DeleteEvent(ctx context.Context, id string) error {
	delq := "DELETE FROM events WHERE id=$1"

	res, err := ds.db.ExecContext(ctx, delq, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return es.ErrDeleteEvent
	}
	return nil
}

func (ds *DBStorage) Check(ctx context.Context) error {
	return ds.db.PingContext(ctx)
}

func (ds *DBStorage) Close() error {
	return ds.db.Close()
}

func (ds *DBStorage) initDB() error {
	q := `CREATE TABLE IF NOT EXISTS events (
		id varchar(50) PRIMARY KEY,
		title varchar(20) NOT NULL,
		description text NULL,
		owner_id varchar(50) NOT NULL,
		start_date timestamp,
		finish_date timestamp,
		notification_time timestamp
	);`

	_, err := ds.db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}
