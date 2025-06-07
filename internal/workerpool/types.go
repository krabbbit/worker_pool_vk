package workerpool

import (
	"context"
	"sync"
)

type Pool struct {
	wg      *sync.WaitGroup
	jobs    chan string
	workers []*worker
	mu      *sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
}

type worker struct {
	id   int
	exit chan struct{}
}

type WorkerPool interface {
	Add() error
	AddJob(job string) error
	Delete() error
	Shutdown()
}
