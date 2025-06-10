package workerpool

import (
	"fmt"
	"github.com/krabbbit/worker_pool_vk/internal/workerpool"
	"strconv"
)

func ExampleWorkerPool() {
	// Create a new worker pool with initial capacity of 5 workers
	wp := workerpool.New(5)

	// Add 5 workers to the pool (total now: 5)
	for i := 0; i < 5; i++ {
		_ = wp.AddWorker()
	}

	// Add 4 jobs to be processed by the worker pool
	for i := 0; i < 4; i++ {
		_ = wp.AddJob(strconv.Itoa(i))
	}

	// Remove 3 workers from the pool (total now: 2)
	for i := 0; i < 3; i++ {
		_ = wp.Delete()
	}

	// Add 4 more jobs to the pool
	for i := 0; i < 4; i++ {
		_ = wp.AddJob(strconv.Itoa(i + 100))
	}

	// Shutdown the worker pool, gracefully terminating all workers
	// after they finish their current jobs
	wp.Shutdown()

	// Attempt to add a new worker after shutdown - this will fail
	err := wp.AddWorker()
	if err != nil {
		// Print error message if worker cannot be added (expected after shutdown)
		fmt.Println("wow, error ", err.Error())
	}

}
