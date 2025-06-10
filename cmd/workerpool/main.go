package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/krabbbit/worker_pool_vk/internal/workerpool"
	"os"
	"strings"
)

func main() {
	jobsNum := flag.Int("cap", 10, "job queue capacity")
	flag.Parse()

	wp := workerpool.New(*jobsNum)
	fmt.Printf("Worker pool started (jobsNum: %d)\n", *jobsNum)
	fmt.Println("Type 'help' for commands")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		cmd := strings.TrimSpace(scanner.Text())
		switch cmd {
		case "addw":
			if err := wp.AddWorker(); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Worker added")
			}

		case "delw":
			if err := wp.Delete(); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Worker deleted")
			}

		case "addj":
			fmt.Print("Enter job: ")
			if !scanner.Scan() {
				break
			}
			job := scanner.Text()
			if err := wp.AddJob(job); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Job added")
			}

		case "shutdown":
			wp.Shutdown()
			fmt.Println("Pool shutdown")
			return

		case "help":
			fmt.Println("Commands:")
			fmt.Println("  addw  - add worker")
			fmt.Println("  delw  - delete worker")
			fmt.Println("  addj  - add job")
			fmt.Println("  shutdown - shutdown pool")
			fmt.Println("  help  - show this help")

		default:
			fmt.Println("Unknown command. Type 'help'")
		}
	}
}
