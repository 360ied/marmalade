package world

import (
	"sync"
)

type ConcurrentSlice struct {
	data []byte
	lock *sync.RWMutex
}

func NewConcurrentSlice(size int) *ConcurrentSlice {
	return &ConcurrentSlice{make([]byte, size), new(sync.RWMutex)}
}

func (c *ConcurrentSlice) Get(index int) byte {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data[index]
}

func (c *ConcurrentSlice) Set(index int, value byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[index] = value
}

func (c *ConcurrentSlice) Snapshot(b []byte) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	copy(b, c.data)
}

func (c *ConcurrentSlice) Len() int {
	return len(c.data)
}
