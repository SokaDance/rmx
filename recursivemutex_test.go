package rmx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type counter struct {
	m RecursiveMutex
	r int
}

func (c *counter) add(n int) {
	c.m.Lock()
	c.r += n
	c.m.Unlock()
}

func (c *counter) augment(n int) {
	c.m.Lock()
	c.add(n)
	c.m.Unlock()
}

func TestRecursiveMutex_Lock(t *testing.T) {
	c := &counter{}
	wg := &sync.WaitGroup{}
	wg.Add(5)
	for i := range 5 {
		index := i
		go func() {
			for range 5 {
				c.augment(index)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, 50, c.r)
}

func TestRecursiveMutex_TryLock(t *testing.T) {
	m := &RecursiveMutex{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	require.True(t, m.TryLock())
	require.True(t, m.TryLock())
	go func() {
		assert.False(t, m.TryLock())
		wg.Done()
	}()
	wg.Wait()
	m.Unlock()
	m.Unlock()
}
