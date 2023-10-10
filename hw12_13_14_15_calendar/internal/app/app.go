package app

import (
	"context"
	"time"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	es "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/mistakes"
	"github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/utils"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type Storage interface {
	Events(ctx context.Context, start, end time.Time) ([]cs.Event, error)
	CreateEvent(ctx context.Context, ne cs.Event) error
	UpdateEvent(ctx context.Context, id string, ne cs.Event) error
	DeleteEvent(ctx context.Context, id string) error
	Check(ctx context.Context) error
	Close() error
}

type App struct {
	log     Logger
	storage Storage
}

func New(log Logger, storage Storage) *App {
	return &App{log, storage}
}

func (a *App) Events(ctx context.Context, date time.Time, period string) ([]cs.Event, error) {
	var end time.Time
	switch period {
	case "day":
		end = date.AddDate(0, 0, 1)

	case "week":
		if m := date.Weekday(); m != time.Monday {
			return nil, es.ErrPeriod
		}
		end = date.AddDate(0, 0, 7)

	case "month":
		if w := date.Day(); w != 1 {
			return nil, es.ErrPeriod
		}
		end = date.AddDate(0, 1, 0)

	default:
		return nil, es.ErrPeriod
	}

	es, err := a.storage.Events(ctx, date, end)
	if err != nil {
		a.log.Error("failed to get events: " + err.Error())
		return nil, err
	}

	return es, nil
}

func (a *App) CreateEvent(ctx context.Context, ne cs.Event) error {
	ne.NotificationDay = utils.ZeroTime(ne.NotificationDay)

	err := a.storage.CreateEvent(ctx, ne)
	if err != nil {
		a.log.Error("failed to create event: " + err.Error())
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, id string, ne cs.Event) error {
	err := a.storage.UpdateEvent(ctx, id, ne)
	if err != nil {
		a.log.Error("failed to update event: " + err.Error())
		return err
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	err := a.storage.DeleteEvent(ctx, id)
	if err != nil {
		a.log.Error("failed to delete event: " + err.Error())
		return err
	}

	return nil
}

func (a *App) Check(ctx context.Context) error {
	err := a.storage.Check(ctx)
	if err != nil {
		a.log.Error("failed to check storage: " + err.Error())
		return err
	}

	return nil
}
