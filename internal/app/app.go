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
	config                           *config.Config
	initializer                      *initializer.Initializer
	storage                          storage.Storage
	services                         *services.Container
	scheduler                        scheduler.Scheduler
	schedulerGeneralRoutesTaskConfig *scheduler.TaskConfig
	schedulerShipmentCloseTaskConfig *scheduler.TaskConfig
	schedulerFinanceRoutesTaskConfig *scheduler.TaskConfig
	schedulerFinanceDailyTaskConfig  *scheduler.TaskConfig
	generalRoutesReporter            reporters.Reporter
	shipmentCloseReporter            reporters.Reporter
	financeRoutesReporter            reporters.Reporter
	financeDailyReporter             reporters.Reporter
	isStarted                        bool
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

	a.initializer = initializer.NewInitializer(cfg, storage, &prompters.CLIInitAppPrompter{})
	dependencies, err := a.initializer.Init()
	if err != nil {
		return errors.Wrap(err, "App.Init()", "Failed to init app dependencies")
	}

	a.storage = storage
	a.services = dependencies.Services
	a.scheduler = dependencies.Scheduler
	a.schedulerGeneralRoutesTaskConfig = dependencies.SchedulerGeneralRoutesTaskConfig
	a.schedulerShipmentCloseTaskConfig = dependencies.SchedulerShipmentCloseTaskConfig
	a.schedulerFinanceRoutesTaskConfig = dependencies.SchedulerFinanceRoutesTaskConfig
	a.schedulerFinanceDailyTaskConfig = dependencies.SchedulerFinanceDailyTaskConfig
	a.generalRoutesReporter = dependencies.GeneralRoutesReporter
	a.shipmentCloseReporter = dependencies.ShipmentCloseReporter
	a.financeRoutesReporter = dependencies.FinanceRoutesReporter
	a.financeDailyReporter = dependencies.FinanceDailyReporter

	logger.Log(logger.INFO, "App.Init()", "Init app successfully")
	return nil
}

func (a *App) Start() error {
	if a.isStarted {
		return errors.New("App.Start()", "Application already started")
	}

	a.isStarted = true
	logger.Log(logger.INFO, "App.Start()", "Start application")

	a.runTasks()

	return nil
}

func (a *App) Stop() {
	logger.Log(logger.INFO, "App.Stop()", "Stop application")
	a.scheduler.Reset()
	a.isStarted = false

	err := a.storage.Save(a.config.Storage().Path())
	if err != nil {
		logger.Logf(logger.ERROR, "App.Stop()", "Failed to save storage: %v", err)
	}

	a.storage.SetEncrypt([]byte("")) // clear password
	a.storage.Clear()
}

func (a *App) Pause() {
	logger.Log(logger.INFO, "App.Pause()", "Pause application")
	a.scheduler.Reset()
	a.isStarted = false
}

func (a *App) runTasks() {
	if a.config.Reports().GeneralRoutes().IsEnabled() {
		logger.Log(logger.INFO, "App.runTasks()", "Start schedule periodic for \"general routes\" reporter")
		a.scheduler.SchedulePeriodic(
			scheduler.NewCallbackTask("general_routes_report", a.generalRoutesHandler),
			a.config.Reports().GeneralRoutes().PollingInterval(),
			*a.schedulerGeneralRoutesTaskConfig,
		)
	}

	if a.config.Reports().ShipmentClose().IsEnabled() {
		logger.Log(logger.INFO, "App.runTasks()", "Start schedule periodic for \"shipment close\" reporter")
		a.scheduler.SchedulePeriodic(
			scheduler.NewCallbackTask("shipment_close_report", a.shipmentCloseHandler),
			a.config.Reports().ShipmentClose().PollingInterval(),
			*a.schedulerShipmentCloseTaskConfig,
		)
	}

	if a.config.Reports().FinanceRoutes().IsEnabled() {
		logger.Log(logger.INFO, "App.runTasks()", "Start schedule periodic for \"finance routes\" reporter")
		a.scheduler.SchedulePeriodic(
			scheduler.NewCallbackTask("finance_routes_report", a.financeRoutesHandler),
			a.config.Reports().FinanceRoutes().PollingInterval(),
			*a.schedulerFinanceRoutesTaskConfig,
		)
	}

	if a.config.Reports().FinanceDaily().IsEnabled() {
		logger.Log(logger.INFO, "App.runTasks()", "Start schedule periodic for \"finance daily\" reporter")
		a.scheduler.SchedulePeriodic(
			scheduler.NewCallbackTask("finance_daily_report", a.financeDailyHandler),
			a.config.Reports().FinanceDaily().PollingInterval(),
			*a.schedulerFinanceDailyTaskConfig,
		)
	}
}

func (a *App) generalRoutesHandler(ctx context.Context) error {
	select {
	case <-ctx.Done():
		logger.Log(logger.ERROR, "App.generalRoutesHandler()", "Cancelling 'general_routes_report' callback task")
		return ctx.Err()
	default:
		err := a.generalRoutesReporter.Run(ctx)
		if err != nil {
			logger.Log(logger.ERROR, "App.generalRoutesHandler()", "Failed to run general routes report")
			a.checkAuthWBLogistic()
			return err
		}
	}
	return nil
}

func (a *App) shipmentCloseHandler(ctx context.Context) error {
	select {
	case <-ctx.Done():
		logger.Log(logger.ERROR, "App.shipmentCloseHandler()", "Cancelling 'shipment_close_report' callback task")
		return ctx.Err()
	default:
		err := a.shipmentCloseReporter.Run(ctx)
		if err != nil {
			logger.Log(logger.ERROR, "App.shipmentCloseHandler()", "Failed to run shipment close report")
			a.checkAuthWBLogistic()
			return err
		}
	}
	return nil
}

func (a *App) financeRoutesHandler(ctx context.Context) error {
	select {
	case <-ctx.Done():
		logger.Log(logger.ERROR, "App.financeRoutesHandler()", "Cancelling 'finance_routes_report' callback task")
		return ctx.Err()
	default:
		err := a.financeRoutesReporter.Run(ctx)
		if err != nil {
			logger.Log(logger.ERROR, "App.financeRoutesHandler()", "Failed to run finance routes report")
			a.checkAuthWBLogistic()
			return err
		}
	}
	return nil
}

func (a *App) financeDailyHandler(ctx context.Context) error {
	select {
	case <-ctx.Done():
		logger.Log(logger.ERROR, "App.financeDailyHandler()", "Cancelling 'finance_daily_report' callback task")
		return ctx.Err()
	default:
		err := a.financeDailyReporter.Run(ctx)
		if err != nil {
			logger.Log(logger.ERROR, "App.financeDailyHandler()", "Failed to run finance daily report")
			a.checkAuthWBLogistic()
			return err
		}
	}
	return nil
}

func (a *App) checkAuthWBLogistic() {
	if a.services.WBLogisticService.IsSessionExpired() {
		logger.Log(logger.ERROR, "App.checkAuthWBLogistic()", "WB logistic session expired")
		go func() {
			a.Pause()

			err := a.initializer.InitDirectWBLogistic()
			if err != nil {
				logger.Logf(logger.ERROR, "App.checkAuthWBLogistic()", "failed to init direct wb logistic: %v", err)
				return
			}

			err = a.Start()
			if err != nil {
				logger.Logf(logger.ERROR, "App.checkAuthWBLogistic()", "failed to start app: %v", err)
				return
			}
		}()
	}
}
