// RetryController повторяет операцию до тех пор, пока она не завершится успешно, но не более 3 раз
// Задержка увеличивается от 1 секунды на 2 секунды каждый повтор:
//
//	delay = delay + nextDelaySecondsDelta.
package retry

import (
	"errors"
	"time"
)

// RetryController хранит список ошибок, при которых требуется повторно выполнить операцию.
type RetryController struct {
	retriableError error
}

const (
	maxRetries            = 3
	firstDelaySeconds     = 1
	nextDelaySecondsDelta = 2
)

// NewRetryController создает новый контроллер для списка ошибок,
// при которых необходимо выполнять операцию повторно.
func NewRetryController(retriableErrors ...error) RetryController {
	return RetryController{retriableError: errors.Join(retriableErrors...)}
}

// Run выполняет операцию fn. В случае, если fn вернет retriableError -
// будет выполнена повторная попытка вызова fn.
func (r RetryController) Run(fn func() error) (err error) {
	// Try
	err = fn()
	if err == nil {
		return nil
	}

	if !errors.Is(err, r.retriableError) {
		return err
	}

	// Retries
	delay := firstDelaySeconds
	for i := 0; i < maxRetries-1; i++ {
		time.Sleep(time.Duration(delay) * time.Second)

		err = fn()
		if err == nil {
			return nil
		}

		if !errors.Is(err, r.retriableError) {
			return err
		}

		delay += nextDelaySecondsDelta
	}
	return
}
