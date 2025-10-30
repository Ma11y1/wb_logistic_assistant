package reporters

import (
	"context"
	"time"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
)

type cacheWithTime[T interface{}] struct {
	value T
	time  time.Time
}

func retryAction(ctx context.Context, source string, attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		logger.Logf(logger.WARN, "Reporters.retryAction()", "failed action %s, attempt %d/%d: %v", source, i+1, attempts, err)

		delay *= time.Duration(1 << i)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "Reporters.retryAction()", "context cancelled while retrying")
		}
	}
	return errors.Wrapf(err, "Reporters.retryAction()", "all %d attempts failed for %s", attempts, source)
}
