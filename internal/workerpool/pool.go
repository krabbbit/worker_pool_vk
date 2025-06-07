package workerpool

import (
	"context"
	"fmt"
	"sync"
)

func New(jobs int) *Pool {
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

func (wp *Pool) Add() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	id := len(wp.workers)
	wp.workers = append(wp.workers, &worker{id: id, exit: make(chan struct{})})

	wp.wg.Add(1)
	go func(id int) {
		defer wp.wg.Done()
		for {
			select {
			case job, ok := <-wp.jobs:
				if !ok {
					return
				}
				fmt.Println(id, job)
			case <-wp.workers[id].exit:
				return
			case <-wp.ctx.Done():
				return
			}
		}
	}(id)
}

func (wp *Pool) Delete() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if len(wp.workers) == 0 {
		return fmt.Errorf("worker pool is empty")
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
