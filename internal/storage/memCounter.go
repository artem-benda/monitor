package storage

type memCounter struct {
	currentValue uint64
}

func (c *memCounter) IncrementAndGet() uint64 {
	c.currentValue++
	return c.currentValue
}

func (c memCounter) Get() uint64 {
	return c.currentValue
}

func (c *memCounter) Reset() {
	c.currentValue = 0
}

var CounterStore Counter = &memCounter{}
