package httpserver

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	hs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/server/http/handlers"
)

type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type application interface {
	Events(ctx context.Context, date time.Time, period string) ([]cs.Event, error)
	CreateEvent(ctx context.Context, ne cs.Event) error
	UpdateEvent(ctx context.Context, id string, ne cs.Event) error
	DeleteEvent(ctx context.Context, id string) error
	Check(ctx context.Context) error
}

type Server struct {
	log logger
	srv *http.Server
	fd  *os.File
}

func New(log logger, app application, logPath, addr string, timeout, idleTimeout time.Duration) (*Server, error) {
	fd, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	s := &Server{log: log, fd: fd}

	r := chi.NewRouter()
	r.Use(s.loggingMiddleware)

	r.Get("/ping", hs.HandleCheck(app))
	r.Get("/event", hs.HandleEvents(app))
	r.Post("/event", hs.HandleCreateEvent(app))
	r.Put("/event", hs.HandleUpdateEvent(app))
	r.Delete("/event", hs.HandleDeleteEvent(app))

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  idleTimeout,
	}

	s.srv = srv
	return s, nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return s.fd.Close()
}
