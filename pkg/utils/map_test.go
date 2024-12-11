package utils

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSyncMap(t *testing.T) {
	mm := &SyncMap[string, int]{}
	t.Run("LoadOrStore", func(t *testing.T) {
		{
			actual, loaded := mm.LoadOrStore("1", 1)
			assert.False(t, loaded)
			assert.Equal(t, 1, actual)
		}
		{
			actual, loaded := mm.LoadOrStore("1", 2)
			assert.True(t, loaded)
			assert.Equal(t, 1, actual)
		}
	})

	t.Run("Store", func(t *testing.T) {
		mm.Store("2", 2)
		assert.Equal(t, mm.Get("2"), 2)
	})

	t.Run("LoadOrStoreFunc", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)
		start := time.Now()
		go func() {
			defer wg.Done()
			actual, err := mm.LoadOrStoreFunc("3", func() (int, error) {
				time.Sleep(time.Millisecond * 100)
				return 3, nil
			})
			assert.NoError(t, err)
			assert.Equal(t, 3, actual)
		}()
		go func() {
			defer wg.Done()
			actual, err := mm.LoadOrStoreFunc("4", func() (int, error) {
				time.Sleep(time.Millisecond * 100)
				return 4, nil
			})
			assert.NoError(t, err)
			assert.Equal(t, 4, actual)
		}()
		wg.Wait()
		assert.Equal(t, time.Since(start).Milliseconds() < 110, true)
	})
}
