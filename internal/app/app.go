package app

import (
	"context"
	"os"
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/initializer"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/prompters"
	"wb_logistic_assistant/internal/reporters"
	"wb_logistic_assistant/internal/scheduler"
	"wb_logistic_assistant/internal/services"
	"wb_logistic_assistant/internal/storage"
)

type App struct {
	config                *config.Config
	storage               storage.Storage
	services              *services.Container
	scheduler             scheduler.Scheduler
	generalRoutesReporter reporters.Reporter
	shipmentCloseReporter reporters.Reporter
	isStarted             bool
}

func NewApp(config *config.Config) *App {
	return &App{config: config}
}

func (a *App) Init() error {
	logger.Log(logger.INFO, "App.Init()", "Start init app")
	cfg := a.config

	storage, err := storage.NewFileStorage(cfg)
	if err != nil {
		return errors.Wrap(err, "App.Init()", "Failed to create storage")
	}

	storagePassword := os.Getenv(cfg.Storage().Env())
	if storagePassword == "" {
		return errors.New("App.Init()", "Storage key is missing from the system")
	}
	storage.SetEncrypt([]byte(storagePassword))

	err = storage.Load(cfg.Storage().Path())
	if err != nil {
		return errors.Wrap(err, "App.Init()", "Failed to load storage")
	}

	init := initializer.NewInitializer(cfg, storage, &prompters.CLIInitAppPrompter{})
	dependencies, err := init.Init()
	if err != nil {
		return errors.Wrap(err, "App.Init()", "Failed to init app dependencies")
	}

	a.storage = storage
	a.services = dependencies.Services
	a.scheduler = dependencies.Scheduler
	a.generalRoutesReporter = dependencies.GeneralRoutesReporter
	a.shipmentCloseReporter = dependencies.ShipmentCloseReporter

	logger.Log(logger.INFO, "App.Init()", "Init app successfully")
	return nil
}

func (a *App) Start() error {
	if a.isStarted {
		return errors.New("App.Start()", "Application already started")
	}

	a.isStarted = true
	logger.Log(logger.INFO, "App.Start()", "Start application")

	//if a.config.Reports().GeneralRoutes().IsEnabled() {
	//	logger.Log(logger.INFO, "App.Start()", "Start schedule periodic for general routes reporter")
	//	a.scheduler.SchedulePeriodic(scheduler.NewCallbackTask(
	//		"general_routes_report",
	//		func(ctx context.Context) error {
	//			select {
	//			case <-ctx.Done():
	//				logger.Log(logger.ERROR, "App.Start()", "Cancelling 'general_routes_report' callback task")
	//				return ctx.Err()
	//			default:
	//				err := a.generalRoutesReporter.Run(ctx)
	//				if err != nil {
	//					logger.Log(logger.ERROR, "App.Start()", "Failed to run general routes report")
	//					return err
	//				}
	//			}
	//			return nil
	//		}),
	//		a.config.Reports().GeneralRoutes().PollingInterval(),
	//		scheduler.TaskConfig{
	//			RetryTaskLimit:        a.config.Reports().GeneralRoutes().ErrRetryTaskLimit(),
	//			Timeout:               a.config.Reports().GeneralRoutes().TaskTimeout(),
	//			IsWaitForPrevious:     true,
	//			IsIntervalAfterFinish: true,
	//		})
	//}

	if a.config.Reports().ShipmentClose().IsEnabled() {
		logger.Log(logger.INFO, "App.Start()", "Start schedule periodic for shipment close reporter")
		a.scheduler.SchedulePeriodic(scheduler.NewCallbackTask(
			"shipment_close_report",
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					logger.Log(logger.ERROR, "App.Start()", "Cancelling 'shipment_close_report' callback task")
					return ctx.Err()
				default:
					err := a.shipmentCloseReporter.Run(ctx)
					if err != nil {
						logger.Log(logger.ERROR, "App.Start()", "Failed to run shipment close report")
						return err
					}
				}
				return nil
			}),
			a.config.Reports().ShipmentClose().PollingInterval(),
			scheduler.TaskConfig{
				RetryTaskLimit:        a.config.Reports().ShipmentClose().ErrRetryTaskLimit(),
				Timeout:               a.config.Reports().ShipmentClose().TaskTimeout(),
				IsWaitForPrevious:     true,
				IsIntervalAfterFinish: true,
			})
	}

	return nil
}

func (a *App) Stop() {
	logger.Log(logger.INFO, "App.Stop()", "Stop application")
	a.scheduler.Cancel()
	a.isStarted = false

	err := a.storage.Save(a.config.Storage().Path())
	if err != nil {
		logger.Logf(logger.ERROR, "App.Stop()", "Failed to save storage: %v", err)
	}

	a.storage.SetEncrypt([]byte(""))
	a.storage.Clear()
}
