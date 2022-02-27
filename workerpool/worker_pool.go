package workerpool

import (
	"context"
	"sync"
	"fmt"
	"time"
)

type Task struct {
	Func func(args ...interface{}) *Result
	Args []interface{}
}

type Result struct {
	Value interface{}
	Err   error
}

type WorkerPool interface {
	Start(ctx context.Context)
	Tasks() chan *Task
	Results() chan *Result
}

type workerPool struct {
	numWorkers int
	tasks      chan *Task
	results    chan *Result
	wg         *sync.WaitGroup
}

var _ WorkerPool = (*workerPool)(nil)

func NewWorkerPool(numWorkers int, bufferSize int) *workerPool {
	return &workerPool{
		numWorkers: numWorkers,
		tasks:      make(chan *Task, bufferSize),
		results:    make(chan *Result, bufferSize),
		wg:         &sync.WaitGroup{},
	}
}

func (wp *workerPool) Start(ctx context.Context) {
	// TODO: implementation
	//
	// Starts numWorkers of goroutines, wait until all jobs are done.
	// Remember to closed the result channel before exit.
	ctx2, cancel := context.WithCancel(ctx)
	ctx3, _ := context.WithTimeout(ctx, 1* time.Second)
	for i:=0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.run(ctx2)
	}
	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutdownctx")
			wp.wg.Wait()
			close(wp.results)
			cancel()
			return
		case <-ctx3.Done():
			fmt.Println("timeout")
			cancel()
			wp.wg.Wait()
			close(wp.results)
			return
		}
	}
}

func (wp *workerPool) Tasks() chan *Task {
	return wp.tasks
}

func (wp *workerPool) Results() chan *Result {
	return wp.results
}

func (wp *workerPool) run(ctx context.Context) {
	// TODO: implementation
	//
	// Keeps fetching task from the task channel, do the task,
	// then makes sure to exit if context is done.
	defer wp.wg.Done()
	for {
		select {
		case shutdown := <-ctx.Done():
			fmt.Println(shutdown)
			return
		case t,ok := <-wp.tasks:
			if ok {
				result := t.Func(t.Args...)
				wp.results <- result
			}
		}
	}
}