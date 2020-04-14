package queue

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	m, err := NewManager(&ManagerOpts{MaxWorker: 1, MaxJob: 1})
	assert.Nil(t, err)
	go m.Start()
	wg := &sync.WaitGroup{}
	var a, b, c int
	wg.Add(3)

	start := time.Now()

	m.Dispatch(JobFunc(func() error {
		a = 1
		wg.Done()
		assert.LessOrEqual(t, int64(0), int64(time.Since(start)))
		return nil
	}))

	m.DispatchAfter(JobFunc(func() error {
		b = 2
		wg.Done()
		assert.LessOrEqual(t, int64(time.Second*2), int64(time.Since(start)))
		return nil
	}), time.Second*2)

	m.DispatchAt(JobFunc(func() error {
		c = 3
		wg.Done()
		assert.LessOrEqual(t, int64(time.Second*3), int64(time.Since(start)))
		return nil
	}), time.Now().Add(time.Second*3))

	wg.Wait()
	assert.Equal(t, 1, a)
	assert.Equal(t, 2, b)
	assert.Equal(t, 3, c)
}
