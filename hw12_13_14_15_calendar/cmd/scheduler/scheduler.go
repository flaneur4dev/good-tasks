package main

import (
	"context"
	"encoding/json"
	"time"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
)

type storage interface {
	Notifications(ctx context.Context, t time.Time) ([]cs.Notification, error)
	Clear(ctx context.Context, t time.Time) error
	Close() error
}

type messageQueue interface {
	Publish(ctx context.Context, body []byte) error
	Close() error
}

type scheduler struct {
	ticker *time.Ticker
	store  storage
	mq     messageQueue
}

func newScheduler(interval string, store storage, mq messageQueue) (*scheduler, error) {
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}

	return &scheduler{
		ticker: time.NewTicker(d),
		store:  store,
		mq:     mq,
	}, nil
}

func (s *scheduler) execute(ctx context.Context, t time.Time) error {
	ns, err := s.store.Notifications(ctx, t)
	if err != nil {
		return err
	}

	for _, n := range ns {
		b, err := json.Marshal(n)
		if err != nil {
			return err
		}

		err = s.mq.Publish(ctx, b)
		if err != nil {
			return err
		}
	}

	err = s.store.Clear(ctx, t.AddDate(-1, 0, 0))
	if err != nil {
		return err
	}

	return nil
}

func (s *scheduler) stop() error {
	s.ticker.Stop()

	err := s.store.Close()
	if err != nil {
		return err
	}

	err = s.mq.Close()
	if err != nil {
		return err
	}

	return nil
}
