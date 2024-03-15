package storage

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/artem-benda/monitor/internal/dto"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/mailru/easyjson"
)

type memStorage struct {
	values         map[model.MetricKey]model.MetricValue
	rw             *sync.RWMutex
	filename       string
	writeImmediate bool
}

func NewMemStorage(saveIntervalSec int, filename string, restore bool) (Storage, error) {
	var metrics = make(map[model.MetricKey]model.MetricValue)
	var savedMetrics dto.MetricsBatch

	if restore {
		if _, err := os.Stat(filename); !errors.Is(err, os.ErrNotExist) {
			var (
				bytes []byte
				err   error
			)
			if bytes, err = os.ReadFile(filename); err != nil {
				return nil, err
			}
			if err = easyjson.Unmarshal(bytes, &savedMetrics); err != nil {
				return nil, err
			}
			for _, m := range savedMetrics {
				switch m.MType {
				case model.GaugeKind:
					{
						metrics[model.MetricKey{Kind: m.MType, Name: m.ID}] = model.MetricValue{Gauge: *m.Value}
					}
				case model.CounterKind:
					{
						metrics[model.MetricKey{Kind: m.MType, Name: m.ID}] = model.MetricValue{Counter: *m.Delta}
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

func (m memStorage) Get(ctx context.Context, key model.MetricKey) (*model.MetricValue, bool, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	value, ok := m.values[key]
	return &value, ok, nil
}

func (m *memStorage) UpsertGauge(ctx context.Context, key model.MetricKey, value model.MetricValue) error {
	if key.Kind != model.GaugeKind {
		return errInvaligArgument
	}

	m.rw.Lock()
	defer m.rw.Unlock()
	m.values[key] = value
	var err error
	if m.writeImmediate {
		err = m.saveToFile()
	}
	return err
}

func (m *memStorage) UpsertCounterAndGet(ctx context.Context, key model.MetricKey, incCounter int64) (int64, error) {
	if key.Kind != model.CounterKind {
		return 0, errInvaligArgument
	}

	m.rw.Lock()
	defer m.rw.Unlock()
	prev := m.values[key]
	nextCounter := prev.Counter + incCounter
	m.values[key] = model.MetricValue{Counter: nextCounter}
	var err error
	if m.writeImmediate {
		err = m.saveToFile()
	}
	return nextCounter, err
}

func (m memStorage) GetAll(ctx context.Context) (map[model.MetricKey]model.MetricValue, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	copy := make(map[model.MetricKey]model.MetricValue)
	for k, v := range m.values {
		copy[k] = v
	}

	return copy, nil
}

func (m memStorage) UpsertBatch(ctx context.Context, metrics []model.MetricKeyWithValue) error {
	for _, v := range metrics {
		switch v.Kind {
		case model.GaugeKind:
			{
				err := m.UpsertGauge(ctx, model.MetricKey{Kind: v.Kind, Name: v.Name}, model.MetricValue{Gauge: v.Gauge, Counter: v.Counter})

				if err != nil {
					return err
				}
			}
		case model.CounterKind:
			{
				_, err := m.UpsertCounterAndGet(ctx, model.MetricKey{Kind: v.Kind, Name: v.Name}, v.Counter)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m memStorage) saveToFile() error {
	dtos := make(dto.MetricsBatch, 0)

	for k, v := range m.values {
		s := model.AsDto(k, v)
		dtos = append(dtos, s)
	}

	bytes, err := easyjson.Marshal(dtos)
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile(m.filename, bytes, 0666); err != nil {
		panic(err)
	}
	return nil
}
