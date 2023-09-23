package memory

import (
	"context"
	"sync"
	"time"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	es "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/mistakes"
)

type MemStore struct {
	mu     sync.RWMutex
	events map[string]cs.Event
}

func New() *MemStore {
	return &MemStore{
		events: map[string]cs.Event{},
	}
}

func (ms *MemStore) Events(ctx context.Context, start, end time.Time) ([]cs.Event, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	es := []cs.Event{}
	for _, v := range ms.events {
		eventDay := v.StartDate.YearDay()
		if start.YearDay() <= eventDay && eventDay < end.YearDay() {
			es = append(es, v)
		}
	}

	return es, nil
}

func (ms *MemStore) CreateEvent(ctx context.Context, ne cs.Event) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	for _, v := range ms.events {
		eventTime := ne.StartDate.Unix()
		if v.StartDate.Unix() <= eventTime && eventTime <= v.FinishDate.Unix() {
			return es.ErrDateBusy
		}
	}

	ms.events[ne.ID] = ne
	return nil
}

func (ms *MemStore) UpdateEvent(ctx context.Context, id string, ne cs.Event) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, ok := ms.events[id]
	if !ok {
		return es.ErrNoEvent
	}

	ms.events[id] = ne
	return nil
}

func (ms *MemStore) DeleteEvent(ctx context.Context, id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	delete(ms.events, id)
	return nil
}

func (ms *MemStore) Check(_ context.Context) error {
	return nil
}

func (ms *MemStore) Close() error {
	return nil
}
