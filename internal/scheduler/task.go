package scheduler

import (
	"context"
	"sync/atomic"
	"time"
	"wb_logistic_assistant/internal/errors"
)

type TaskConfig struct {
	RetryTaskLimit        int           // how many times to repeat; if 0, the default value of the scheduler is used
	Timeout               time.Duration // maximum execution time; if 0, no timeouts
	IsWaitForPrevious     bool          // do  need to wait for the completion of the previous task
	IsIntervalAfterFinish bool          // "wait for interval after task completion" mode
}

type Task interface {
	ID() uint64
	Name() string
	Execute(context.Context) error
}

var globalTaskID uint64

type CallbackTask struct {
	id       uint64
	name     string
	callback func(context.Context) error
}

func NewCallbackTask(name string, callback func(context.Context) error) *CallbackTask {
	return &CallbackTask{
		id:       atomic.AddUint64(&globalTaskID, 1),
		name:     name,
		callback: callback,
	}
}

func (t *CallbackTask) ID() uint64 {
	return t.id
}

func (t *CallbackTask) Name() string {
	if t.name == "" {
		return "Callback"
	}
	return t.name
}

func (t *CallbackTask) Execute(ctx context.Context) error {
	if t.callback == nil {
		return errors.New("CallbackTask.Execute()", "callback func is nil")
	}
	return t.callback(ctx)
}
