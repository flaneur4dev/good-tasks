package memory

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	cs "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/contracts"
	es "github.com/flaneur4dev/good-tasks/hw12_13_14_15_calendar/internal/mistakes"
)

func TestMemStore(t *testing.T) {
	ms := New()

	t.Run("concurrency create", func(t *testing.T) {
		events := [...]cs.Event{
			{
				ID:               "1",
				Title:            "event 1",
				Description:      "description 1",
				OwnerID:          "11",
				StartDate:        time.Date(2023, time.September, 01, 04, 0, 0, 0, time.UTC),
				FinishDate:       time.Date(2023, time.September, 01, 13, 0, 0, 0, time.UTC),
				NotificationTime: time.Date(2023, time.September, 01, 10, 0, 0, 0, time.UTC),
			},
			{
				ID:               "2",
				Title:            "event 2",
				Description:      "description 2",
				OwnerID:          "12",
				StartDate:        time.Date(2023, time.September, 04, 10, 5, 0, 0, time.UTC),
				FinishDate:       time.Date(2023, time.September, 05, 23, 0, 0, 0, time.UTC),
				NotificationTime: time.Date(2023, time.September, 05, 20, 0, 0, 0, time.UTC),
			},
			{
				ID:               "3",
				Title:            "event 3",
				Description:      "description 3",
				OwnerID:          "13",
				StartDate:        time.Date(2023, time.August, 28, 07, 0, 0, 0, time.UTC),
				FinishDate:       time.Date(2023, time.August, 28, 13, 0, 0, 0, time.UTC),
				NotificationTime: time.Date(2023, time.August, 28, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:               "4",
				Title:            "event 4",
				Description:      "description 4",
				OwnerID:          "14",
				StartDate:        time.Date(2023, time.October, 11, 4, 0, 0, 0, time.UTC),
				FinishDate:       time.Date(2023, time.October, 12, 16, 0, 0, 0, time.UTC),
				NotificationTime: time.Date(2023, time.October, 12, 15, 0, 0, 0, time.UTC),
			},
			{
				ID:               "5",
				Title:            "event 5",
				Description:      "description 5",
				OwnerID:          "15",
				StartDate:        time.Date(2023, time.September, 23, 4, 0, 0, 0, time.UTC),
				FinishDate:       time.Date(2023, time.September, 23, 13, 0, 0, 0, time.UTC),
				NotificationTime: time.Date(2023, time.September, 23, 10, 0, 0, 0, time.UTC),
			},
		}

		workersCount := 10
		var maxErrorsCount int64
		wg := sync.WaitGroup{}
		wg.Add(workersCount)

		for i := 0; i < workersCount; i++ {
			go func() {
				defer wg.Done()

				for _, ev := range events {
					if err := ms.CreateEvent(context.Background(), ev); errors.Is(err, es.ErrDateBusy) {
						atomic.AddInt64(&maxErrorsCount, 1)
					}
				}
			}()
		}

		wg.Wait()

		expected := workersCount*len(events) - len(events)
		if int(maxErrorsCount) != expected {
			t.Fatalf("[create error] want: %d, but got: %d", expected, maxErrorsCount)
		}
	})

	t.Run("update", func(t *testing.T) {
		tests := [...]struct {
			input  cs.Event
			output error
		}{
			{
				input: cs.Event{
					ID:               "1",
					Title:            "event 1",
					Description:      "new description",
					OwnerID:          "11",
					StartDate:        time.Date(2023, time.September, 01, 04, 0, 0, 0, time.UTC),
					FinishDate:       time.Date(2023, time.September, 01, 13, 0, 0, 0, time.UTC),
					NotificationTime: time.Date(2023, time.September, 01, 10, 0, 0, 0, time.UTC),
				},
				output: nil,
			},
			{
				input: cs.Event{
					ID:               "42",
					Title:            "event 1",
					Description:      "description 1",
					OwnerID:          "11",
					StartDate:        time.Date(2020, time.December, 01, 04, 0, 0, 0, time.UTC),
					FinishDate:       time.Date(2020, time.December, 01, 13, 0, 0, 0, time.UTC),
					NotificationTime: time.Date(2020, time.December, 01, 10, 0, 0, 0, time.UTC),
				},
				output: es.ErrNoEvent,
			},
		}

		for _, tt := range tests {
			if err := ms.UpdateEvent(context.Background(), tt.input.ID, tt.input); !errors.Is(err, tt.output) {
				t.Fatalf("[update error] want: %v, but got: %v", tt.output, err)
			}
		}
	})
}
