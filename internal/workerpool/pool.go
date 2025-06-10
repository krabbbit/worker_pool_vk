// Package workerpool provides concurrent job processing using a configurable worker pool.
package workerpool

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// New creates a new WorkerPool with the specified job queue capacity.
// The pool starts with zero workers - call AddWorker() to add them.
func New(jobs int) WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &pool{
		jobs:    make(chan string, jobs),
		wg:      &sync.WaitGroup{},
		workers: []*worker{},
		mu:      &sync.Mutex{},
		ctx:     ctx,
		cancel:  cancel,
	}
}

// GetWorkersCount returns the current number of active workers in the pool.
func (wp *pool) GetWorkersCount() int {
	return len(wp.workers)
}

// AddWorker creates and adds a new worker to the pool.
// Returns error if the pool is closed.
func (wp *pool) AddWorker() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	select {
	case <-wp.ctx.Done():
		return fmt.Errorf("workerpool is closed")
	default:
	}

	id := len(wp.workers)
	exit := make(chan struct{})
	wp.workers = append(wp.workers, &worker{id: id, exit: exit})
	wp.wg.Add(1)

	go wp.work(id, exit, doJob)

	return nil
}

func doJob(job string, id int) {
	time.Sleep(time.Duration(500+rand.Intn(1500)) * time.Millisecond)
	fmt.Printf("worker_id: %d, job: %s \n", id, job)
}

func (wp *pool) work(id int, exit chan struct{}, doJob func(job string, id int)) {
	defer wp.wg.Done()
	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			doJob(job, id)
		case <-exit:
			return
		case <-wp.ctx.Done():
			return
		}
	}
}

// AddJob adds a new job to the worker pool queue.
// Returns error if the pool is closed or full.
func (wp *pool) AddJob(job string) error {
	select {
	case <-wp.ctx.Done():
		return fmt.Errorf("workerpool is closed")
	default:
	}

	select {
	case wp.jobs <- job:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("workerpool is closed")
	}
}

// Delete removes the most recently added worker from the pool.
// Returns error if the pool is empty.
func (wp *pool) Delete() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if len(wp.workers) == 0 {
		return fmt.Errorf("workerpool is empty")
	}

	last := len(wp.workers) - 1
	close(wp.workers[last].exit)
	wp.workers = wp.workers[:last]
	return nil
}

// Shutdown gracefully stops all workers and cleans up resources.
// Waits for all in-progress jobs to complete.
func (wp *pool) Shutdown() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.cancel()
	for _, w := range wp.workers {
		close(w.exit)
	}
	close(wp.jobs)
	wp.workers = nil
	wp.wg.Wait()
}

// Wait blocks until all workers have finished processing.
func (wp *pool) Wait() {
	wp.wg.Wait()
}
