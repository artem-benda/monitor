package storage

import (
	"context"
	"database/sql"
	"errors"
	"sort"

	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/retry"
	srverror "github.com/artem-benda/monitor/internal/server/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type dbStorage struct {
	dbpool *pgxpool.Pool
	c      retry.RetryController
}

func NewDBStorage(dbpool *pgxpool.Pool, r retry.RetryController) Storage {
	return dbStorage{dbpool: dbpool, c: r}
}

func (s dbStorage) Get(ctx context.Context, key model.MetricKey) (m model.MetricValue, ok bool, err error) {
	err = s.c.Run(func() error {
		m, ok, err = s.get(ctx, key)
		return mapError(err)
	})
	return
}

func (s dbStorage) UpsertGauge(ctx context.Context, key model.MetricKey, value model.MetricValue) (err error) {
	err = s.c.Run(func() error {
		err = s.upsertGauge(ctx, key, value)
		return mapError(err)
	})
	return
}

func (s dbStorage) UpsertCounterAndGet(ctx context.Context, key model.MetricKey, incCounter int64) (cnt int64, err error) {
	err = s.c.Run(func() error {
		cnt, err = s.upsertCounterAndGet(ctx, key, incCounter)
		return mapError(err)
	})
	return
}

func (s dbStorage) GetAll(ctx context.Context) (m map[model.MetricKey]model.MetricValue, err error) {
	err = s.c.Run(func() error {
		m, err = s.getAll(ctx)
		return mapError(err)
	})
	return
}

func (s dbStorage) UpsertBatch(ctx context.Context, metrics []model.MetricKeyWithValue) (err error) {
	err = s.c.Run(func() error {
		err = s.upsertBatch(ctx, metrics)
		return mapError(err)
	})
	return
}

func (s dbStorage) get(ctx context.Context, key model.MetricKey) (model.MetricValue, bool, error) {
	var (
		mtype   string
		mname   string
		gauge   sql.NullFloat64
		counter sql.NullInt64
	)
	err := s.dbpool.QueryRow(
		ctx,
		"SELECT mtype, mname, gauge, counter FROM metrics WHERE mtype = $1 AND mname = $2",
		key.Kind,
		key.Name,
	).Scan(&mtype, &mname, &gauge, &counter)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		{
			logger.Log.Debug("metric not found")
			return model.MetricValue{}, false, nil
		}
	case err != nil:
		{
			logger.Log.Debug("get metric error", zap.String("Kind", key.Kind), zap.String("Name", key.Name), zap.Error(err))
			return model.MetricValue{}, false, err
		}
	case key.Kind == model.GaugeKind && gauge.Valid:
		{
			return model.MetricValue{Gauge: gauge.Float64}, true, nil
		}
	case key.Kind == model.CounterKind && counter.Valid:
		{
			return model.MetricValue{Counter: counter.Int64}, true, err
		}
	}
	return model.MetricValue{}, false, errInvalidData
}

func (s dbStorage) upsertGauge(ctx context.Context, key model.MetricKey, value model.MetricValue) error {
	if key.Kind != model.GaugeKind {
		return errInvaligArgument
	}

	gauge := sql.NullFloat64{Float64: value.Gauge}

	upsertMetricQuery := `INSERT INTO metrics(mtype, mname, gauge) VALUES ($1, $2, $3) 
		ON CONFLICT (mtype, mname) DO UPDATE SET gauge = EXCLUDED.gauge`

	_, err := s.dbpool.Exec(
		ctx,
		upsertMetricQuery,
		key.Kind,
		key.Name,
		gauge,
	)

	return err
}

func (s dbStorage) upsertCounterAndGet(ctx context.Context, key model.MetricKey, incCounter int64) (int64, error) {
	if key.Kind != model.CounterKind {
		return 0, errInvaligArgument
	}

	// Сюда получим актуальное значение счетчика после его обновления
	var counter sql.NullInt64

	upsertMetricQuery := `INSERT INTO metrics(mtype, mname, counter) VALUES ($1, $2, $3) 
		ON CONFLICT (mtype, mname) DO UPDATE SET counter = COALESCE(metrics.counter, 0) + EXCLUDED.counter 
		RETURNING counter`
	err := s.dbpool.QueryRow(
		ctx,
		upsertMetricQuery,
		key.Kind,
		key.Name,
		sql.NullInt64{Int64: incCounter, Valid: true},
	).Scan(&counter)

	if err != nil {
		return int64(0), err
	}

	if !counter.Valid {
		return int64(0), errNullCounterValue
	}

	return counter.Int64, err
}

func (s dbStorage) getAll(ctx context.Context) (map[model.MetricKey]model.MetricValue, error) {
	rows, err := s.dbpool.Query(ctx, "SELECT mtype, mname, gauge, counter FROM metrics ORDER BY mtype, mname")
	if err != nil {
		return nil, err
	}

	// Результирующая мапа
	m := make(map[model.MetricKey]model.MetricValue)

	// Переменные для чтения строк
	var (
		mtype   string
		mname   string
		gauge   sql.NullFloat64
		counter sql.NullInt64
	)

	for rows.Next() {
		rows.Scan(&mtype, &mname, &gauge, &counter)

		switch {
		case mtype == model.CounterKind && counter.Valid:
			{
				m[model.MetricKey{Kind: mtype, Name: mname}] = model.MetricValue{Counter: counter.Int64}
			}
		case mtype == model.GaugeKind && gauge.Valid:
			{
				m[model.MetricKey{Kind: mtype, Name: mname}] = model.MetricValue{Gauge: gauge.Float64}
			}
		default:
			{
				return m, errInvalidData
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func (s dbStorage) upsertBatch(ctx context.Context, metrics []model.MetricKeyWithValue) error {
	// Сортировка для исключения deadlocks при параллельном обновлении пачек метрик
	sort.Slice(metrics, func(i, j int) bool {
		if metrics[i].Name != metrics[j].Name {
			return metrics[i].Name < metrics[j].Name
		} else {
			return metrics[i].Kind < metrics[j].Kind
		}
	})

	tx, err := s.dbpool.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(
		ctx,
		"upsert-metrics",
		`INSERT INTO metrics(mtype, mname, gauge, counter) VALUES ($1, $2, $3, $4) 
		ON CONFLICT (mtype, mname) DO UPDATE
		SET counter = COALESCE(metrics.counter, 0) + EXCLUDED.counter, gauge = EXCLUDED.gauge`,
	)

	if err != nil {
		return err
	}

	for _, m := range metrics {
		_, err := tx.Exec(ctx, "upsert-metrics", m.Kind, m.Name, m.Gauge, m.Counter)

		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ConnectionException {
		return srverror.ErrStorageConnection{Err: pgErr}
	} else {
		return err
	}
}
