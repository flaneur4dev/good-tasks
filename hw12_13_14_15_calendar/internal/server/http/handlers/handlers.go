package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	es "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/mistakes"
)

type (
	Events interface {
		Events(ctx context.Context, date time.Time, period string) ([]cs.Event, error)
	}
	Creator interface {
		CreateEvent(ctx context.Context, ne cs.Event) error
	}
	Updater interface {
		UpdateEvent(ctx context.Context, id string, ne cs.Event) error
	}
	Deleter interface {
		DeleteEvent(ctx context.Context, id string) error
	}
	Checker interface {
		Check(ctx context.Context) error
	}
)

type (
	ResponseEvents struct {
		Events []cs.Event `json:"events"`
	}
	ResponseEvent struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}
)

var periods = map[string]struct{}{
	"day":   {},
	"week":  {},
	"month": {},
}

func HandleEvents(srv Events) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := time.Parse(time.DateTime, r.URL.Query().Get("date"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := r.URL.Query().Get("period")
		if _, ok := periods[p]; !ok {
			http.Error(w, "incorrect period", http.StatusBadRequest)
			return
		}

		evs, err := srv.Events(r.Context(), d, p)
		if err != nil {
			sc := http.StatusInternalServerError
			if errors.Is(err, es.ErrPeriod) {
				sc = http.StatusBadRequest
			}
			http.Error(w, err.Error(), sc)
			return
		}

		v := ResponseEvents{evs}
		b, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleCreateEvent(srv Creator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var e cs.Event
		err = json.Unmarshal(reqBody, &e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = srv.CreateEvent(r.Context(), e)
		if err != nil {
			sc := http.StatusInternalServerError
			if errors.Is(err, es.ErrDateBusy) {
				sc = http.StatusBadRequest
			}
			http.Error(w, err.Error(), sc)
			return
		}

		v := ResponseEvent{e.ID, "created"}
		b, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleUpdateEvent(srv Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var e cs.Event
		err = json.Unmarshal(reqBody, &e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = srv.UpdateEvent(r.Context(), e.ID, e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v := ResponseEvent{e.ID, "updated"}
		b, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleDeleteEvent(srv Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		err := srv.DeleteEvent(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v := ResponseEvent{id, "deleted"}
		b, err := json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleCheck(srv Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := srv.Check(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}
