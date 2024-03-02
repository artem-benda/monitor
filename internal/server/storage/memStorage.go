package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/artem-benda/monitor/internal/model"
)

type memStorage struct {
	values         map[model.Metric]any
	rw             *sync.RWMutex
	filename       string
	writeImmediate bool
}

func NewStorage(saveIntervalSec int, filename string, restore bool) (Storage, error) {
	var metrics map[model.Metric]any = make(map[model.Metric]any)
	var savedMetrics []model.SaveableMetricValue
	if restore {
		if _, err := os.Stat(filename); !errors.Is(err, os.ErrNotExist) {
			var (
				bytes []byte
				err   error
			)
			if bytes, err = os.ReadFile(filename); err != nil {
				return nil, err
			}
			if err = json.Unmarshal(bytes, &savedMetrics); err != nil {
				return nil, err
			}
			for _, m := range savedMetrics {
				switch m.Kind {
				case model.GaugeKind:
					{
						metrics[model.Metric{Kind: m.Kind, Name: m.Name}] = m.Float64Value
					}
				case model.CounterKind:
					{
						metrics[model.Metric{Kind: m.Kind, Name: m.Name}] = m.Int64Value
					}
				}
			}
		}
	}

	s := memStorage{
		values:         metrics,
		rw:             &sync.RWMutex{},
		filename:       filename,
		writeImmediate: saveIntervalSec == 0,
	}

	if saveIntervalSec > 0 {
		ticker := time.NewTicker(time.Duration(saveIntervalSec) * time.Second)
		go func() {
			for range ticker.C {
				s.rw.Lock()
				s.saveToFile()
				s.rw.Unlock()
			}
		}()
	}
	return &s, nil
}

func (m memStorage) Get(key model.Metric) (any, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	value, ok := m.values[key]
	return value, ok
}

func (m *memStorage) Put(key model.Metric, value any) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.values[key] = value
	var err error
	if m.writeImmediate {
		err = m.saveToFile()
	}
	return err
}

func (m *memStorage) UpdateAndGetFunc(key model.Metric, fn func(prev any) any) (any, error) {
	m.rw.Lock()
	defer m.rw.Unlock()
	prev := m.values[key]
	next := fn(prev)
	m.values[key] = next
	var err error
	if m.writeImmediate {
		err = m.saveToFile()
	}
	return next, err
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

func (m memStorage) saveToFile() error {
	var saveableMetrics []model.SaveableMetricValue
	for k, v := range m.values {
		s, err := model.AsSaveableMetric(k, v)
		if err != nil {
			return err
		}
		saveableMetrics = append(saveableMetrics, s)
	}
	bytes, err := json.Marshal(saveableMetrics)
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile(m.filename, bytes, 0666); err != nil {
		panic(err)
	}
	return nil
}
