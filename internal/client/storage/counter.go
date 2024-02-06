package storage

type Counter interface {
	IncrementAndGet() uint64
	Get() uint64
	Reset()
}
