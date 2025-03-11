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
	mx                *sync.RWMutex
	workerCount       int
	tasks             chan Task
	failed            bool
	stopped           bool
	errorsCount       atomic.Int32
	acceptErrorsCount int
}

func NewWorkerPool(workerCount int, tasks chan Task, acceptErrorsCount int) *WorkerPool {
	return &WorkerPool{
		workerCount:       workerCount,
		tasks:             tasks,
		acceptErrorsCount: acceptErrorsCount,
		wg:                &sync.WaitGroup{},
		mx:                &sync.RWMutex{},
	}
}

func (wp *WorkerPool) addWorker() {
	defer wp.wg.Done()
	for task := range wp.tasks {
		wp.mx.RLock()
		stoppedOrFailed := wp.stopped || wp.failed
		wp.mx.RUnlock()
		if stoppedOrFailed {
			return
		}
		err := task()
		if err != nil {
			wp.errorsCount.Add(1)
		}
	}
}

func (wp *WorkerPool) runErrorsChecker() {
	if wp.acceptErrorsCount <= 0 {
		return
	}
	for {
		wp.mx.Lock()
		if wp.stopped {
			defer wp.mx.Unlock()
			return
		}

		wp.failed = wp.errorsCount.Load() >= int32(wp.acceptErrorsCount)
		wp.mx.Unlock()
	}
}

func (wp *WorkerPool) Stop() {
	wp.mx.Lock()
	wp.stopped = true
	wp.mx.Unlock()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	runTasks := func() chan Task {
		tasksChan := make(chan Task, n)
		go func() {
			for _, task := range tasks {
				tasksChan <- task
			}
			close(tasksChan)
		}()
		return tasksChan
	}

	tasksChan := runTasks()

	wp := NewWorkerPool(n, tasksChan, m)

	for range n {
		wp.wg.Add(1)
		go wp.addWorker()
	}

	go wp.runErrorsChecker()

	wp.wg.Wait()
	wp.Stop()

	if wp.failed {
		return ErrErrorsLimitExceeded
	}

	return nil
}
