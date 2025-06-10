package workerpool

import (
	"context"
	"sync"
)

// pool is the internal implementation of WorkerPool.
// It manages:
// - wg: wait group for tracking active workers
// - jobs: buffered channel for incoming jobs
// - workers: slice of active worker instances
// - mu: mutex for concurrent access protection
// - ctx/cancel: context for graceful shutdown
type pool struct {
	wg      *sync.WaitGroup
	jobs    chan string
	workers []*worker
	mu      *sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// worker represents an individual worker in the pool.
// Contains:
// - id: unique identifier for the worker
// - exit: channel for signaling worker termination
type worker struct {
	id   int
	exit chan struct{}
}

// WorkerPool manages a pool of workers for concurrent job processing.
// Handles worker lifecycle, job distribution, and graceful shutdown.
// - AddWorker: adds new worker to the pool
// - AddJob: enqueues job for processing
// - Delete: removes most recent worker
// - Shutdown: initiates graceful shutdown
// - Wait: blocks until all workers complete
// - GetWorkersCount: returns current worker count
type WorkerPool interface {
	AddWorker() error
	AddJob(job string) error
	Delete() error
	Shutdown()
	Wait()
	GetWorkersCount() int
}
