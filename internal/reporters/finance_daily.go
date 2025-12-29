package reporters

import (
	"context"
	"fmt"
	"time"
	wb_models "wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/internal/models"

	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/prompters"
	"wb_logistic_assistant/internal/report_renderers"
	"wb_logistic_assistant/internal/reports"
	"wb_logistic_assistant/internal/services"
	"wb_logistic_assistant/internal/storage"
)

type FinanceDailyReporterData struct {
	DateStart          time.Time
	DateEnd            time.Time
	RouteID            int
	Parking            int
	Flights            int
	OpenedFlights      int
	ShippedBarcodes    int
	Tare               int
	ShippedTare        int
	ReturnTare         int
	Income             float64
	IncomeReturn       float64
	Fine               float64
	TotalSalaryRate    float64
	SalaryRate         float64
	ExtendedSalaryRate float64
	Marriage           float64
	PercentMarriage    float64
	Tax                float64
	PercentTax         float64
	Margin             float64
	WaySheetIDs        []string
	OpenedWaySheetIDs  []string
}

type FinanceDailyReporter struct {
	config        *config.Config
	storage       storage.Storage
	services      *services.Container
	reportRoute   *reports.FinanceDailyRouteReport
	reportGeneral *reports.FinanceDailyGeneralReport
	prompter      prompters.FinanceDailyReporterPrompter

	messageQueueTG *Queue[string]
	rendererTG     report_renderers.ReportRenderer[[]string]
	isRenderTG     bool
	tgChatID       int64

	officeID                int
	suppliers               map[int]struct{} // supplier id -> struct{}
	skipRoutes              map[int]struct{} // route id -> struct{}
	intervalUpdateWaySheets time.Duration
	prevTimeUpdateWaySheets time.Time
	dayOffset               int
	salaryRatePercent       map[int]float64 // route id -> salary rate
	salaryRate              map[int]float64 // route id -> salary rate
	percentTax              float64
	taxRate                 float64
	percentMarriage         float64
	marriageRate            float64
	expensesDaily           float64
	timeLastRender          time.Time
	isRender                bool

	chData map[int]*FinanceDailyReporterData
}

func NewFinanceDailyReporter(config *config.Config, storage storage.Storage, service *services.Container, prompter prompters.FinanceDailyReporterPrompter) *FinanceDailyReporter {
	expensesDaily := 0.0
	if config.Logistic().Office().Expenses() != 0 && config.Logistic().Office().ExpensesPeriod() != 0 {
		expensesDaily = config.Logistic().Office().Expenses() / float64(config.Logistic().Office().ExpensesPeriod())
	}
	return &FinanceDailyReporter{
		config:        config,
		storage:       storage,
		services:      service,
		prompter:      prompter,
		reportRoute:   &reports.FinanceDailyRouteReport{},
		reportGeneral: &reports.FinanceDailyGeneralReport{},

		messageQueueTG: New[string](300),
		rendererTG:     &report_renderers.TelegramBotRenderer{Mode: report_renderers.TelegramBotRenderHTML},
		isRenderTG:     config.Reports().FinanceDaily().IsRenderTelegramBot(),
		tgChatID:       config.Telegram().FinanceDaily().ChatID(),

		officeID:          config.Logistic().Office().ID(),
		suppliers:         config.Logistic().Office().SuppliersMap(),
		skipRoutes:        config.Logistic().Office().SkipRoutesMap(),
		salaryRatePercent: config.Logistic().Office().SalaryRatePercent(),
		salaryRate:        config.Logistic().Office().SalaryRate(),
		percentTax:        config.Logistic().Office().PercentTax(),
		taxRate:           config.Logistic().Office().PercentTax() / 100,
		percentMarriage:   config.Logistic().Office().PercentMarriage(),
		marriageRate:      config.Logistic().Office().PercentMarriage() / 100,
		expensesDaily:     expensesDaily,
		dayOffset:         config.Reports().FinanceDaily().DayOffset(),
		isRender:          config.Reports().FinanceDaily().RenderAtStart(),

		chData: map[int]*FinanceDailyReporterData{},
	}
}

func (r *FinanceDailyReporter) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "FinanceDailyReporter.Run()", "task was preliminarily completed")
	}

	now := time.Now()
	r.prompter.PromptStart(time.Date(now.Year(), now.Month(), now.Day()+r.dayOffset, 0, 0, 0, 0, time.UTC))

	err := r.sendRemainingMessagesTelegramBot(ctx)
	if err != nil {
		r.prompter.PromptError("Failed send remains messages to Telegram bot")
		return errors.Wrap(err, "FinanceDailyReporter.Run()", "failed sending remaining messages to Telegram bot")
	}

	// isRender is true if it was originally set this way in the configuration or if the day has changed since the function was last run
	if !r.timeLastRender.IsZero() && r.timeLastRender.Day() != now.Day() {
		r.isRender = true
	}
	r.timeLastRender = now

	if !r.isRender {
		r.prompter.PromptFinish(time.Since(now))
		return nil
	}

	err = r.processWaySheets(ctx)
	if err != nil {
		return errors.Wrap(err, "FinanceDailyReporter.Run()", "failed processing way sheets")
	}

	err = r.processReports(ctx)
	if err != nil {
		return errors.Wrap(err, "FinanceDailyReporter.Run()", "failed processing reports")
	}

	r.isRender = false
	r.prompter.PromptFinish(time.Since(now))
	return nil
}

func (r *FinanceDailyReporter) processWaySheets(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "FinanceDailyReporter.processWaySheets()", "task was preliminarily completed")
	}

	logger.Log(logger.INFO, "FinanceDailyReporter.processWaySheets()", "start process way sheets")

	now := time.Now()
	timeStart := time.Date(now.Year(), now.Month(), now.Day()+r.dayOffset, 0, 0, 0, 0, time.UTC)
	timeEnd := time.Date(now.Year(), now.Month(), now.Day()+r.dayOffset, 23, 59, 59, 999999999, time.UTC)

	waySheets, err := r.loadWaySheets(ctx)
	if err != nil {
		return errors.Wrap(err, "FinanceDailyReporter.processWaySheets()", "failed process way sheets")
	}

	clear(r.chData)

	targetClosed := 0
	targetOpened := 0
	for _, waySheet := range waySheets {
		if waySheet == nil {
			continue
		}

		if waySheet.OpenDt.Before(timeStart) || waySheet.OpenDt.After(timeEnd) {
			continue
		}

		supplierID := atoiSafe(waySheet.SupplierID)
		if !r.isValidSupplier(supplierID) {
			continue
		}

		routeID := atoiSafe(waySheet.RouteCarID)
		data, ok := r.chData[routeID]
		if !ok {
			data = &FinanceDailyReporterData{DateStart: timeStart, DateEnd: timeEnd, RouteID: routeID}
			r.chData[routeID] = data

			time.Sleep(200 * time.Millisecond)
			waySheetID := atoiSafe(waySheet.WaySheetID)
			info, err := r.loadWaySheetInfo(ctx, waySheetID)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed loading way sheet info for route %d, way sheet %s", routeID, waySheet.WaySheetID))
				logger.Logf(logger.ERROR, "FinanceDailyReporter.processWaySheets()", "failed load way sheet info for way sheet %d: %v", waySheetID, err)
			} else {
				parking, err := r.loadParking(ctx, info.DstOffices)
				if err != nil {
					r.prompter.PromptError(fmt.Sprintf("Failed load parking for route %d, way sheet %s", routeID, waySheet.WaySheetID))
					logger.Logf(logger.ERROR, "FinanceDailyReporter.processWaySheets()", "failed loading parking info for way sheet %d: %v", waySheetID, err)
				}
				data.Parking = parking
			}
		}

		data.Flights++

		if waySheet.CloseDt.IsZero() {
			data.OpenedWaySheetIDs = append(data.OpenedWaySheetIDs, waySheet.WaySheetID)
			data.OpenedFlights++
			targetOpened++
			continue
		}
		targetClosed++

		data.ShippedBarcodes += atoiSafe(waySheet.CountBarcodes)

		tare := atoiSafe(waySheet.CountBox)
		shippedTare := atoiSafe(waySheet.CountArrivalBox)
		data.Tare += tare
		data.ShippedTare += shippedTare
		data.ReturnTare += tare - shippedTare

		salaryRate := data.SalaryRate
		if salaryRate == 0 {
			if v, ok := r.salaryRate[routeID]; ok {
				salaryRate = v
			} else if v, ok := r.salaryRatePercent[routeID]; ok {
				salaryRate = waySheet.TotalPrice * (v / 100) // calculated if the rate is in percentages and not fixed
			} else {
				r.prompter.PromptError(fmt.Sprintf("There is no salary rate data for the route %d, way sheet %s", routeID, waySheet.WaySheetID))
				logger.Logf(logger.ERROR, "FinanceDailyReporter.processWaySheets()", "there is no salary rate data for the route %d, way sheet %s", routeID, waySheet.WaySheetID)
			}
			data.SalaryRate = salaryRate
		}

		totalPriceSubFine := waySheet.TotalPrice - waySheet.SumFine
		marriage := totalPriceSubFine * r.marriageRate
		tax := (totalPriceSubFine - marriage) * r.taxRate

		data.Income += waySheet.TotalPrice
		data.IncomeReturn += waySheet.SumReturn
		data.Fine += waySheet.SumFine
		data.TotalSalaryRate += salaryRate
		data.ExtendedSalaryRate += salaryRate + marriage + tax
		data.Marriage += marriage
		data.PercentMarriage = r.percentMarriage
		data.Tax += tax
		data.PercentTax = r.percentTax
		data.Margin += totalPriceSubFine - (salaryRate + marriage + tax)
		data.WaySheetIDs = append(data.WaySheetIDs, waySheet.WaySheetID)
	}

	r.prompter.PromptCountWaySheet(len(waySheets), targetClosed, targetOpened)

	return nil
}

func (r *FinanceDailyReporter) processReports(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "FinanceDailyReporter.processReports()", "task was preliminarily completed")
	}

	logger.Log(logger.INFO, "FinanceDailyReporter.processReports()", "start process reports")

	var dateStart, dateEnd time.Time
	var flights, openedFlights, shippedBarcodes, shippedTare, tare, returnedTare int
	var income, incomeReturn, fine, salaryRate, extendedSalaryRate, tax, marriage, margin float64
	var openedWaySheets []string

	for routeID, data := range r.chData {
		if data == nil {
			delete(r.chData, routeID)
			continue
		}

		dateStart = data.DateStart
		dateEnd = data.DateEnd
		flights += data.Flights
		openedFlights += data.OpenedFlights
		shippedBarcodes += data.ShippedBarcodes
		tare += data.Tare
		shippedTare += data.ShippedTare
		returnedTare += data.ReturnTare
		income += data.Income
		incomeReturn += data.IncomeReturn
		fine += data.Fine
		salaryRate += data.TotalSalaryRate
		extendedSalaryRate += data.ExtendedSalaryRate
		tax += data.Tax
		marriage += data.Marriage
		margin += data.Margin

		if len(data.OpenedWaySheetIDs) > 0 {
			openedWaySheets = append(openedWaySheets, data.OpenedWaySheetIDs...)
		}

		renderData, err := r.renderRouteReport(&reports.FinanceDailyRouteReportData{
			Date:               dateStart,
			RouteID:            routeID,
			Parking:            data.Parking,
			Flights:            data.Flights,
			ShippedBarcodes:    data.ShippedBarcodes,
			Tare:               data.Tare,
			ShippedTare:        data.ShippedTare,
			ReturnedTare:       data.ReturnTare,
			Income:             data.Income,
			IncomeReturn:       data.IncomeReturn,
			Fine:               data.Fine,
			TotalSalaryRate:    data.TotalSalaryRate,
			SalaryRate:         data.SalaryRate,
			ExtendedSalaryRate: data.ExtendedSalaryRate,
			Marriage:           data.Marriage,
			PercentMarriage:    data.PercentMarriage,
			Tax:                data.Tax,
			PercentTax:         data.PercentTax,
			Margin:             data.Margin,
			WaySheetIDs:        data.WaySheetIDs,
			OpenedWaySheets:    data.OpenedWaySheetIDs,
		})
		if err != nil {
			r.prompter.PromptError(fmt.Sprintf("Failed render route report, route id %d: %v", routeID, err))
			logger.Logf(logger.ERROR, "FinanceDailyReporter.processReports()", "failed render route report, route id %d: %v", routeID, err)
			continue
		}

		err = r.sendReport(ctx, renderData)
		if err != nil {
			r.prompter.PromptError(fmt.Sprintf("Failed send route report, route id %d: %v", routeID, err))
			logger.Logf(logger.ERROR, "FinanceDailyReporter.processReports()", "failed send route report, route id %d: %v", routeID, err)

		} else {
			r.prompter.PromptSendReport(routeID)
			logger.Logf(logger.INFO, "FinanceDailyReporter.processReports()", "send report for route id %d", routeID)
		}

		time.Sleep(3 * time.Second)
	}

	renderData, err := r.renderGeneralReport(&reports.FinanceDailyGeneralReportData{
		DateStart:          dateStart,
		DateEnd:            dateEnd,
		Flights:            flights,
		OpenedFlights:      openedFlights,
		ShippedBarcodes:    shippedBarcodes,
		ShippedTare:        shippedTare,
		Tare:               tare,
		ReturnedTare:       returnedTare,
		Income:             income,
		IncomeReturn:       incomeReturn,
		Fine:               fine,
		SalaryRate:         salaryRate,
		ExtendedSalaryRate: extendedSalaryRate,
		Marriage:           marriage,
		PercentMarriage:    r.percentMarriage,
		Tax:                tax,
		PercentTax:         r.percentTax,
		Margin:             margin,
		Expenses:           r.expensesDaily,
		TotalMargin:        margin - r.expensesDaily,
		OpenedWaySheets:    openedWaySheets,
	})
	if err != nil {
		return errors.Wrap(err, "FinanceDailyReporter.processReports()", "failed render general report")
	}

	err = r.sendReport(ctx, renderData)
	if err != nil {
		return errors.Wrap(err, "FinanceDailyReporter.processReports()", "failed send general report")
	}

	return nil
}

func (r *FinanceDailyReporter) isValidSupplier(supplierID int) bool {
	if _, ok := r.suppliers[supplierID]; !ok {
		return false
	}
	return true
}

func (r *FinanceDailyReporter) loadParking(ctx context.Context, dstOffices []*wb_models.WaySheetDestinationOffice) (int, error) {
	if dstOffices == nil || len(dstOffices) == 0 {
		return 0, nil
	}

	offices := make([]int, 0, len(dstOffices))
	for _, office := range dstOffices {
		if office == nil {
			continue
		}
		id := atoiSafe(office.ID)
		if id > 0 {
			offices = append(offices, id)
		}
	}

	remainsTares, err := r.loadRemainsTares(ctx, offices)
	if err != nil {
		return 0, errors.Wrapf(err, "FinanceDailyReporter.loadParking()", "failed load remains tares")
	}

	if len(remainsTares) == 0 {
		return 0, nil
	}

	_, parking := SpNameToGateParking(remainsTares[0].SpName)

	return parking, nil
}

func (r *FinanceDailyReporter) loadRemainsTares(ctx context.Context, dstOfficeIDs []int) (tares []*wb_models.TareForOffice, err error) {
	err = retryAction(ctx, "FinanceDailyReporter.loadRemainsTares", 3, 1*time.Second, func() error {
		// isDrive = false is default value
		tares, err = r.services.WBLogisticService.GetTaresForOffices(ctx, r.officeID, dstOfficeIDs, false)
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "FinanceDailyReporter.loadRemainsTares()", "failed load tares for offices")
	}
	return tares, nil
}

func (r *FinanceDailyReporter) loadWaySheets(ctx context.Context) (waySheets []*wb_models.WaySheet, err error) {
	now := time.Now()
	for supplierID := range r.suppliers {
		var page *wb_models.WaySheetsPage
		err = retryAction(ctx, "FinanceDailyReporter.loadWaySheets", 3, 1*time.Second, func() error {
			page, err = r.services.WBLogisticService.GetWaySheets(ctx, &models.WBLogisticGetWaySheetsParamsRequest{
				DateOpen:    time.Date(now.Year(), now.Month(), now.Day()+r.dayOffset, 0, 0, 0, 0, time.UTC),
				DateClose:   time.Date(now.Year(), now.Month(), now.Day()+r.dayOffset, 23, 59, 59, 999999999, time.UTC),
				SupplierID:  supplierID,
				SrcOfficeID: r.officeID,
				Offset:      0,
				Limit:       1000,
				WayTypeID:   0,
			})
			return err
		})
		if page == nil || err != nil {
			r.prompter.PromptError(fmt.Sprintf("failed load way sheets for supplier %d", supplierID))
			logger.Logf(logger.ERROR, "FinanceDailyReporter.loadWaySheets()", "failed load way sheets for supplier %d: %v", supplierID, err)
			continue
		}
		if len(page.WaySheets) > 0 {
			if len(waySheets) == 0 {
				waySheets = page.WaySheets
			} else {
				waySheets = append(waySheets, page.WaySheets...)
			}
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "FinanceDailyReporter.loadWaySheets()", "failed load way sheets")
	}
	return waySheets, nil
}

func (r *FinanceDailyReporter) loadWaySheetInfo(ctx context.Context, waySheetID int) (waySheetInfo *wb_models.WaySheetInfo, err error) {
	err = retryAction(ctx, "FinanceDailyReporter.loadWaySheetInfo", 3, 1*time.Second, func() error {
		waySheetInfo, err = r.services.WBLogisticService.GetWaySheetInfo(ctx, waySheetID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "FinanceDailyReporter.loadWaySheetInfo()", "failed load way sheet %d info", waySheetID)
	}

	if waySheetInfo == nil {
		return nil, errors.Newf("FinanceDailyReporter.loadWaySheetInfo()", "way sheet %d info returned empty value without error", waySheetID)
	}

	return waySheetInfo, nil
}

func (r *FinanceDailyReporter) renderRouteReport(reportData *reports.FinanceDailyRouteReportData) (*reports.ReportData, error) {
	if reportData == nil {
		return nil, errors.New("FinanceDailyReporter.renderRouteReport()", "route report data is empty")
	}
	report, err := r.reportRoute.Render(reportData)
	if err != nil {
		r.prompter.PromptError("failed to render route report")
		return nil, errors.Wrapf(err, "FinanceDailyReporter.renderRouteReport()", "failed render route report, route id: %d", reportData.RouteID)
	}
	return report, nil
}

func (r *FinanceDailyReporter) renderGeneralReport(reportData *reports.FinanceDailyGeneralReportData) (*reports.ReportData, error) {
	if reportData == nil {
		return nil, errors.New("FinanceDailyReporter.renderGeneralReport()", "general report data is empty")
	}
	report, err := r.reportGeneral.Render(reportData)
	if err != nil {
		r.prompter.PromptError("failed to render general report")
		return nil, errors.Wrapf(err, "FinanceDailyReporter.renderGeneralReport()", "failed render general report")
	}
	return report, nil
}

func (r *FinanceDailyReporter) sendReport(ctx context.Context, data *reports.ReportData) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "FinanceDailyReporter.sendReport()", "task was preliminarily completed")
	}

	if r.isRenderTG {
		messages, err := r.rendererTG.Render(data)
		if err != nil {
			r.prompter.PromptError("failed to render report Telegram Bot")
			logger.Logf(logger.ERROR, "FinanceDailyReporter.sendReport()", "failed render report for Telegram Bot: %v", err)
		} else {
			if len(messages) == 0 {
				return nil
			}

			for _, message := range messages {
				if message != "" {
					r.messageQueueTG.Push(message)
				}
			}

			if err = r.sendTelegramBot(ctx); err != nil {
				r.prompter.PromptError("failed to send report Telegram Bot")
				logger.Logf(logger.ERROR, "FinanceDailyReporter.sendReport()", "failed send report to Telegram Bot: %v", err)
			}
		}
	}

	return nil
}

func (r *FinanceDailyReporter) sendTelegramBot(ctx context.Context) error {
	for r.messageQueueTG.Len() > 0 {
		err := retryAction(ctx, "FinanceDailyReporter.sendTelegramBot", 3, 1*time.Second, func() error {
			message, ok := r.messageQueueTG.Peek()
			if !ok {
				return errors.New("FinanceDailyReporter.sendTelegramBot()", "failed to get message from telegram message queue")
			}
			return r.services.TelegramBotService.SendMessage(r.tgChatID, message, "HTML")
		})
		if err != nil {
			return errors.Wrapf(err, "FinanceDailyReporter.sendTelegramBot()", "failed send data to chat %d", r.tgChatID)
		}
		r.messageQueueTG.Pop()
	}

	return nil
}

func (r *FinanceDailyReporter) sendRemainingMessagesTelegramBot(ctx context.Context) error {
	if r.messageQueueTG.Len() <= 0 {
		return nil
	}

	err := r.sendTelegramBot(ctx)
	if err != nil {
		return errors.Wrap(err, "FinanceDailyReporter.sendRemainingMessagesTelegramBot)", "failed sending remains messages to Telegram bot")
	}
	logger.Log(logger.INFO, "FinanceDailyReporter.sendRemainingMessagesTelegramBot()", "send remains messages to Telegram bot")

	return nil
}
