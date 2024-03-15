package retry

import (
	"errors"
	"time"
)

type RetryController struct {
	retriableError error
}

const (
	maxRetries            = 3
	firstDelaySeconds     = 1
	nextDelaySecondsDelta = 2
)

func NewRetryController(retriableErrors ...error) RetryController {
	return RetryController{retriableError: errors.Join(retriableErrors...)}
}

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
