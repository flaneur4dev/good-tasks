package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/server/grpc/pb"
	hs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/server/http/handlers"
)

type CalendarSuite struct {
	suite.Suite

	db             *sql.DB
	httpURL        string
	httpClient     *http.Client
	conn           *grpc.ClientConn
	grpcClient     pb.CalendarClient
	requestTimeout time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
}

func (s *CalendarSuite) SetupSuite() {
	s.setupDelay()
	s.setupDB()
	s.setupHTTP()
	s.setupGRPC()
	s.setupTimeout()
}

func (s *CalendarSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)

	err = s.conn.Close()
	s.Require().NoError(err)
}

func (s *CalendarSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), s.requestTimeout)
}

func (s *CalendarSuite) TearDownTest() {
	s.cancel()
}

func (s *CalendarSuite) TestHTTP() {
	newEvent := cs.Event{
		ID:              "0c2e2081",
		Title:           "event 1",
		Description:     "description 1",
		OwnerID:         "2081",
		StartDate:       time.Date(2023, time.November, 1, 10, 10, 0, 0, time.UTC),
		FinishDate:      time.Date(2023, time.November, 2, 11, 0, 0, 0, time.UTC),
		NotificationDay: time.Date(2023, time.November, 2, 0, 0, 0, 0, time.UTC),
	}

	updatedEvent := cs.Event{
		ID:              "0c2e2081",
		Title:           "event 1",
		Description:     "new description 1",
		OwnerID:         "2081",
		StartDate:       time.Date(2023, time.November, 10, 10, 10, 0, 0, time.UTC),
		FinishDate:      time.Date(2023, time.November, 12, 11, 0, 0, 0, time.UTC),
		NotificationDay: time.Date(2023, time.November, 12, 0, 0, 0, 0, time.UTC),
	}

	tests := [...]struct {
		name      string
		path      string
		method    string
		body      cs.Event
		qres      hs.ResponseEvents
		qexpected hs.ResponseEvents
		cres      hs.ResponseEvent
		cexpected hs.ResponseEvent
	}{
		{
			name:      "get empty events",
			path:      "/event?date=2023-11-01%2006:05:00&period=month",
			method:    http.MethodGet,
			qexpected: hs.ResponseEvents{Events: []cs.Event{}},
		},
		{
			name:      "create event",
			path:      "/event",
			method:    http.MethodPost,
			body:      newEvent,
			cexpected: hs.ResponseEvent{ID: "0c2e2081", Message: "created"},
		},
		{
			name:      "get new events",
			path:      "/event?date=2023-11-01%2006:05:00&period=month",
			method:    http.MethodGet,
			qexpected: hs.ResponseEvents{Events: []cs.Event{newEvent}},
		},
		{
			name:      "update event",
			path:      "/event",
			method:    http.MethodPut,
			body:      updatedEvent,
			cexpected: hs.ResponseEvent{ID: "0c2e2081", Message: "updated"},
		},
		{
			name:      "get update events",
			path:      "/event?date=2023-11-01%2006:05:00&period=month",
			method:    http.MethodGet,
			qexpected: hs.ResponseEvents{Events: []cs.Event{updatedEvent}},
		},
		{
			name:      "delete event",
			path:      "/event?id=0c2e2081",
			method:    http.MethodDelete,
			cexpected: hs.ResponseEvent{ID: "0c2e2081", Message: "deleted"},
		},
		{
			name:      "get empty events",
			path:      "/event?date=2023-11-01%2006:05:00&period=month",
			method:    http.MethodGet,
			qexpected: hs.ResponseEvents{Events: []cs.Event{}},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			var body io.Reader
			if tt.method == http.MethodPost || tt.method == http.MethodPut {
				b, err := json.Marshal(tt.body)
				s.Require().NoError(err)
				body = bytes.NewReader(b)
			}

			req, err := http.NewRequestWithContext(s.ctx, tt.method, s.httpURL+tt.path, body)
			s.Require().NoError(err)

			res, err := s.httpClient.Do(req)
			s.Require().NoError(err)
			defer res.Body.Close()

			s.Require().Equal(http.StatusOK, res.StatusCode)
			s.Require().Equal("application/json", res.Header.Get("Content-Type"))

			out, err := io.ReadAll(res.Body)
			s.Require().NoError(err)

			if tt.method == http.MethodGet {
				err = json.Unmarshal(out, &tt.qres)
				s.Require().NoError(err)
				s.Require().Equal(tt.qexpected, tt.qres)
				return
			}

			err = json.Unmarshal(out, &tt.cres)
			s.Require().NoError(err)
			s.Require().Equal(tt.cexpected, tt.cres)

			// проверка, что приложение произвело изменения именно в нужной базе данных, а не где-либо ещё
			e, err := s.event(tt.cres.ID)
			switch tt.method {
			case http.MethodPost, http.MethodPut:
				s.Require().NoError(err)
				s.Require().Equal(tt.body, e)
			case http.MethodDelete:
				s.Require().Equal(sql.ErrNoRows, err)
			}
		})
	}
}

func (s *CalendarSuite) TestGRPC() {
	req := pb.EventsRequest{
		Date:   timestamppb.New(time.Date(2023, time.December, 1, 4, 10, 0, 0, time.UTC)),
		Period: "month",
	}

	newEvent := pb.Event{
		Id:              "1c2e2081",
		Title:           "grpc event 1",
		Description:     "grpc description",
		OwnerId:         "130c866",
		StartDate:       timestamppb.New(time.Date(2023, time.December, 1, 9, 10, 0, 0, time.UTC)),
		FinishDate:      timestamppb.New(time.Date(2023, time.December, 2, 14, 10, 0, 0, time.UTC)),
		NotificationDay: timestamppb.New(time.Date(2023, time.December, 2, 0, 0, 0, 0, time.UTC)),
	}

	updatedEvent := pb.Event{
		Id:              "1c2e2081",
		Title:           "new grpc event 1",
		Description:     "new grpc description",
		OwnerId:         "130c866",
		StartDate:       timestamppb.New(time.Date(2023, time.December, 11, 4, 10, 0, 0, time.UTC)),
		FinishDate:      timestamppb.New(time.Date(2023, time.December, 12, 14, 10, 0, 0, time.UTC)),
		NotificationDay: timestamppb.New(time.Date(2023, time.December, 12, 0, 0, 0, 0, time.UTC)),
	}

	s.Run("get empty events", func() {
		res, err := s.grpcClient.Events(s.ctx, &req)
		s.Require().NoError(err)
		s.Require().Len(res.GetEvents(), 0)
	})

	s.Run("create event", func() {
		res, err := s.grpcClient.CreateEvent(s.ctx, &newEvent)
		s.Require().NoError(err)
		s.Require().Equal("1c2e2081", res.GetId())
		s.Require().Equal("created", res.GetMessage())

		_, err = s.event(res.GetId())
		s.Require().NoError(err)
	})

	s.Run("get new event", func() {
		res, err := s.grpcClient.Events(s.ctx, &req)
		s.Require().NoError(err)

		ev := res.GetEvents()[0]
		s.Require().Equal(newEvent.GetId(), ev.GetId())
		s.Require().Equal(newEvent.GetTitle(), ev.GetTitle())
		s.Require().Equal(newEvent.GetDescription(), ev.GetDescription())
		s.Require().Equal(newEvent.GetOwnerId(), ev.GetOwnerId())
		s.Require().Equal(newEvent.GetStartDate().AsTime(), ev.GetStartDate().AsTime())
		s.Require().Equal(newEvent.GetFinishDate().AsTime(), ev.GetFinishDate().AsTime())
		s.Require().Equal(newEvent.GetNotificationDay().AsTime(), ev.GetNotificationDay().AsTime())
	})

	s.Run("update event", func() {
		res, err := s.grpcClient.UpdateEvent(s.ctx, &updatedEvent)
		s.Require().NoError(err)
		s.Require().Equal("1c2e2081", res.GetId())
		s.Require().Equal("updated", res.GetMessage())
	})

	s.Run("delete event", func() {
		res, err := s.grpcClient.DeleteEvent(s.ctx, &pb.DeleteRequest{Id: "1c2e2081"})
		s.Require().NoError(err)
		s.Require().Equal("1c2e2081", res.GetId())
		s.Require().Equal("deleted", res.GetMessage())

		_, err = s.event(res.GetId())
		s.Require().Equal(sql.ErrNoRows, err)
	})

	s.Run("get empty events", func() {
		res, err := s.grpcClient.Events(s.ctx, &req)
		s.Require().NoError(err)
		s.Require().Len(res.GetEvents(), 0)
	})
}

func (s *CalendarSuite) setupDelay() {
	delay := os.Getenv("TEST_DELAY")
	if delay == "" {
		return
	}

	d, err := time.ParseDuration(delay)
	s.Require().NoError(err)

	s.T().Logf("wait %s for service availability...", delay)
	time.Sleep(d)
}

func (s *CalendarSuite) setupDB() {
	dsn := os.Getenv("TEST_DB")
	if dsn == "" {
		dsn = "postgresql://db_user:super_secret_password_42@localhost:5432/calendardb"
	}

	db, err := sql.Open("pgx", dsn)
	s.Require().NoError(err)

	err = db.Ping()
	s.Require().NoError(err)

	s.db = db
}

func (s *CalendarSuite) setupHTTP() {
	url := os.Getenv("TEST_HTTP")
	if url == "" {
		url = "http://localhost:3000"
	}

	s.httpURL = url
	s.httpClient = http.DefaultClient
}

func (s *CalendarSuite) setupGRPC() {
	host := os.Getenv("TEST_GRPC")
	if host == "" {
		host = ":50051"
	}

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.conn = conn
	s.grpcClient = pb.NewCalendarClient(conn)
}

func (s *CalendarSuite) setupTimeout() {
	d := os.Getenv("TEST_REQUEST_TIMEOUT")
	if d == "" {
		d = "100ms"
	}

	t, err := time.ParseDuration(d)
	s.Require().NoError(err)

	s.requestTimeout = t
}

func (s *CalendarSuite) event(eventID string) (cs.Event, error) {
	var (
		id, title, desc, oID string
		st, ft, nd           sql.NullTime
	)
	selq := `SELECT id, title, description, owner_id, start_date, finish_date, notification_day
		FROM events
		WHERE id=$1 LIMIT 1
	`

	err := s.db.QueryRow(selq, eventID).Scan(&id, &title, &desc, &oID, &st, &ft, &nd)
	if err != nil {
		return cs.Event{}, err
	}

	return cs.Event{
		ID:              id,
		Title:           title,
		Description:     desc,
		OwnerID:         oID,
		StartDate:       st.Time,
		FinishDate:      ft.Time,
		NotificationDay: nd.Time,
	}, nil
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}
