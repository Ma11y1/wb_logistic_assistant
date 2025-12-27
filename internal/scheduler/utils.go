package scheduler

import (
	"context"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
)

func cfgOrDefault(cfg []TaskConfig) TaskConfig {
	if len(cfg) > 0 {
		return cfg[0]
	}
	return TaskConfig{}
}

func safeRunTask(task Task, ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Newf("BaseScheduler.safeRunTask()", "panic in task: %v", r)
		}
	}()
	err = task.Execute(ctx)
	return err
}

func safeRunGoroutine(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Logf(logger.ERROR, "BaseScheduler.safeRunGoroutine()", "panic in goroutine: %v", r)
			}
		}()
		f()
	}()
}
