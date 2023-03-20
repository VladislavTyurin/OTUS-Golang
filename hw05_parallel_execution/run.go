package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	currentTaskID := 0
	errorsCount := 0

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				mutex.Lock()
				if currentTaskID < len(tasks) && errorsCount < m {
					t := tasks[currentTaskID]
					currentTaskID++
					mutex.Unlock()

					if t() != nil {
						mutex.Lock()
						errorsCount++
						mutex.Unlock()
					}
				} else {
					mutex.Unlock()
					break
				}
			}
		}()
	}
	wg.Wait()

	if errorsCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
