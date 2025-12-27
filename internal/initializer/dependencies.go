package initializer

import (
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/reporters"
	"wb_logistic_assistant/internal/scheduler"
	"wb_logistic_assistant/internal/services"
	"wb_logistic_assistant/internal/storage"
)

type AppDependencies struct {
	Config                           *config.Config
	Storage                          storage.Storage
	Services                         *services.Container
	Scheduler                        scheduler.Scheduler
	SchedulerGeneralRoutesTaskConfig *scheduler.TaskConfig
	SchedulerShipmentCloseTaskConfig *scheduler.TaskConfig
	SchedulerFinanceRoutesTaskConfig *scheduler.TaskConfig
	SchedulerFinanceDailyTaskConfig  *scheduler.TaskConfig
	GeneralRoutesReporter            reporters.Reporter
	ShipmentCloseReporter            reporters.Reporter
	FinanceRoutesReporter            reporters.Reporter
	FinanceDailyReporter             reporters.Reporter
}
