package reporters

import (
	"context"
	"fmt"
	"time"

	wb_models "wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/models"
	"wb_logistic_assistant/internal/prompters"
	"wb_logistic_assistant/internal/report_renderers"
	"wb_logistic_assistant/internal/reports"
	"wb_logistic_assistant/internal/services"
	"wb_logistic_assistant/internal/storage"
)

type ShipmentCloseReporter struct {
	config   *config.Config
	storage  storage.Storage
	services *services.Container
	report   *reports.ShipmentCloseReport
	prompter prompters.ShipmentCloseReporterPrompter

	messageQueueTG        *Queue[string]
	counterErrMessageSend int
	limitErrMessageSend   int
	rendererGS            report_renderers.ReportRenderer[[][]interface{}]
	rendererTG            report_renderers.ReportRenderer[[]string]
	isRenderGS            bool
	isRenderTG            bool
	tgChatID              int64

	officeID                int
	suppliers               map[int]struct{} // supplier id -> struct{}
	skipRoutes              map[int]struct{} // route id -> struct{}
	intervalUpdateShipments time.Duration
	prevTimeUpdateShipments time.Time

	spreadsheetID string
	sheetName     string
	sheetPosition string

	chOpenedShipments map[int]int // shipment id -> route id
}

func NewShipmentCloseReporter(config *config.Config, storage storage.Storage, service *services.Container, prompter prompters.ShipmentCloseReporterPrompter) *ShipmentCloseReporter {
	return &ShipmentCloseReporter{
		config:   config,
		storage:  storage,
		services: service,
		prompter: prompter,
		report:   &reports.ShipmentCloseReport{},

		messageQueueTG:      New[string](300),
		limitErrMessageSend: 3,
		rendererGS:          &report_renderers.GoogleSheetsRenderer{},
		rendererTG:          &report_renderers.TelegramBotRenderer{Mode: report_renderers.TelegramBotRenderHTML},
		isRenderGS:          config.Reports().ShipmentClose().IsRenderGoogleSheets(),
		isRenderTG:          config.Reports().ShipmentClose().IsRenderTelegramBot(),
		tgChatID:            config.Telegram().ShipmentClose().ChatID(),

		officeID:                config.Logistic().Office().ID(),
		suppliers:               config.Logistic().Office().SuppliersMap(),
		skipRoutes:              config.Logistic().Office().SkipRoutesMap(),
		intervalUpdateShipments: config.Reports().ShipmentClose().IntervalUpdateShipments(),

		spreadsheetID: config.GoogleSheets().ReportSheets().GeneralRoutes().SpreadsheetID(),
		sheetName:     config.GoogleSheets().ReportSheets().GeneralRoutes().SheetName(),
		sheetPosition: "A1",

		chOpenedShipments: map[int]int{},
	}
}

func (r *ShipmentCloseReporter) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "ShipmentCloseReporter.Run()", "task was preliminarily completed")
	}

	r.prompter.PromptStart()

	timeStart := time.Now()

	if r.messageQueueTG.Len() > 0 {
		err := r.sendTelegramBot(ctx)
		if err != nil {
			r.prompter.PromptError("Failed sending remains messages to Telegram bot")
			return errors.Wrap(err, "ShipmentCloseReporter.Run()", "failed sending remains messages to Telegram bot")
		}
		logger.Log(logger.INFO, "ShipmentCloseReporter.Run()", "send remains messages to Telegram bot")
	}

	if timeStart.After(r.prevTimeUpdateShipments.Add(r.intervalUpdateShipments)) {
		if err := r.findOpenedShipments(ctx); err != nil {
			r.prompter.PromptError("Failed finding opened shipments")
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.Run()", "failed finding opened shipments: %v", err)
		}
		r.prevTimeUpdateShipments = timeStart
	}

	if err := r.processOpenedShipments(ctx); err != nil {
		r.prompter.PromptError("Failed processing opened shipments")
		logger.Logf(logger.ERROR, "ShipmentCloseReporter.Run()", "failed processing opened shipments: %v", err)
	}

	r.prompter.PromptFinish(time.Since(timeStart))
	return nil
}

func (r *ShipmentCloseReporter) findOpenedShipments(ctx context.Context) error {
	logger.Log(logger.INFO, "ShipmentCloseReporter.findOpenedShipments()", "start find opened shipments")

	remainsReport, err := r.loadRemainsLastMileReport(ctx)
	if err != nil {
		return errors.Wrap(err, "ShipmentCloseReporter.findOpenedShipments()", "failed load remains last miles report")
	}

	for _, route := range remainsReport.Routes {
		if !r.isValidRoute(route) {
			continue
		}

		routeID := route.CarID

		routeInfo, err := r.loadRemainsLastMileReportRouteInfo(ctx, routeID)
		if routeInfo == nil || err != nil {
			r.prompter.PromptError(fmt.Sprintf("Failed load route info for route %d", routeID))
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.findOpenedShipments()", "failed load route info for route %d: %v", routeID, err)
			continue
		}
		if len(routeInfo) == 0 || routeInfo[0].DestinationOfficeName == "" {
			continue
		}

		shipments, err := r.loadShipments(ctx, routeID, routeInfo[0].DestinationOfficeName)
		if err != nil {
			r.prompter.PromptError(fmt.Sprintf("Failed load shipments for route %d", routeID))
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.findOpenedShipments()", "failed load shipments for route %d: %v", routeID, err)
			continue
		}
		if len(shipments) == 0 {
			continue
		}

		// Checking all shipments for opening. Here information about open shipments is collected, which will then be tracked separately.
		for _, shipment := range shipments {
			if !shipment.CloseDt.IsZero() {
				continue
			}
			r.chOpenedShipments[shipment.ShipmentID] = routeID
			r.prompter.PromptShipmentOpened(routeID, shipment.ShipmentID, len(r.chOpenedShipments))
			logger.Logf(logger.INFO, "ShipmentCloseReporter.findOpenedShipments()", "open shipment: route %d, shipment %d. Cache size opened shipments: %d", routeID, shipment.ShipmentID, len(r.chOpenedShipments))
		}
	}
	return nil
}

func (r *ShipmentCloseReporter) processOpenedShipments(ctx context.Context) error {
	logger.Log(logger.INFO, "ShipmentCloseReporter.processOpenedShipments()", "start process opened shipments")

	for shipmentID, routeID := range r.chOpenedShipments {
		time.Sleep(100 * time.Millisecond)

		info, err := r.loadShipmentInfo(ctx, shipmentID)
		if err != nil {
			r.prompter.PromptError(fmt.Sprintf("Failed loading shipment info for shipment %d", shipmentID))
			delete(r.chOpenedShipments, shipmentID)
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.processOpenedShipments()", "failed load shipment info, route %d, shipment %d: %v", routeID, shipmentID, err)
			continue
		}

		if !info.CloseDt.IsZero() {
			// Transfer boxes
			transferBoxes, err := r.loadShipmentTransfersBoxes(ctx, shipmentID)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed loading shipment transfers boxes for shipment %d", shipmentID))
				logger.Logf(logger.ERROR, "ShipmentCloseReporter.processOpenedShipments()", "failed load shipment transfer boxes, shipment %d: %v", shipmentID, err)
			}

			totalTransferBarcodes := 0
			for _, box := range transferBoxes {
				if box == nil {
					logger.Logf(logger.WARN, "ShipmentCloseReporter.processOpenedShipments()", "transfer box is nil, shipment %d", shipmentID)
					continue
				}
				totalTransferBarcodes += box.CountBarcodes
			}

			// Remains tares
			dstOfficesIDs := make([]int, len(info.DestinationOfficesInfo))
			for i, v := range info.DestinationOfficesInfo {
				if v == nil {
					continue
				}
				dstOfficesIDs[i] = v.DstOfficeID
			}

			remainsTares, err := r.loadRemainsTares(ctx, dstOfficesIDs)
			if err != nil {
				remainsTares = make([]*wb_models.TareForOffice, 0)
				r.prompter.PromptError(fmt.Sprintf("Failed loading remains tares for shipment %d", shipmentID))
				logger.Logf(logger.ERROR, "ShipmentCloseReporter.processOpenedShipments()", "failed load remains tares, shipment %d: %v", shipmentID, err)
			}

			remainsTaresInfo := make([]*reports.ShipmentCloseRemainsTareInfo, len(remainsTares))
			totalRemainsBarcodes := 0
			for i, t := range remainsTares {
				if t == nil {
					logger.Logf(logger.WARN, "ShipmentCloseReporter.processOpenedShipments()", "remains tare is nil for shipment %d, total tares %d", shipmentID, len(remainsTares))
					remainsTaresInfo[i] = &reports.ShipmentCloseRemainsTareInfo{}
					continue
				}

				remainsTaresInfo[i] = &reports.ShipmentCloseRemainsTareInfo{
					ID:              t.ID,
					DstOfficeID:     t.DstOfficeID,
					DstOfficeName:   t.DstOfficeName,
					CountBarcodes:   t.CountBarcodes,
					LastOperationDt: t.LastOperationDt,
				}
				totalRemainsBarcodes += t.CountBarcodes
			}

			spName := ""
			parking := 0
			if len(remainsTares) > 0 {
				spName = remainsTares[0].SpName
				_, parking = SpNameToGateParking(spName)
			}

			r.prompter.PromptShipmentClose(routeID, shipmentID)
			logger.Logf(logger.INFO, "ShipmentCloseReporter.processOpenedShipments()", "shipment %d is closed on route %d, waysheet %d", shipmentID, routeID, info.WaySheetID)
			err = r.sendReport(ctx, &reports.ShipmentCloseReportData{
				RouteID:               routeID,
				ShipmentID:            shipmentID,
				WaySheetID:            info.WaySheetID,
				Parking:               parking,
				DriverName:            info.DriverName,
				VehicleNumberPlate:    info.VehicleNumberPlate,
				TotalRemainsBarcodes:  totalRemainsBarcodes,
				TotalRemainsTares:     len(remainsTares),
				TotalTransferBarcodes: totalTransferBarcodes,
				TotalTransferTares:    len(transferBoxes),
				Date:                  info.CreateDt,
				DateCreate:            info.CreateDt,
				DateClose:             info.CloseDt,
				SpName:                spName,
				RemainsTaresInfo:      remainsTaresInfo,
			})
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed send report on route %d, shipment: %d", routeID, shipmentID))
				logger.Logf(logger.ERROR, "ShipmentCloseReporter.processOpenedShipments()", "failed send report on route %d, shipment %d: %v", routeID, shipmentID, err)
				continue
			}
			delete(r.chOpenedShipments, shipmentID)
		}
	}
	return nil
}

func (r *ShipmentCloseReporter) isValidRoute(route *wb_models.Route) bool {
	if route == nil || len(route.Suppliers) == 0 {
		return false
	}

	isValidSupplier := false
	for _, supplier := range route.Suppliers {
		if _, ok := r.suppliers[supplier.ID]; ok {
			isValidSupplier = true
			break
		}
	}
	if !isValidSupplier {
		return false
	}

	if _, ok := r.skipRoutes[route.CarID]; ok {
		return false
	}
	return true
}

func (r *ShipmentCloseReporter) loadRemainsLastMileReport(ctx context.Context) (remainsLastMileReport *wb_models.RemainsLastMileReport, err error) {
	err = retryAction(ctx, "ShipmentCloseReporter.loadRemainsLastMileReport", 3, 1*time.Second, func() error {
		remainsLastMileReport, err = r.services.WBLogisticService.GetRemainsLastMileReportByOfficeID(ctx, r.officeID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "ShipmentCloseReporter.loadRemainsLastMileReport()", "failed load remains last mile report for office %d", r.officeID)
	}
	return remainsLastMileReport, nil
}

func (r *ShipmentCloseReporter) loadRemainsLastMileReportRouteInfo(ctx context.Context, routeID int) (remainsLastMileReportsRouteInfo []*wb_models.RemainsLastMileReportsRouteInfo, err error) {
	err = retryAction(ctx, "ShipmentCloseReporter.loadRemainsLastMileReportRouteInfo", 3, 1*time.Second, func() error {
		remainsLastMileReportsRouteInfo, err = r.services.WBLogisticService.GetRemainsLastMileReportsRouteInfo(ctx, routeID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "ShipmentCloseReporter.loadRemainsLastMileReportRouteInfo()", "failed load remains last miles reports route %d info", routeID)
	}

	if remainsLastMileReportsRouteInfo == nil || len(remainsLastMileReportsRouteInfo) == 0 {
		return nil, errors.Newf("ShipmentCloseReporter.loadRemainsLastMileReportRouteInfo()", "remains last miles reports route %d info is empty", routeID)
	}

	return remainsLastMileReportsRouteInfo, nil
}

func (r *ShipmentCloseReporter) loadShipments(ctx context.Context, routeID int, destinationAddressName string) (shipments []*wb_models.Shipment, err error) {
	now := time.Now()
	for supplierID := range r.suppliers {
		var res []*wb_models.Shipment
		err = retryAction(ctx, "ShipmentCloseReporter.loadShipments", 3, 1*time.Second, func() error {
			res, _, err = r.services.WBLogisticService.GetShipments(ctx, &models.WBLogisticGetShipmentsParamsRequest{
				DataStart:       now.AddDate(0, 0, -2),
				DataEnd:         now,
				SrcOfficeID:     r.officeID,
				PageIndex:       0,
				Limit:           3,
				SupplierID:      supplierID,
				Direction:       -1,
				Sorter:          "updated_at",
				FilterDstOffice: destinationAddressName,
			})
			return err
		})
		if err != nil {
			r.prompter.PromptError(fmt.Sprintf("failed load way sheets for supplier %d", supplierID))
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.loadShipments()", "failed load shipments on route %d for supplier %d: %v", routeID, supplierID, err)
			continue
		}
		if len(res) > 0 {
			if len(shipments) == 0 {
				shipments = res
			} else {
				shipments = append(shipments, res...)
			}
		}
	}
	if err != nil {
		return nil, errors.Wrapf(err, "ShipmentCloseReporter.loadShipments()", "failed load shipments on route %d", routeID)
	}
	return shipments, nil
}

func (r *ShipmentCloseReporter) loadShipmentInfo(ctx context.Context, shipmentID int) (shipmentInfo *wb_models.ShipmentInfo, err error) {
	err = retryAction(ctx, "ShipmentCloseReporter.loadShipmentInfo()", 3, 1*time.Second, func() error {
		shipmentInfo, err = r.services.WBLogisticService.GetShipmentInfo(ctx, shipmentID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "ShipmentCloseReporter.loadShipmentInfo()", "failed load shipment %d info", shipmentID)
	}

	if shipmentInfo == nil {
		return nil, errors.Newf("ShipmentCloseReporter.loadShipmentInfo()", "shipment %d info returned empty value without error", shipmentID)
	}

	if shipmentInfo.DestinationOfficesInfo == nil {
		shipmentInfo.DestinationOfficesInfo = make([]*wb_models.ShipmentInfoDestinationOfficeInfo, 0)
	}

	return shipmentInfo, nil
}

func (r *ShipmentCloseReporter) loadRemainsTares(ctx context.Context, dstOfficeIDs []int) (tares []*wb_models.TareForOffice, err error) {
	if len(dstOfficeIDs) == 0 {
		return nil, errors.New("ShipmentCloseReporter.loadRemainsTares()", "failed load tares for offices, dst office ids is empty")
	}
	err = retryAction(ctx, "ShipmentCloseReporter.loadRemainsTares", 3, 1*time.Second, func() error {
		// isDrive = false is default value
		tares, err = r.services.WBLogisticService.GetTaresForOffices(ctx, r.officeID, dstOfficeIDs, false)
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "ShipmentCloseReporter.loadRemainsTares()", "failed load tares for offices")
	}
	return tares, nil
}

func (r *ShipmentCloseReporter) loadShipmentTransfersBoxes(ctx context.Context, shipmentID int) ([]*wb_models.ShipmentTransferBox, error) {
	var err error
	var transfers *wb_models.ShipmentTransfers
	err = retryAction(ctx, "ShipmentCloseReporter.loadShipmentTransfers", 3, 1*time.Second, func() error {
		transfers, err = r.services.WBLogisticService.GetShipmentTransfers(ctx, shipmentID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "ShipmentCloseReporter.loadShipmentTransfers()", "failed load transfers of shipment %d", shipmentID)
	}

	if transfers.TransferBoxes == nil {
		return nil, errors.Newf("ShipmentCloseReporter.loadShipmentTransfers()", "failed load transfers boxes, boxes is empty of shipment %d", shipmentID)
	}

	return transfers.TransferBoxes, nil
}

func (r *ShipmentCloseReporter) sendReport(ctx context.Context, reportData *reports.ShipmentCloseReportData) error {
	report, err := r.report.Render(reportData)
	if err != nil {
		r.prompter.PromptError("Failed to render report")
		return errors.Wrapf(err, "ShipmentCloseReporter.sendReport()", "failed render report, route id: %d shipment id: %d, way sheet id: %d", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID)
	}

	if r.isRenderGS {
		data, err := r.rendererGS.Render(report)
		if err != nil {
			r.prompter.PromptError("Failed to render report Google Sheet")
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.sendReport()", "failed render report for Google Sheets, route id: %d shipment id: %d, way sheet id: %d: %v", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID, err)
		} else {
			if err = r.sendGoogleSheets(ctx, data); err != nil {
				r.prompter.PromptError("Failed to send report Google Sheet")
				logger.Logf(logger.ERROR, "ShipmentCloseReporter.sendReport()", "failed send report to Google Sheets, route id: %d shipment id: %d, way sheet id: %d: %v", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID, err)
			} else {
				r.prompter.PromptSendReport(reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID)
				logger.Logf(logger.INFO, "ShipmentCloseReporter.sendReport()", "send report to Google Sheets, route id: %d, shipment id: %d, waysheet id: %d", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID)
			}
		}
	}

	if r.isRenderTG {
		messages, err := r.rendererTG.Render(report)
		if err != nil {
			r.prompter.PromptError("failed to render report Telegram Bot")
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.sendReport()", "failed render report for Telegram Bot, route id: %d shipment id: %d, way sheet id: %d: %v", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID, err)
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
				logger.Logf(logger.ERROR, "ShipmentCloseReporter.sendReport()", "failed send report to Telegram Bot, route id: %d shipment id: %d, way sheet id: %d: %v", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID, err)
			} else {
				r.prompter.PromptSendReport(reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID)
				logger.Logf(logger.INFO, "ShipmentCloseReporter.sendReport()", "send report to Telegram Bot, route id: %d, shipment id: %d, waysheet id: %d", reportData.RouteID, reportData.ShipmentID, reportData.WaySheetID)
			}
		}
	}

	return nil
}

func (r *ShipmentCloseReporter) sendGoogleSheets(ctx context.Context, data [][]interface{}) error {
	err := retryAction(ctx, "ShipmentCloseReporter.sendGoogleSheets", 3, 1*time.Second, func() error {
		err := r.services.GoogleSheetsService.ClearValues(r.spreadsheetID, r.sheetName, "A:Z")
		if err != nil {
			return errors.Wrapf(err, "ShipmentCloseReporter.sendGoogleSheets()", "failed clear sheet %s, page %s", r.spreadsheetID, r.sheetName)
		}
		return r.services.GoogleSheetsService.UpdateValues(r.spreadsheetID, r.sheetName, r.sheetPosition, data, false)
	})
	if err != nil {
		return errors.Wrapf(err, "ShipmentCloseReporter.sendGoogleSheets()", "failed update sheet %s, page %s to position %s", r.spreadsheetID, r.sheetName, r.sheetPosition)
	}
	return nil
}

func (r *ShipmentCloseReporter) sendTelegramBot(ctx context.Context) error {
	if r.messageQueueTG.Len() <= 0 {
		return nil
	}

	for r.messageQueueTG.Len() > 0 {
		err := retryAction(ctx, "ShipmentCloseReporter.sendTelegramBot", 3, 1*time.Second, func() error {
			message, ok := r.messageQueueTG.Peek()
			if !ok {
				return errors.New("ShipmentCloseReporter.sendTelegramBot()", "failed to get message from Telegram message queue")
			}
			return r.services.TelegramBotService.SendMessage(r.tgChatID, message, "HTML")
		})
		if err != nil {
			r.counterErrMessageSend++
			// if an error occurs in any of the messages or Telegram refuses to accept the message, it is better to reset the message queue
			if r.counterErrMessageSend >= r.limitErrMessageSend {
				r.messageQueueTG.Reset()
				r.counterErrMessageSend = 0
			}
			return errors.Wrapf(err, "ShipmentCloseReporter.sendTelegramBot()", "failed send data to chat %d", r.tgChatID)
		}

		r.counterErrMessageSend = 0
		r.messageQueueTG.Pop()
	}

	return nil
}
