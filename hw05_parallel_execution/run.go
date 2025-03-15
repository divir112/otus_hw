package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type WorkerPool struct {
	wg                *sync.WaitGroup
	tasks             chan Task
	errorsCount       atomic.Int32
	acceptErrorsCount int
}

func NewWorkerPool(tasks chan Task, acceptErrorsCount int) *WorkerPool {
	return &WorkerPool{
		tasks:             tasks,
		acceptErrorsCount: acceptErrorsCount,
		wg:                &sync.WaitGroup{},
		errorsCount:       atomic.Int32{},
	}
}

func (wp *WorkerPool) addWorker() {
	defer wp.wg.Done()
	for task := range wp.tasks {
		err := task()
		if err != nil {
			wp.errorsCount.Add(1)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 0 || m < 0 {
		return errors.New("n and m must not be less than zero")
	}
	tasksChan := make(chan Task)

	wp := NewWorkerPool(tasksChan, m)

	for i := 0; i < n; i++ {
		wp.wg.Add(1)
		go wp.addWorker()
	}

	for _, task := range tasks {
		if int(wp.errorsCount.Load()) >= wp.acceptErrorsCount {
			break
		}
		tasksChan <- task
	}
	close(tasksChan)

	wp.wg.Wait()

	if int(wp.errorsCount.Load()) >= wp.acceptErrorsCount {
		return ErrErrorsLimitExceeded
	}

	return nil
}
