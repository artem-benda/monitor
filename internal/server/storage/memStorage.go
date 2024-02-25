package storage

import (
	"sync"

	"github.com/artem-benda/monitor/internal/model"
)

type memStorage struct {
	values map[model.Metric]any
	rw     *sync.RWMutex
}

func NewStorage() Storage {
	return &memStorage{make(map[model.Metric]any), &sync.RWMutex{}}
}

func (m memStorage) Get(key model.Metric) (any, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	value, ok := m.values[key]
	return value, ok
}

func (m *memStorage) Put(key model.Metric, value any) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.values[key] = value
}

func (m *memStorage) UpdateAndGetFunc(key model.Metric, fn func(prev any) any) any {
	m.rw.Lock()
	defer m.rw.Unlock()
	prev := m.values[key]
	next := fn(prev)
	m.values[key] = next
	return next
}

func (m memStorage) GetAll() map[model.Metric]any {
	m.rw.RLock()
	defer m.rw.RUnlock()
	copy := make(map[model.Metric]any)
	for k, v := range m.values {
		copy[k] = v
	}

	return copy
}

var Store = NewStorage()
