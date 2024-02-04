package service

import (
	"testing"

	"github.com/artem-benda/monitor/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestReadMetrics(t *testing.T) {
	counter := storage.NewCounter()
	result := ReadMetrics(counter)
	assert.Equal(t, 29, len(result))
}
