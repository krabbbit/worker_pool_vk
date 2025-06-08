package main

import (
	"fmt"
	"github.com/krabbbit/worker_pool_vk/internal/workerpool"
	"strconv"
)

func ExampleWorkerPool() {
	wp := workerpool.New(5)
	for i := 0; i < 5; i++ {
		_ = wp.AddWorker()
	}
	for i := 0; i < 4; i++ {
		_ = wp.AddJob(strconv.Itoa(i))
	}
	for i := 0; i < 3; i++ {
		_ = wp.Delete()
	}
	for i := 0; i < 4; i++ {
		_ = wp.AddJob(strconv.Itoa(i + 100))
	}
	wp.Shutdown()
	err := wp.AddWorker()
	if err != nil {
		fmt.Println("wow, error ", err.Error())
	}

}
