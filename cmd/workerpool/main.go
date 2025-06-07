// package workerpool
package main

import (
	"fmt"
	"strconv"
	"worker_pool_vk/internal/workerpool"
)

func main() {
	wp := workerpool.NewWorkerpool(5)
	for i := 0; i < 5; i++ {
		_ = wp.Add()
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
	err := wp.Add()
	if err != nil {
		fmt.Println("gdgf")
	}

}
