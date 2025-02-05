package email

import (
	"fmt"
	"go.uber.org/zap"
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
			s.logger.Info("Retrying email sending",
				zap.Int("attempt", attempt+1),
				zap.Int("maxAttempts", maxRetries),
			)
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
