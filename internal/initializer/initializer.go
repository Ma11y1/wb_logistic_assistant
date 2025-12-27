package initializer

import (
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/initializer/google_sheets"
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
	services     *services.Container
	wbLogistic   *wb_logistic.Initializer
	googleSheets *google_sheets.Initializer
	telegramBot  *telegram_bot.Initializer
}

func NewInitializer(config *config.Config, storage storage.Storage, prompter prompters.InitializeAppPrompter) *Initializer {
	services := &services.Container{}
	return &Initializer{
		config:   config,
		storage:  storage,
		prompter: prompter,
		dependencies: &AppDependencies{
			Config:   config,
			Storage:  storage,
			Services: services,
		},
		services:     services,
		wbLogistic:   wb_logistic.NewInitializer(config, storage, prompter),
		googleSheets: google_sheets.NewInitializer(config, storage, prompter),
		telegramBot:  telegram_bot.NewInitializer(config, storage, prompter),
	}
}

func (i *Initializer) Init() (*AppDependencies, error) {
	logger.Log(logger.INFO, "Initializer.Init()", "Start init application dependencies")

	//// Services
	err := i.initWBLogistic(i.config.Logistic().CacheTTL())
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.Init()", "")
	}

	err = i.initGoogleSheets()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.Init()", "")
	}

	err = i.initTelegramBot()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.Init()", "")
	}

	i.initScheduler()

	i.initOrchestrators()

	logger.Log(logger.INFO, "Initializer.Init()", "Finish init application dependencies, successfully initialized")
	i.prompter.PromptInitFinish()
	return i.dependencies, nil
}

func (i *Initializer) initWBLogistic(ttl *config.LogisticCacheTTL) error {
	if i.config.Reports().GeneralRoutes().IsEnabled() ||
		i.config.Reports().ShipmentClose().IsEnabled() ||
		i.config.Reports().FinanceRoutes().IsEnabled() ||
		i.config.Reports().FinanceDaily().IsEnabled() {

		wbLogisticClient, wbLogisticSession, err := i.wbLogistic.Init()
		if err != nil {
			return errors.Wrap(err, "Initializer.initWBLogistic()", "Failed to init WB Logistic client")
		}
		i.services.WBLogisticService = services.NewBaseWBLogisticService(wbLogisticClient, wbLogisticSession, &models.WBLogisticTTlParams{
			UserInfo:                        ttl.UserInfo(),
			RemainsLastMileReports:          ttl.RemainsLastMileReports(),
			RemainsLastMileReportsRouteInfo: ttl.RemainsLastMileReportsRouteInfo(),
			JobsScheduling:                  ttl.JobsScheduling(),
			ShipmentInfo:                    ttl.ShipmentInfo(),
			ShipmentTransfers:               ttl.ShipmentTransfers(),
			WaySheetInfo:                    ttl.WaySheetInfo(),
			WaySheetFinanceDetails:          ttl.WaySheetFinanceDetails(),
		})
	} else {
		return errors.New("Initializer.initWBLogistic()", "All reports are disabled")
	}
	return nil
}

// InitDirectWBLogistic init without storage data
func (i *Initializer) InitDirectWBLogistic() error {
	if i.config.Reports().GeneralRoutes().IsEnabled() ||
		i.config.Reports().ShipmentClose().IsEnabled() ||
		i.config.Reports().FinanceRoutes().IsEnabled() ||
		i.config.Reports().FinanceDaily().IsEnabled() {

		wbLogisticClient, wbLogisticSession, err := i.wbLogistic.InitDirect()
		if err != nil {
			return errors.Wrap(err, "Initializer.initWBLogistic()", "Failed to init WB Logistic client")
		}
		i.services.WBLogisticService.SetClient(wbLogisticClient)
		i.services.WBLogisticService.SetSession(wbLogisticSession)
	} else {
		return errors.New("Initializer.initWBLogistic()", "All reports are disabled")
	}
	return nil
}

func (i *Initializer) initGoogleSheets() error {
	if (i.config.Reports().GeneralRoutes().IsEnabled() && i.config.Reports().GeneralRoutes().IsRenderGoogleSheets()) ||
		(i.config.Reports().ShipmentClose().IsEnabled() && i.config.Reports().ShipmentClose().IsRenderGoogleSheets()) {
		googleSheetsClient, googleSheetsActor, err := i.googleSheets.Init()
		if err != nil {
			return errors.Wrap(err, "Initializer.initGoogleSheets()", "Failed to init Google Sheets client")
		}
		i.services.GoogleSheetsService = services.NewBaseGoogleSheetsService(googleSheetsClient, googleSheetsActor)
	}
	return nil
}

func (i *Initializer) initTelegramBot() error {
	if (i.config.Reports().ShipmentClose().IsEnabled() && i.config.Reports().ShipmentClose().IsRenderTelegramBot()) ||
		(i.config.Reports().FinanceRoutes().IsEnabled() && i.config.Reports().FinanceRoutes().IsRenderTelegramBot()) ||
		(i.config.Reports().FinanceDaily().IsEnabled() && i.config.Reports().FinanceDaily().IsRenderTelegramBot()) {
		telegramBot, err := i.telegramBot.Init()
		if err != nil {
			return errors.Wrap(err, "Initializer.initTelegramBot()", "Failed to init Telegram Bot client")
		}
		i.services.TelegramBotService = services.NewTelegramBotAPIService(telegramBot)
	}
	return nil
}

func (i *Initializer) initScheduler() {
	logger.Log(logger.INFO, "Initializer.initServices()", "Start init application scheduler")
	i.dependencies.Scheduler = scheduler.NewBaseScheduler(i.config.Internal().SchedulerMaxWorkers(), i.config.Internal().SchedulerRetryTaskLimit())

	i.dependencies.SchedulerGeneralRoutesTaskConfig = &scheduler.TaskConfig{
		RetryTaskLimit:        i.config.Reports().GeneralRoutes().ErrRetryTaskLimit(),
		Timeout:               i.config.Reports().GeneralRoutes().TaskTimeout(),
		IsWaitForPrevious:     true,
		IsIntervalAfterFinish: true,
	}

	i.dependencies.SchedulerShipmentCloseTaskConfig = &scheduler.TaskConfig{
		RetryTaskLimit:        i.config.Reports().ShipmentClose().ErrRetryTaskLimit(),
		Timeout:               i.config.Reports().ShipmentClose().TaskTimeout(),
		IsWaitForPrevious:     true,
		IsIntervalAfterFinish: true,
	}

	i.dependencies.SchedulerFinanceRoutesTaskConfig = &scheduler.TaskConfig{
		RetryTaskLimit:        i.config.Reports().FinanceRoutes().ErrRetryTaskLimit(),
		Timeout:               i.config.Reports().FinanceRoutes().TaskTimeout(),
		IsWaitForPrevious:     true,
		IsIntervalAfterFinish: true,
	}

	i.dependencies.SchedulerFinanceDailyTaskConfig = &scheduler.TaskConfig{
		RetryTaskLimit:        i.config.Reports().FinanceDaily().ErrRetryTaskLimit(),
		Timeout:               i.config.Reports().FinanceDaily().TaskTimeout(),
		IsWaitForPrevious:     true,
		IsIntervalAfterFinish: true,
	}
	logger.Log(logger.INFO, "Initializer.initServices()", "Finish init application scheduler, successfully initialized")
}

func (i *Initializer) initOrchestrators() {
	logger.Log(logger.INFO, "Initializer.initServices()", "Start init application orchestrators")
	i.dependencies.GeneralRoutesReporter = reporters.NewGeneralRoutesReporter(i.config, i.storage, i.services, &prompters.CLIReporterGeneralRoutesPrompter{})
	i.dependencies.ShipmentCloseReporter = reporters.NewShipmentCloseReporter(i.config, i.storage, i.services, &prompters.CLIReporterShipmentClosePrompter{})
	i.dependencies.FinanceRoutesReporter = reporters.NewFinanceRoutesReporter(i.config, i.storage, i.services, &prompters.CLIReporterFinanceRoutesPrompter{})
	i.dependencies.FinanceDailyReporter = reporters.NewFinanceDailyReporter(i.config, i.storage, i.services, &prompters.CLIReporterFinanceDailyPrompter{})
	logger.Log(logger.INFO, "Initializer.initServices()", "Finish init application orchestrators, successfully initialized")

}
