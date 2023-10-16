package main

import (
	"context"
	"io"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
)

type messageQueue interface {
	Consume(ctx context.Context) (<-chan cs.NotificationMessage, error)
	Close() error
}

type sender struct {
	mq messageQueue
	w  io.Writer
}

func newSender(mq messageQueue, w io.Writer) *sender {
	return &sender{
		mq: mq,
		w:  w,
	}
}

func (s *sender) start(ctx context.Context) error {
	messages, err := s.mq.Consume(ctx)
	if err != nil {
		return err
	}

	for m := range messages {
		if _, err := s.w.Write(m.Body); err != nil {
			return err
		}
	}

	return nil
}

func (s *sender) stop() error {
	return s.mq.Close()
}
