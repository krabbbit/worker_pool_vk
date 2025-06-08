package workerpool

import (
	"context"
	"fmt"
	"sync"
)

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

func (wp *pool) GetWorkersCount() int {
	return len(wp.workers)
}

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

	go wp.work(id, exit)

	return nil
}

func (wp *pool) work(id int, exit chan struct{}) {
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
}

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

func (wp *pool) Wait() {
	wp.wg.Wait()
}
