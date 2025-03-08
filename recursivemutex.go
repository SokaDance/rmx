package rmx

import (
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

type RecursiveMutex struct {
	mutex          sync.Mutex
	id             atomic.Int64
	recursionCount atomic.Int64
}

func (m *RecursiveMutex) Lock() {
	id := goid.Get()
	if !m.tryRecursiveLock(id) {
		m.mutex.Lock()
		m.id.Store(id)
		m.recursionCount.Store(1)
	}
}

func (m *RecursiveMutex) TryLock() bool {
	id := goid.Get()
	return m.tryRecursiveLock(id) || m.tryBasicLock(id)
}

func (m *RecursiveMutex) Unlock() {
	if m.recursionCount.Add(-1) == 0 {
		m.id.Store(0)
		m.mutex.Unlock()
	}
}

func (m *RecursiveMutex) tryRecursiveLock(id int64) bool {
	if m.id.Load() == id {
		m.recursionCount.Add(1)
		return true
	}
	return false
}

func (m *RecursiveMutex) tryBasicLock(id int64) bool {
	if m.mutex.TryLock() {
		m.id.Store(id)
		m.recursionCount.Store(1)
		return true
	}
	return false
}
