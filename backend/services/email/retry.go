package email

import (
	"fmt"
	"log"
	"time"
)

const (
	maxRetries    = 3
	retryInterval = 5 * time.Second
)

// RetryableError represents an error that can be retried
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

// withRetry wraps a function with retry mechanism
func (s *Service) withRetry(operation func() error) error {
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryInterval)
			log.Printf("Retrying email sending: attempt %d/%d", attempt+1, maxRetries)
		}

		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err
		if _, ok := err.(*RetryableError); !ok {
			return err
		}
	}

	return fmt.Errorf("failed after %d attempts, last error: %v", maxRetries, lastErr)
}