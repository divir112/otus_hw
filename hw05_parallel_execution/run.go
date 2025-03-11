package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type WorkerPool struct {
	wg                *sync.WaitGroup
	mx                *sync.Mutex
	workerCount       int
	tasks             chan Task
	hasError          bool
	stopped           bool
	errorsCount       int
	acceptErrorsCount int
}

func NewWorkerPool(workerCount int, tasks chan Task, acceptErrorsCount int) *WorkerPool {
	return &WorkerPool{
		workerCount:       workerCount,
		tasks:             tasks,
		acceptErrorsCount: acceptErrorsCount,
		wg:                &sync.WaitGroup{},
		mx:                &sync.Mutex{},
	}
}

func (wp *WorkerPool) addWorker() {
	defer wp.wg.Done()
	for task := range wp.tasks {
		if wp.hasError || wp.stopped {
			return
		}
		err := task()
		if err != nil {
			wp.mx.Lock()
			wp.errorsCount++
			wp.mx.Unlock()
		}
	}
}

func (wp *WorkerPool) runErrorsChecker() {
	if wp.acceptErrorsCount <= 0 {
		return
	}
	for {
		if wp.errorsCount >= wp.acceptErrorsCount || wp.stopped {
			wp.hasError = true
			return
		}
	}
}

func (wp *WorkerPool) Stop() {
	wp.stopped = true
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

	if wp.hasError {
		return ErrErrorsLimitExceeded
	}

	return nil
}
