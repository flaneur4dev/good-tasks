package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) (err error) {
	if len(tasks) == 0 {
		return
	}
	if n < 1 {
		n = 1
	}

	tasksCh := make(chan Task)
	mistakes := make(chan struct{})
	done := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for t := range tasksCh {
				if err := t(); err != nil && m > 0 {
					select {
					case mistakes <- struct{}{}:
						continue
					case <-done:
						return
					}
				}
			}
		}()
	}

	count := 0

loop:
	for _, t := range tasks {
		select {
		case tasksCh <- t:
			continue
		case <-mistakes:
			count++
			if count == m {
				err = ErrErrorsLimitExceeded
				break loop
			}
			tasksCh <- t
		}
	}

	close(tasksCh)
	close(done)
	wg.Wait()

	return
}
