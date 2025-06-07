package workerpool

import (
	"context"
	"fmt"
	"sync"
)

func NewWorkerpool(jobs int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		jobs:    make(chan string, jobs),
		wg:      &sync.WaitGroup{},
		workers: []*worker{},
		mu:      &sync.Mutex{},
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (wp *Pool) GetWorkers() []*worker {
	return wp.workers
}

func (wp *Pool) Add() error {
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
	go func(id int, exit chan struct{}) {
		defer wp.wg.Done()
		for {
			select {
			case job, ok := <-wp.jobs:
				if !ok {
					return
				}
				fmt.Printf("worker_id: %d, job: %s \n", id, job)
			case <-exit:
				return
			case <-wp.ctx.Done():
				return
			}
		}
	}(id, exit)
	return nil
}

func (wp *Pool) AddJob(job string) error {
	select {
	case wp.jobs <- job:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("workerpool is closed")
	}
}

func (wp *Pool) Delete() error {
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

func (wp *Pool) Shutdown() {
	wp.cancel()
	wp.mu.Lock()
	defer wp.mu.Unlock()
	for _, w := range wp.workers {
		close(w.exit)
	}
	close(wp.jobs)
	wp.wg.Wait()
}

func (wp *Pool) Wait() {
	wp.wg.Wait()
}
