package storage

import "github.com/artem-benda/monitor/internal/model"

type memStorage struct {
	values map[model.Metric]any
}

func NewStorage() Storage {
	return &memStorage{make(map[model.Metric]any)}
}

func (m memStorage) Get(key model.Metric) (any, bool) {
	value, ok := m.values[key]
	return value, ok
}

func (m *memStorage) Put(key model.Metric, value any) {
	m.values[key] = value
}

func (m *memStorage) UpdateFunc(key model.Metric, fn func(prev any) any) {
	prev := m.values[key]
	m.values[key] = fn(prev)
}

func (m memStorage) GetAll() map[model.Metric]any {
	return m.values
}

var Store = NewStorage()
