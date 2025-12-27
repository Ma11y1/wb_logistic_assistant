package config

import "runtime"

type Internal struct {
	schedulerMaxWorkers     int
	schedulerRetryTaskLimit int
}

func newInternal() *Internal {
	return &Internal{
		schedulerMaxWorkers:     runtime.NumCPU(),
		schedulerRetryTaskLimit: 3,
	}
}

func (i *Internal) SchedulerMaxWorkers() int { return i.schedulerMaxWorkers }

func (i *Internal) SchedulerRetryTaskLimit() int { return i.schedulerRetryTaskLimit }
