package initializer

import (
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/initializer/telegram_bot"
	"wb_logistic_assistant/internal/initializer/wb_logistic"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/models"
	"wb_logistic_assistant/internal/prompters"
	"wb_logistic_assistant/internal/reporters"
	"wb_logistic_assistant/internal/scheduler"
	"wb_logistic_assistant/internal/services"
	"wb_logistic_assistant/internal/storage"
)

type Initializer struct {
	config       *config.Config
	storage      storage.Storage
	prompter     prompters.InitializeAppPrompter
	dependencies *AppDependencies
}

func NewInitializer(config *config.Config, storage storage.Storage, prompter prompters.InitializeAppPrompter) *Initializer {
	return &Initializer{
		config:   config,
		storage:  storage,
		prompter: prompter,
		dependencies: &AppDependencies{
			Config:  config,
			Storage: storage,
		},
	}
}

func (i *Initializer) Init() (*AppDependencies, error) {
	logger.Log(logger.INFO, "Initializer.Init()", "Start init application dependencies")

	//// Services
	servicesContainer := &services.Container{}
	servicesTTL := i.config.Logistic().CacheTTL()

	//googleSheetsInitializer := google_sheets.NewInitializer(i.config, i.storage, i.prompter)
	//googleSheetsClient, googleSheetsActor, err := googleSheetsInitializer.Init()
	//if err != nil {
	//	return nil, errors.Wrap(err, "Initializer.Init()", "Failed to init Google Sheets client")
	//}
	//servicesContainer.GoogleSheetsService = services.NewBaseGoogleSheetsService(googleSheetsClient, googleSheetsActor)
	//
	//err = i.updateStorage()
	//if err != nil {
	//	return nil, err
	//}

	wbLogisticInitializer := wb_logistic.NewInitializer(i.config, i.storage, i.prompter)
	wbLogisticClient, wbLogisticSession, err := wbLogisticInitializer.Init()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.Init()", "Failed to init WB Logistic client")
	}
	servicesContainer.WBLogisticService = services.NewBaseWBLogisticService(wbLogisticClient, wbLogisticSession, &models.WBLogisticTTlParams{
		UserInfo:                        servicesTTL.UserInfo(),
		RemainsLastMileReports:          servicesTTL.RemainsLastMileReports(),
		RemainsLastMileReportsRouteInfo: servicesTTL.RemainsLastMileReportsRouteInfo(),
		JobsScheduling:                  servicesTTL.JobsScheduling(),
		ShipmentInfo:                    servicesTTL.ShipmentInfo(),
		ShipmentTransfers:               servicesTTL.ShipmentTransfers(),
		WaySheetInfo:                    servicesTTL.WaySheetInfo(),
	})

	telegramBotInitializer := telegram_bot.NewInitializer(i.config, i.storage, i.prompter)
	telegramBot, err := telegramBotInitializer.Init()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.Init()", "Failed to init Telegram Bot client")
	}
	servicesContainer.TelegramBotService = services.NewTelegramBotAPIService(telegramBot)
	i.dependencies.Services = servicesContainer

	err = i.updateStorage()
	if err != nil {
		return nil, err
	}

	//// BaseScheduler
	logger.Log(logger.INFO, "Initializer.initServices()", "Start init application scheduler")
	i.dependencies.Scheduler = scheduler.NewBaseScheduler(i.config.Internal().SchedulerMaxWorkers(), i.config.Internal().SchedulerRetryTaskLimit())
	logger.Log(logger.INFO, "Initializer.initServices()", "Finish init application scheduler, successfully initialized")

	///// Orchestrators
	logger.Log(logger.INFO, "Initializer.initServices()", "Start init application orchestrators")
	i.dependencies.GeneralRoutesReporter = reporters.NewGeneralRoutesReporter(i.config, i.storage, servicesContainer, &prompters.CLIReporterGeneralRoutesPrompter{})
	i.dependencies.ShipmentCloseReporter = reporters.NewShipmentCloseReporter(i.config, i.storage, servicesContainer, &prompters.CLIReporterShipmentClosePrompter{})
	logger.Log(logger.INFO, "Initializer.initServices()", "Finish init application orchestrators, successfully initialized")

	logger.Log(logger.INFO, "Initializer.Init()", "Finish init application dependencies, successfully initialized")
	return i.dependencies, nil
}

func (i *Initializer) updateStorage() error {
	err := i.storage.Save(i.config.Storage().Path())
	if err != nil {
		return errors.Wrap(err, "Initializer.updateStorage()", "Failed to update storage")
	}
	return nil
}
