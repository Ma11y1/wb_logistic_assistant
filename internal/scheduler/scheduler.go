package scheduler

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
)

type Scheduler interface {
	ScheduleNow(task Task, cfg ...TaskConfig)
	ScheduleAsync(task Task, cfg ...TaskConfig)
	ScheduleAfter(task Task, delay time.Duration, cfg ...TaskConfig)
	SchedulePeriodic(task Task, interval time.Duration, cfg ...TaskConfig)
	Reset()
}

type BaseScheduler struct {
	ctx                   context.Context
	cancelFunc            context.CancelFunc
	isCanceled            atomic.Bool
	wg                    sync.WaitGroup
	workerPool            chan struct{}
	defaultRetryTaskLimit int
}

func NewBaseScheduler(maxWorkers, defaultRetryTaskLimit int) *BaseScheduler {
	if defaultRetryTaskLimit <= 0 {
		defaultRetryTaskLimit = 3
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseScheduler{
		ctx:                   ctx,
		cancelFunc:            cancel,
		workerPool:            make(chan struct{}, maxWorkers),
		defaultRetryTaskLimit: defaultRetryTaskLimit,
	}
}

func (s *BaseScheduler) ScheduleNow(task Task, cfg ...TaskConfig) {
	s.runTaskAsync(task, cfgOrDefault(cfg), "Now")
}

func (s *BaseScheduler) ScheduleAsync(task Task, cfg ...TaskConfig) {
	safeRunGoroutine(func() {
		s.runTaskAsync(task, cfgOrDefault(cfg), "Async")
	})
}

func (s *BaseScheduler) ScheduleAfter(task Task, delay time.Duration, cfg ...TaskConfig) {
	safeRunGoroutine(func() {
		select {
		case <-time.After(delay):
			s.runTaskAsync(task, cfgOrDefault(cfg), "After")
		case <-s.ctx.Done():
		}
	})
}

func (s *BaseScheduler) SchedulePeriodic(task Task, interval time.Duration, cfg ...TaskConfig) {
	config := cfgOrDefault(cfg)
	if config.IsIntervalAfterFinish {
		s.schedulePeriodicSequential(task, interval, config)
	} else {
		s.schedulePeriodicTicker(task, interval, config)
	}
}

// fixed interval mode
func (s *BaseScheduler) schedulePeriodicTicker(task Task, interval time.Duration, cfg TaskConfig) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var mtx sync.Mutex
		var isRunning bool

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				if cfg.IsWaitForPrevious {
					mtx.Lock()
					if isRunning {
						mtx.Unlock()
						continue
					}
					isRunning = true
					mtx.Unlock()

					s.wg.Add(1)
					go func() {
						defer s.wg.Done()
						s.runTask(task, cfg, "Periodic")
						mtx.Lock()
						isRunning = false
						mtx.Unlock()
					}()
				} else {
					s.runTaskAsync(task, cfg, "Periodic")
				}
			}
		}
	}()
}

// "interval after completion" mode
func (s *BaseScheduler) schedulePeriodicSequential(task Task, interval time.Duration, cfg TaskConfig) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				s.runTask(task, cfg, "PeriodicSeq")

				timer := time.NewTimer(interval)
				select {
				case <-timer.C:
				case <-s.ctx.Done():
					timer.Stop()
					return
				}
				timer.Stop()
			}
		}
	}()
}

func (s *BaseScheduler) Reset() {
	if s.isCanceled.CompareAndSwap(false, true) {
		s.cancelFunc()

		logger.Log(logger.INFO, "BaseScheduler.Reset()", "waiting for all tasks to complete...")
		s.wg.Wait()
		logger.Log(logger.INFO, "BaseScheduler.Reset()", "all tasks completed")

		newCtx, newCancel := context.WithCancel(context.Background())
		s.ctx = newCtx
		s.cancelFunc = newCancel

		s.isCanceled.Store(false)

		logger.Log(logger.INFO, "BaseScheduler.Reset()", "scheduler context reset")
	}
}

func (s *BaseScheduler) runTask(task Task, cfg TaskConfig, source string) {
	s.wg.Add(1)
	defer s.wg.Done()

	select {
	case s.workerPool <- struct{}{}:
		defer func() { <-s.workerPool }()
		s.executeWithRetry(task, cfg, source)
	case <-s.ctx.Done():
		logger.Logf(logger.ERROR, "BaseScheduler.runTask()", "%s task %d '%s' skipped due to shutdown", source, task.ID(), task.Name())
	}
}

func (s *BaseScheduler) runTaskAsync(task Task, cfg TaskConfig, source string) {
	s.wg.Add(1)
	safeRunGoroutine(func() {
		defer s.wg.Done()

		select {
		case s.workerPool <- struct{}{}:
			defer func() { <-s.workerPool }()
			s.executeWithRetry(task, cfg, source)
		case <-s.ctx.Done():
			logger.Logf(logger.ERROR, "BaseScheduler.runTaskAsync()", "%s task '%s' skipped due to shutdown, because scheduler has been cancelled", source, task.Name())
		}
	})
}

func (s *BaseScheduler) executeWithRetry(task Task, cfg TaskConfig, source string) {
	retries := cfg.RetryTaskLimit
	if retries <= 0 {
		retries = s.defaultRetryTaskLimit
	}

	for attempt := 1; attempt <= retries; attempt++ {
		if s.ctx.Err() != nil {
			logger.Logf(logger.ERROR, "BaseScheduler", "%s task %d '%s' cancelled before start", source, task.ID(), task.Name())
			return
		}

		start := time.Now()
		logger.Logf(logger.INFO, "BaseScheduler", "%s task %d '%s' started (attempt %d)", source, task.ID(), task.Name(), attempt)

		runCtx := s.ctx
		var cancel context.CancelFunc
		if cfg.Timeout > 0 {
			runCtx, cancel = context.WithTimeout(s.ctx, cfg.Timeout)
		}

		err := safeRunTask(task, runCtx)

		if cancel != nil {
			cancel()
		}

		if err == nil {
			logger.Logf(logger.INFO, "BaseScheduler", "%s task %d '%s' finished in %v", source, task.ID(), task.Name(), time.Since(start))
			return
		}

		if errors.Is(err, context.Canceled) {
			logger.Logf(logger.INFO, "BaseScheduler", "%s task %d '%s' cancelled", source, task.ID(), task.Name())
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			logger.Logf(logger.ERROR, "BaseScheduler", "%s task %d '%s' timed out", source, task.ID(), task.Name())
			return
		}

		logger.Logf(logger.ERROR, "BaseScheduler", "%s task %d '%s' failed attempt %d: %v", source, task.ID(), task.Name(), attempt, err)

		if attempt == retries {
			logger.Logf(logger.ERROR, "BaseScheduler", "%s task %d '%s' permanently failed after %d attempts", source, task.ID(), task.Name(), attempt)
			return
		}

		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		logger.Logf(logger.INFO, "BaseScheduler", "%s task %d '%s' retrying in %v...", source, task.ID(), task.Name(), backoff)

		select {
		case <-time.After(backoff):
		case <-s.ctx.Done():
			logger.Logf(logger.ERROR, "BaseScheduler", "%s task %d '%s' retry aborted by cancellation", source, task.ID(), task.Name())
			return
		}
	}
}
