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

type FinanceRoutesReporter struct {
	config   *config.Config
	storage  storage.Storage
	services *services.Container
	report   *reports.FinanceRoutesReport
	prompter prompters.FinanceRoutesReporterPrompter

	messageQueueTG *Queue[string]
	rendererTG     report_renderers.ReportRenderer[[]string]
	tgChatID       int64
	isRenderTG     bool

	officeID          int
	suppliers         map[int]struct{} // supplier id -> struct{}
	skipRoutes        map[int]struct{} // route id -> struct{}
	salaryRatePercent map[int]float64  // route id -> salary rate
	salaryRate        map[int]float64  // route id -> salary rate
	percentTax        float64
	taxRate           float64
	percentMarriage   float64
	marriageRate      float64

	chOpenedWaySheets map[string]*wb_models.WaySheet // way sheet id -> way sheet
}

func NewFinanceRoutesReporter(config *config.Config, storage storage.Storage, service *services.Container, prompter prompters.FinanceRoutesReporterPrompter) *FinanceRoutesReporter {
	return &FinanceRoutesReporter{
		config:   config,
		storage:  storage,
		services: service,
		prompter: prompter,
		report:   &reports.FinanceRoutesReport{},

		messageQueueTG: New[string](300),
		rendererTG:     &report_renderers.TelegramBotRenderer{Mode: report_renderers.TelegramBotRenderHTML},
		tgChatID:       config.Telegram().FinanceRoutes().ChatID(),
		isRenderTG:     config.Reports().FinanceRoutes().IsRenderTelegramBot(),

		officeID:          config.Logistic().Office().ID(),
		suppliers:         config.Logistic().Office().SuppliersMap(),
		skipRoutes:        config.Logistic().Office().SkipRoutesMap(),
		salaryRatePercent: config.Logistic().Office().SalaryRatePercent(),
		salaryRate:        config.Logistic().Office().SalaryRate(),
		percentTax:        config.Logistic().Office().PercentTax(),
		taxRate:           config.Logistic().Office().PercentTax() / 100,
		percentMarriage:   config.Logistic().Office().PercentMarriage(),
		marriageRate:      config.Logistic().Office().PercentMarriage() / 100,

		chOpenedWaySheets: map[string]*wb_models.WaySheet{},
	}
}

func (r *FinanceRoutesReporter) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "FinanceRoutesReporter.Run()", "task was preliminarily completed")
	}

	r.prompter.PromptStart()

	err := r.sendRemainingMessagesTelegramBot(ctx)
	if err != nil {
		r.prompter.PromptError("Failed to send remaining messages to Telegram bot")
		return errors.Wrap(err, "FinanceRoutesReporter.Run()", "failed sending remaining messages to Telegram bot")
	}

	now := time.Now()

	if err := r.findOpenedWaySheets(ctx); err != nil {
		r.prompter.PromptError("Failed finding opened way sheets")
		return errors.Wrap(err, "FinanceRoutesReporter.Run()", "failed finding opened way sheets")
	}

	if err := r.processOpenedWaySheets(ctx); err != nil {
		r.prompter.PromptError("Failed processing opened way sheets")
		return errors.Wrap(err, "FinanceRoutesReporter.Run()", "failed processing opened way sheets")
	}

	r.prompter.PromptFinish(time.Since(now))
	return nil
}

func (r *FinanceRoutesReporter) findOpenedWaySheets(ctx context.Context) error {
	waySheets, err := r.loadWaySheets(ctx)
	if err != nil {
		return errors.Wrap(err, "FinanceRoutesReporter.findOpenedWaySheets()", "failed loading way sheets")
	}

	for _, waySheet := range waySheets {
		if waySheet == nil {
			continue
		}
		_, ok := r.chOpenedWaySheets[waySheet.WaySheetID]
		if !ok && !waySheet.CloseDt.IsZero() {
			continue
		}
		if r.isValidRoute(atoiSafe(waySheet.RouteCarID), atoiSafe(waySheet.SupplierID)) {
			r.chOpenedWaySheets[waySheet.WaySheetID] = waySheet
		}
	}

	r.prompter.PromptCountWaySheet(len(r.chOpenedWaySheets))
	logger.Logf(logger.INFO, "FinanceRoutesReporter.findOpenedWaySheets()", "opened way sheets: %d", len(r.chOpenedWaySheets))

	return nil
}

func (r *FinanceRoutesReporter) processOpenedWaySheets(ctx context.Context) error {
	for id, waySheet := range r.chOpenedWaySheets {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if waySheet.TotalPrice == 0 {
			// financial data may appear with a delay, skips if it has not yet appeared
			continue
		}

		if !waySheet.CloseDt.IsZero() {
			waySheetID := atoiSafe(id)
			info, err := r.loadWaySheetInfo(ctx, waySheetID)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed loading way sheet info for way sheet %d", waySheetID))
				logger.Logf(logger.ERROR, "FinanceRoutesReporter.processOpenedWaySheets()", "failed load way sheet info for way sheet %d: %v", waySheetID, err)
				continue
			}

			routeID := atoiSafe(info.Route.RouteCarID)
			shipmentID := ""
			if len(info.Shippings) > 0 {
				shipmentID = info.Shippings[0].ID
			}

			parking, err := r.loadParking(ctx, info.DstOffices)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed loading parking for route %d, way sheet %s", routeID, waySheet.WaySheetID))
				logger.Logf(logger.ERROR, "FinanceRoutesReporter.processOpenedWaySheets()", "failed loading parking info for way sheet %d: %v", waySheetID, err)
			}

			var driverName string
			if len(info.Drivers) > 0 {
				for _, driver := range info.Drivers {
					if driverName != "" {
						driverName += ", "
					}
					driverName += driver.DriverName
				}
			}

			var totalReturnTare, currentReturnTare int
			for _, tare := range info.Tares {
				if tare == nil {
					continue
				}
				if tare.IsReturn {
					totalReturnTare++
					if !tare.DtArrival.IsZero() {
						currentReturnTare++
					}
				}
			}

			salaryRate := 0.0
			if v, ok := r.salaryRate[routeID]; ok {
				salaryRate = v
			} else if v, ok := r.salaryRatePercent[routeID]; ok {
				salaryRate = waySheet.TotalPrice * (v / 100) // calculated if the rate is in percentages and not fixed
			} else {
				r.prompter.PromptError(fmt.Sprintf("There is no salary rate data for the route %d, way sheet %s", routeID, waySheet.WaySheetID))
				logger.Logf(logger.ERROR, "FinanceRoutesReporter.processOpenedWaySheets()", "there is no salary rate data for the route %d, way sheet %d", routeID, waySheetID)
			}

			totalPriceSubFine := waySheet.TotalPrice - waySheet.SumFine
			marriage := totalPriceSubFine * r.marriageRate
			tax := (totalPriceSubFine - marriage) * r.taxRate
			extendedSalaryRate := salaryRate + marriage + tax
			margin := totalPriceSubFine - (salaryRate + marriage + tax)

			r.prompter.PromptCloseWaySheet(routeID, waySheet.WaySheetID, shipmentID)
			logger.Logf(logger.INFO, "FinanceRoutesReporter.processOpenedWaySheets()", "way sheet %d is closed on route %d, shipment %s", waySheetID, routeID, shipmentID)
			err = r.sendReport(ctx, &reports.FinanceRoutesReportData{
				RouteID:            routeID,
				ShipmentID:         shipmentID,
				WaySheetID:         id,
				Parking:            parking,
				DateOpen:           info.DateOpen,
				DriverName:         driverName,
				VehicleNumberPlate: info.Vehicles.ShippingCarNumber,
				ShippedBarcodes:    info.TotalBarcodesCount,
				ShippedTare:        len(info.Tares),
				TotalReturnTare:    totalReturnTare,
				CurrentReturnTare:  currentReturnTare,
				Income:             waySheet.TotalPrice,
				IncomeReturn:       waySheet.SumReturn,
				Fine:               waySheet.SumFine,
				SalaryRate:         salaryRate,
				ExtendedSalaryRate: extendedSalaryRate,
				Marriage:           marriage,
				PercentMarriage:    r.percentMarriage,
				Tax:                tax,
				PercentTax:         r.percentTax,
				Margin:             margin,
			})
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed send report, route %d, shipment: %s, way sheet: %d", routeID, shipmentID, waySheetID))
				logger.Logf(logger.ERROR, "FinanceRoutesReporter.processOpenedWaySheets()", "failed send report, route %d, shipment: %s, way sheet: %d: %v", routeID, shipmentID, waySheetID, err)
				continue
			}
			delete(r.chOpenedWaySheets, id)
		}
	}

	return nil
}

func (r *FinanceRoutesReporter) isValidRoute(routeID, supplierID int) bool {
	if _, ok := r.suppliers[supplierID]; !ok {
		return false
	}
	return true
}

func (r *FinanceRoutesReporter) loadParking(ctx context.Context, dstOffices []*wb_models.WaySheetDestinationOffice) (int, error) {
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
		return 0, errors.Wrapf(err, "FinanceRoutesReporter.loadParking()", "failed load remains tares")
	}
	if len(remainsTares) == 0 {
		return 0, nil
	}

	_, parking := SpNameToGateParking(remainsTares[0].SpName)

	return parking, nil
}

func (r *FinanceRoutesReporter) loadRemainsTares(ctx context.Context, dstOfficeIDs []int) (tares []*wb_models.TareForOffice, err error) {
	err = retryAction(ctx, "FinanceRoutesReporter.loadRemainsTares", 3, 1*time.Second, func() error {
		// isDrive = false is default value
		tares, err = r.services.WBLogisticService.GetTaresForOffices(ctx, r.officeID, dstOfficeIDs, false)
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "FinanceRoutesReporter.loadRemainsTares()", "failed load remains tares for offices")
	}
	return tares, nil
}

func (r *FinanceRoutesReporter) loadWaySheets(ctx context.Context) (waySheets []*wb_models.WaySheet, err error) {
	now := time.Now()
	for supplierID := range r.suppliers {
		var page *wb_models.WaySheetsPage
		err = retryAction(ctx, "FinanceRoutesReporter.loadWaySheets", 3, 1*time.Second, func() error {
			page, err = r.services.WBLogisticService.GetWaySheets(ctx, &models.WBLogisticGetWaySheetsParamsRequest{
				DateOpen:    time.Date(now.Year(), now.Month(), now.Day()-3, 0, 0, 0, 0, time.UTC),
				DateClose:   now,
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
			logger.Logf(logger.ERROR, "FinanceRoutesReporter.loadWaySheets()", "failed load way sheets for supplier %d: %v", supplierID, err)
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
		return nil, errors.Wrap(err, "FinanceRoutesReporter.loadWaySheets()", "failed load way sheets")
	}
	return waySheets, nil
}

func (r *FinanceRoutesReporter) loadWaySheetInfo(ctx context.Context, waySheetID int) (waySheetInfo *wb_models.WaySheetInfo, err error) {
	err = retryAction(ctx, "FinanceRoutesReporter.loadWaySheetInfo", 3, 1*time.Second, func() error {
		waySheetInfo, err = r.services.WBLogisticService.GetWaySheetInfo(ctx, waySheetID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "FinanceRoutesReporter.loadWaySheetInfo()", "failed load way sheet %d info", waySheetID)
	}

	if waySheetInfo == nil {
		return nil, errors.Newf("FinanceRoutesReporter.loadWaySheetInfo()", "way sheet %d info returned empty value without error", waySheetID)
	}

	return waySheetInfo, nil
}

func (r *FinanceRoutesReporter) sendReport(ctx context.Context, reportData *reports.FinanceRoutesReportData) error {
	report, err := r.report.Render(reportData)
	if err != nil {
		return errors.Wrapf(err, "FinanceRoutesReporter.sendReport()", "failed render report, route id: %d shipment id: %s, way sheet id: %s", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID)
	}

	if r.isRenderTG {
		messages, err := r.rendererTG.Render(report)
		if err != nil {
			r.prompter.PromptError(fmt.Sprintf("Failed to render report for Telegram bot, route id: %d, way sheet id: %s", reportData.RouteID, reportData.WaySheetID))
			logger.Logf(logger.ERROR, "FinanceRoutesReporter.sendReport()", "failed render report for Telegram bot, route id: %d shipment id: %s, way sheet id: %s: %v", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID, err)
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
				r.prompter.PromptError(fmt.Sprintf("Failed to send report Telegram Bot, route id: %d way sheet id: %s", reportData.RouteID, reportData.WaySheetID))
				logger.Logf(logger.ERROR, "FinanceRoutesReporter.sendReport()", "failed send report to Telegram bot, route id: %d, shipment id: %s, way sheet id: %s: %v", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID, err)
			} else {
				r.prompter.PromptSendReport(reportData.RouteID, reportData.WaySheetID, reportData.ShipmentID)
				logger.Logf(logger.INFO, "FinanceRoutesReporter.sendReport()", "send report to Telegram bot, route id: %d, shipment id: %s, waysheet id: %s", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID)
			}
		}
	}

	return nil
}

func (r *FinanceRoutesReporter) sendTelegramBot(ctx context.Context) error {
	for r.messageQueueTG.Len() > 0 {
		err := retryAction(ctx, "FinanceRoutesReporter.sendTelegramBot", 3, 1*time.Second, func() error {
			message, ok := r.messageQueueTG.Peek()
			if !ok {
				return errors.New("FinanceRoutesReporter.sendTelegramBot()", "failed to get message from Telegram message queue")
			}
			return r.services.TelegramBotService.SendMessage(r.tgChatID, message, "HTML")
		})
		if err != nil {
			return errors.Wrapf(err, "FinanceRoutesReporter.sendTelegramBot()", "failed send data to chat %d", r.tgChatID)
		}
		r.messageQueueTG.Pop()
	}

	return nil
}

func (r *FinanceRoutesReporter) sendRemainingMessagesTelegramBot(ctx context.Context) error {
	if r.messageQueueTG.Len() <= 0 {
		return nil
	}

	err := r.sendTelegramBot(ctx)
	if err != nil {
		return errors.Wrap(err, "FinanceRoutesReporter.sendRemainingMessagesTelegramBot)", "failed sending remains messages to Telegram bot")
	}
	logger.Log(logger.INFO, "FinanceRoutesReporter.sendRemainingMessagesTelegramBot()", "send remains messages to Telegram bot")

	return nil
}
