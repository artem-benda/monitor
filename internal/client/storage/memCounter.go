package storage

import "sync"

type memCounter struct {
	currentValue uint64
	rw           *sync.RWMutex
}

func NewCounter() Counter {
	return &memCounter{0, &sync.RWMutex{}}
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

var CounterStore Counter = &memCounter{}
