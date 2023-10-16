package schedulerdb

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" //nolint:revive

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
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

	return &DBStorage{db}, nil
}

func (ds *DBStorage) Notifications(ctx context.Context, nt time.Time) ([]cs.Notification, error) {
	selq := "SELECT id, title, owner_id, start_date FROM events WHERE notification_day=$1"

	rows, err := ds.db.QueryContext(ctx, selq, nt.Format(time.DateTime))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ns := make([]cs.Notification, 0)
	for rows.Next() {
		var (
			id, title, oID string
			st             sql.NullTime
		)

		if err := rows.Scan(&id, &title, &oID, &st); err != nil {
			return nil, err
		}

		ns = append(ns, cs.Notification{
			UserID:     oID,
			EventID:    id,
			EventTitle: title,
			EventDate:  st.Time,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ns, nil
}

func (ds *DBStorage) Clear(ctx context.Context, t time.Time) error {
	delq := "DELETE FROM events WHERE start_date<=$1"

	res, err := ds.db.ExecContext(ctx, delq, t.Format(time.DateTime))
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (ds *DBStorage) Close() error {
	return ds.db.Close()
}
