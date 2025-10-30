package initializer

import (
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/reporters"
	"wb_logistic_assistant/internal/scheduler"
	"wb_logistic_assistant/internal/services"
	"wb_logistic_assistant/internal/storage"
)

type AppDependencies struct {
	Config                *config.Config
	Storage               storage.Storage
	Services              *services.Container
	Scheduler             scheduler.Scheduler
	GeneralRoutesReporter reporters.Reporter
	ShipmentCloseReporter reporters.Reporter
}
