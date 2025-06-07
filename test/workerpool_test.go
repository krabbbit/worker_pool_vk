package workerpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worker_pool_vk/internal/workerpool"
)

func TestNewWorkerpool(t *testing.T) {
	t.Run("should create wp with correct capacity", func(t *testing.T) {
		wp := workerpool.NewWorkerpool(10)
		assert.NotNil(t, wp)
		defer wp.Shutdown()
	})
}

func TestAddWorker(t *testing.T) {
	t.Run("should add worker successfully", func(t *testing.T) {
		wp := workerpool.NewWorkerpool(10)
		defer wp.Shutdown()

		err := wp.Add()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(wp.GetWorkers()))
	})

	t.Run("should return error when wp is closed", func(t *testing.T) {
		wp := workerpool.NewWorkerpool(10)
		wp.Shutdown()

		err := wp.Add()
		assert.Error(t, err)
		assert.Equal(t, "workerpool is closed", err.Error())
	})
}

func TestAddJob(t *testing.T) {
	t.Run("should add job successfully", func(t *testing.T) {
		wp := workerpool.NewWorkerpool(10)
		defer wp.Shutdown()
		_ = wp.Add()

		err := wp.AddJob("test job")
		assert.NoError(t, err)
	})

	t.Run("should return error when wp is closed", func(t *testing.T) {
		wp := workerpool.NewWorkerpool(10)
		wp.Shutdown()

		err := wp.AddJob("test job")
		assert.Error(t, err)
		assert.Equal(t, "workerpool is closed", err.Error())
	})
}

func TestDeleteWorker(t *testing.T) {
	t.Run("should delete worker successfully", func(t *testing.T) {
		wp := workerpool.NewWorkerpool(10)
		defer wp.Shutdown()
		_ = wp.Add()
		_ = wp.Add()

		err := wp.Delete()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(wp.GetWorkers()))
	})

	t.Run("should return error when no workers", func(t *testing.T) {
		wp := workerpool.NewWorkerpool(10)
		defer wp.Shutdown()

		err := wp.Delete()
		assert.Error(t, err)
		assert.Equal(t, "workerpool is empty", err.Error())
	})
}
