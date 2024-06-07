package storage

import "sync"

type memCounter struct {
	rw           *sync.RWMutex
	currentValue uint64
}

func NewCounter() Counter {
	return &memCounter{&sync.RWMutex{}, 0}
}

func (c *memCounter) IncrementAndGet() uint64 {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.currentValue++
	return c.currentValue
}

func (c memCounter) Get() uint64 {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.currentValue
}

func (c *memCounter) Reset() {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.currentValue = 0
}

var CounterStore Counter = &memCounter{&sync.RWMutex{}, 0}
