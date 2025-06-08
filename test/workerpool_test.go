package workerpool

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/krabbbit/worker_pool_vk/internal/workerpool"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	t.Run("NewWorkerpool", TestNewWorkerpool)
	t.Run("AddWorker", TestAddWorker)
	t.Run("AddJob", TestAddJob)
	t.Run("DeleteWorker", TestDeleteWorker)
	t.Run("Shutdown", TestShutdown)
	t.Run("ConcurrentOperations", TestConcurrentOperations)
}

func TestNewWorkerpool(t *testing.T) {
	t.Run("should create wp with correct capacity", func(t *testing.T) {
		wp := workerpool.New(10)
		assert.NotNil(t, wp)
		defer wp.Shutdown()
	})
}

func TestAddWorker(t *testing.T) {
	t.Run("should add worker successfully", func(t *testing.T) {
		wp := workerpool.New(10)
		defer wp.Shutdown()

		err := wp.AddWorker()
		assert.NoError(t, err)
		assert.Equal(t, 1, wp.GetWorkersCount())
	})

	t.Run("should return error when wp is closed", func(t *testing.T) {
		wp := workerpool.New(10)
		wp.Shutdown()

		err := wp.AddWorker()
		assert.Error(t, err)
		assert.Equal(t, "workerpool is closed", err.Error())
	})
}

func TestAddJob(t *testing.T) {
	t.Run("should add job successfully", func(t *testing.T) {
		wp := workerpool.New(10)
		defer wp.Shutdown()
		_ = wp.AddWorker()

		err := wp.AddJob("test job")
		assert.NoError(t, err)
	})

	t.Run("should return error when wp is closed", func(t *testing.T) {
		wp := workerpool.New(10)
		wp.Shutdown()

		err := wp.AddJob("test job")
		assert.Error(t, err)
		assert.Equal(t, "workerpool is closed", err.Error())
	})
}

func TestDeleteWorker(t *testing.T) {
	t.Run("should delete worker successfully", func(t *testing.T) {
		wp := workerpool.New(10)
		defer wp.Shutdown()
		_ = wp.AddWorker()
		_ = wp.AddWorker()

		err := wp.Delete()
		assert.NoError(t, err)
		assert.Equal(t, 1, wp.GetWorkersCount())
	})

	t.Run("should return error when no workers", func(t *testing.T) {
		wp := workerpool.New(10)
		defer wp.Shutdown()

		err := wp.Delete()
		assert.Error(t, err)
		assert.Equal(t, "workerpool is empty", err.Error())
	})
}

func TestShutdown(t *testing.T) {
	t.Run("should shutdown gracefully", func(t *testing.T) {
		wp := workerpool.New(10)
		_ = wp.AddWorker()
		_ = wp.AddWorker()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			wp.Shutdown()
		}()

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("Shutdown took too long")
		}

		err := wp.AddJob("test")
		assert.Error(t, err)
	})
}

func TestConcurrentOperations(t *testing.T) {
	wp := workerpool.New(100)
	defer wp.Shutdown()

	var wg sync.WaitGroup
	workers := 10
	jobs := 50

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := wp.AddWorker()
			assert.NoError(t, err)
		}()
	}

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := wp.AddJob(fmt.Sprintf("job%d", i))
			assert.NoError(t, err)
		}(i)
	}

	for i := 0; i < workers/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := wp.Delete()
			if err != nil && !errors.Is(err, errors.New("workerpool is empty")) {
				assert.NoError(t, err)
			}
		}()
	}

	wg.Wait()
	assert.True(t, wp.GetWorkersCount() > 0)
}
