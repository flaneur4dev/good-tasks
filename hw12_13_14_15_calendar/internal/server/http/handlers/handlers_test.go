package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/app"
	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestHandlersStepByStep(t *testing.T) {
	a := app.New(slog.Default(), memory.New())

	newEvent := cs.Event{
		ID:               "0c2e2081",
		Title:            "event 1",
		Description:      "description 1",
		OwnerID:          "2081",
		StartDate:        time.Date(2023, time.September, 1, 10, 10, 0, 0, time.UTC),
		FinishDate:       time.Date(2023, time.September, 2, 11, 0, 0, 0, time.UTC),
		NotificationTime: time.Date(2023, time.September, 2, 10, 0, 0, 0, time.UTC),
	}

	updatedEvent := cs.Event{
		ID:               "0c2e2081",
		Title:            "event 1",
		Description:      "new description 1",
		OwnerID:          "2081",
		StartDate:        time.Date(2023, time.September, 10, 10, 10, 0, 0, time.UTC),
		FinishDate:       time.Date(2023, time.September, 12, 11, 0, 0, 0, time.UTC),
		NotificationTime: time.Date(2023, time.September, 12, 10, 0, 0, 0, time.UTC),
	}

	tests := [...]struct {
		name      string
		path      string
		method    string
		handler   http.Handler
		body      cs.Event
		qres      ResponseEvents
		qexpected ResponseEvents
		cres      ResponseEvent
		cexpected ResponseEvent
	}{
		{
			name:      "get empty events",
			path:      "/event?date=2023-09-01%2006:05:00&period=month",
			method:    http.MethodGet,
			handler:   HandleEvents(a),
			qexpected: ResponseEvents{[]cs.Event{}},
		},
		{
			name:      "create event",
			path:      "/event",
			method:    http.MethodPost,
			handler:   HandleCreateEvent(a),
			body:      newEvent,
			cexpected: ResponseEvent{"0c2e2081", "created"},
		},
		{
			name:      "get new events",
			path:      "/event?date=2023-09-01%2006:05:00&period=month",
			method:    http.MethodGet,
			handler:   HandleEvents(a),
			qexpected: ResponseEvents{[]cs.Event{newEvent}},
		},
		{
			name:      "update event",
			path:      "/event",
			method:    http.MethodPut,
			handler:   HandleUpdateEvent(a),
			body:      updatedEvent,
			cexpected: ResponseEvent{"0c2e2081", "updated"},
		},
		{
			name:      "get update events",
			path:      "/event?date=2023-09-01%2006:05:00&period=month",
			method:    http.MethodGet,
			handler:   HandleEvents(a),
			qexpected: ResponseEvents{[]cs.Event{updatedEvent}},
		},
		{
			name:      "delete event",
			path:      "/event?id=0c2e2081",
			method:    http.MethodDelete,
			handler:   HandleDeleteEvent(a),
			cexpected: ResponseEvent{"0c2e2081", "deleted"},
		},
		{
			name:      "get empty events",
			path:      "/event?date=2023-09-01%2006:05:00&period=month",
			method:    http.MethodGet,
			handler:   HandleEvents(a),
			qexpected: ResponseEvents{[]cs.Event{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			var body io.Reader
			if tt.method == http.MethodPost || tt.method == http.MethodPut {
				b, err := json.Marshal(tt.body)
				require.NoError(t, err)
				body = bytes.NewReader(b)
			}

			req, err := http.NewRequest(tt.method, ts.URL+tt.path, body)
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, "application/json", res.Header.Get("Content-Type"))

			out, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if tt.method == http.MethodGet {
				err = json.Unmarshal(out, &tt.qres)
				require.NoError(t, err)
				require.Equal(t, tt.qexpected, tt.qres)
			} else {
				err = json.Unmarshal(out, &tt.cres)
				require.NoError(t, err)
				require.Equal(t, tt.cexpected, tt.cres)
			}
		})
	}
}
