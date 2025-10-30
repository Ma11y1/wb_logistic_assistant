package reporters

import (
	"context"
	"fmt"
	"sort"
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

type GeneralRoutesReporter struct {
	config                             *config.Config
	storage                            storage.Storage
	services                           *services.Container
	generalRoutesReport                *reports.GeneralRoutesReport
	shipmentCloseReport                *reports.ShipmentCloseReport
	shipmentCloseReportDataRenderQueue []*reports.ShipmentCloseReportData
	googleSheetsRenderer               report_renderers.ReportRenderer[[][]interface{}]
	telegramBotRenderer                report_renderers.ReportRenderer[[]string]
	prompter                           prompters.GeneralRoutesReporterPrompter
	isShipmentCloseReport              bool
	isNeedRenderReport                 bool
	isNeedProcessShipments             bool

	// params
	parking         map[int]int      // [route id]parking
	suppliers       map[int]struct{} // [supplier id]struct{}
	skipRoutes      map[int]struct{} // [route id]struct{}
	sortColumn      int
	isSort          bool
	isSortAscending bool

	// google sheets
	spreadsheetID string
	sheetName     string
	sheetPosition string

	// telegram
	telegramChatID int64

	// report
	routesMetaData *reports.GeneralRoutesReportMetaData // general routes report
	routesData     []*reports.GeneralRoutesReportData   // general routes report

	// cache
	cacheRoutesData                      map[int]*reports.GeneralRoutesReportData // [route.CarID]route data
	cacheRemainsLastMileReport           *wb_models.RemainsLastMileReport
	cacheRemainsLastMileReportsRouteInfo map[int][]*wb_models.RemainsLastMileReportsRouteInfo
	cacheShipments                       map[int]*cacheWithTime[[]*wb_models.Shipment]
	cacheLastShipments                   map[int]*wb_models.Shipment
	cacheJobsScheduling                  *wb_models.JobsScheduling

	// timings
	tllDurationResetCountChangeBarcodes  time.Duration
	ttlDurationLoadRemainsLastMileReport time.Duration
	ttlDurationLoadJobsScheduling        time.Duration
	ttlDurationLoadShipments             time.Duration
	timeResetCountChangeBarcodes         time.Time
	timeLastResetCountChangeBarcodes     time.Time
	timeLoadRemainsLastMileReport        time.Time
	timeLoadJobsScheduling               time.Time
}

func NewGeneralRoutesReporter(config *config.Config, storage storage.Storage, service *services.Container, prompter prompters.GeneralRoutesReporterPrompter) *GeneralRoutesReporter {
	grcfg := config.Reports().GeneralRoutes()
	suppliers := make(map[int]struct{})
	for _, supplier := range config.Logistic().Office().Suppliers() {
		suppliers[supplier] = struct{}{}
	}
	skipRoutes := make(map[int]struct{})
	for _, route := range config.Reports().GeneralRoutes().SkipRoutes() {
		skipRoutes[route] = struct{}{}
	}
	return &GeneralRoutesReporter{
		config:                               config,
		storage:                              storage,
		services:                             service,
		generalRoutesReport:                  reports.NewGeneralRoutesReport(),
		shipmentCloseReport:                  reports.NewShipmentCloseReport(),
		shipmentCloseReportDataRenderQueue:   make([]*reports.ShipmentCloseReportData, 0),
		googleSheetsRenderer:                 &report_renderers.GoogleSheetsRenderer{},
		telegramBotRenderer:                  &report_renderers.TelegramBotRenderer{},
		prompter:                             prompter,
		spreadsheetID:                        config.GoogleSheets().ReportSheets().GeneralRoutes().SpreadsheetID(),
		sheetName:                            config.GoogleSheets().ReportSheets().GeneralRoutes().SheetName(),
		sheetPosition:                        "A1",
		routesMetaData:                       &reports.GeneralRoutesReportMetaData{},
		routesData:                           make([]*reports.GeneralRoutesReportData, 0),
		cacheRoutesData:                      map[int]*reports.GeneralRoutesReportData{},
		cacheRemainsLastMileReportsRouteInfo: map[int][]*wb_models.RemainsLastMileReportsRouteInfo{},
		cacheShipments:                       map[int]*cacheWithTime[[]*wb_models.Shipment]{},
		cacheLastShipments:                   map[int]*wb_models.Shipment{},
		telegramChatID:                       config.Telegram().ChatID(),
		tllDurationResetCountChangeBarcodes:  grcfg.TTLResetChangeBarcodes(),
		ttlDurationLoadRemainsLastMileReport: grcfg.TTLLoadRemainsLastMileReport(),
		ttlDurationLoadJobsScheduling:        grcfg.TTLLoadJobsScheduling(),
		ttlDurationLoadShipments:             grcfg.TTLLoadShipments(),
		parking:                              config.Logistic().Office().Parking(),
		suppliers:                            suppliers,
		skipRoutes:                           skipRoutes,
		sortColumn:                           grcfg.SortColumn(),
		isSort:                               grcfg.IsSort(),
		isSortAscending:                      grcfg.IsAscending(),
	}
}

func (r *GeneralRoutesReporter) Run(ctx context.Context) error {
	err := r.loadRemainsLastMileReport(ctx)
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.Run()", "failed load remains last miles reports")
	}

	//err = r.loadJobsScheduling(ctx)
	//if err != nil {
	//	logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed load jobs scheduling: %v", err)
	//}

	err = r.process(ctx)
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.Run()", "failed process data")
	}

	if !r.isNeedRenderReport {
		return nil
	}

	// General routes report
	generalRoutesReport, err := r.generalRoutesReport.Render(r.routesMetaData, r.routesData)
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.Run()", "failed render general routes report")
	}

	generalRoutesReportData, err := r.googleSheetsRenderer.Render(generalRoutesReport)
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.Run()", "failed render general routes report for google sheets")
	}

	err = r.sendGoogleSheets(generalRoutesReportData)
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.Run()", "failed send general routes report to google sheets")
	}

	// Shipment close report
	if r.isShipmentCloseReport {
		for _, value := range r.shipmentCloseReportDataRenderQueue {
			if ctx.Err() != nil {
				return errors.Wrap(ctx.Err(), "GeneralRoutesReporter.Run()", "sending a report to the Telegram bot has been suspended")
			}

			shipmentCloseReport, err := r.shipmentCloseReport.Render(value)
			if err != nil {
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed render shipment close report: %v", err)
				continue
			}

			shipmentCloseReportData, err := r.telegramBotRenderer.Render(shipmentCloseReport)
			if err != nil {
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed render shipment close report for telegram bot: %v", err)
				continue
			}

			fmt.Println(shipmentCloseReportData)
			//err = r.sendTelegramBot(shipmentCloseReportData)
			//if err != nil {
			//	logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed send shipment close report to telegram bot: %v", err)
			//	continue
			//}
			time.Sleep(1 * time.Second)
		}

		if len(r.shipmentCloseReportDataRenderQueue) > 0 {
			r.shipmentCloseReportDataRenderQueue = make([]*reports.ShipmentCloseReportData, 0)
		}
	}

	r.prompter.PromptRender()
	r.isNeedRenderReport = false

	return nil
}

func (r *GeneralRoutesReporter) process(ctx context.Context) error {
	if !r.isNeedRenderReport {
		return nil
	}

	now := time.Now()
	routes := r.cacheRemainsLastMileReport.Routes
	countRoutes := len(routes)

	if cap(r.routesData) < countRoutes {
		r.routesData = make([]*reports.GeneralRoutesReportData, countRoutes)
	} else {
		r.routesData = r.routesData[:countRoutes]
	}

	isNeedResetCountChangeBarcode := false
	if r.timeResetCountChangeBarcodes.IsZero() || now.After(r.timeResetCountChangeBarcodes) {
		r.timeResetCountChangeBarcodes = now.Add(r.tllDurationResetCountChangeBarcodes)
		r.timeLastResetCountChangeBarcodes = now
		r.routesMetaData.TimeLastResetCountBarcodes = r.timeLastResetCountChangeBarcodes
		isNeedResetCountChangeBarcode = true
	}

	indexRouteData := 0
	// Used separately if route object is nil and the loop index is out of range
	for i := 0; i < countRoutes; i++ {
		route := routes[i]
		if !r.isValidRoute(route) {
			continue
		}
		routeID := route.CarID
		parking, ok := r.parking[routeID]
		if !ok {
			parking = -1
		}

		// Route data
		routeData := r.cacheRoutesData[routeID]
		if routeData == nil {
			routeData = &reports.GeneralRoutesReportData{}
			r.cacheRoutesData[routeID] = routeData
		}
		r.routesData[indexRouteData] = routeData
		indexRouteData++

		routeData.ID = routeID
		routeData.Parking = parking

		routeData.CountTares = route.CountTares
		routeData.Barcodes = route.CountShk

		if isNeedResetCountChangeBarcode {
			routeData.CountChangeBarcodes = 0
		}
		if routeData.LastBarcodes > 0 {
			routeData.CountChangeBarcodes = routeData.CountChangeBarcodes + (route.CountShk - routeData.LastBarcodes)
		}
		routeData.LastBarcodes = route.CountShk

		////  Remains last mile route info
		err := r.loadRemainsLastMileReportsRouteInfo(ctx, routeID)
		if err != nil {
			// Without interrupting the entire loop, if an error returns, there is no point in continuing without this data
			logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed load remains last mile reports route info: %v", err)
			continue
		}

		remainsLastMileReportsRouteInfo := r.cacheRemainsLastMileReportsRouteInfo[routeID]
		if remainsLastMileReportsRouteInfo == nil {
			logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "after loading remains last mile reports route info, it is missing for route %d", route.CarID)
			continue
		}

		// To search for shipments on the current route, you need to know at least one destination address
		destinationOfficeName := remainsLastMileReportsRouteInfo[0].DestinationOfficeName
		if destinationOfficeName == "" {
			continue
		}

		////  Shipments
		err = r.loadShipments(ctx, routeID, route.Suppliers[0].ID, destinationOfficeName)
		if err != nil {
			logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed load shipments: %v", err)
			continue
		}

		if r.isNeedProcessShipments {
			shipments := r.cacheShipments[routeID]
			if shipments == nil || shipments.value == nil {
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "after loadin shipments, it is missing for route %d", routeID)
				continue
			}

			if len(shipments.value) == 0 {
				logger.Logf(logger.WARN, "GeneralRoutesReporter.Run()", "after loadin shipments, it is empty for route %d", routeID)
				continue
			}

			lastShipment := r.cacheLastShipments[routeID]
			currentShipment := shipments.value[0]
			r.cacheLastShipments[routeID] = currentShipment

			routeData.ShipmentID = currentShipment.ShipmentID
			routeData.ShipmentDateCreate = currentShipment.CreateDt
			routeData.ShipmentDateClose = currentShipment.CloseDt

			//Checking if the last shipment has been closed, and if it has been closed or a new shipment has been opened, the condition is met
			if lastShipment != nil {
				for _, shipment := range shipments.value {
					if (shipment.ShipmentID == lastShipment.ShipmentID) && (lastShipment.CloseDt.IsZero() && !shipment.CloseDt.IsZero()) {
						shipmentsCloseRoutesInfo := make([]*reports.ShipmentCloseRemainsTareInfo, len(remainsLastMileReportsRouteInfo))

						remainsBarcodes := 0

						//for k, remainsLastMileReportRouteInfo := range remainsLastMileReportsRouteInfo {
						//	if remainsLastMileReportRouteInfo == nil {
						//		shipmentsCloseRoutesInfo[i] = &reports.ShipmentCloseRemainsTareInfo{}
						//		continue
						//	}
						//	shipmentsCloseRouteInfoTare := make([]*reports.ShipmentCloseReportDataRouteInfoTare, len(remainsLastMileReportRouteInfo.Tares))
						//	for j, tare := range remainsLastMileReportRouteInfo.Tares {
						//		if tare == nil {
						//			shipmentsCloseRouteInfoTare[j] = &reports.ShipmentCloseReportDataRouteInfoTare{}
						//			continue
						//		}
						//
						//		shipmentsCloseRouteInfoTare[j] = &reports.ShipmentCloseReportDataRouteInfoTare{
						//			ID:            tare.ID,
						//			CountBarcodes: tare.CountBarcodes,
						//			PrepareDate:   tare.PrepareDate,
						//		}
						//
						//		remainsBarcodes += tare.CountBarcodes
						//	}
						//	shipmentsCloseRoutesInfo[k] = &reports.ShipmentCloseRemainsTareInfo{
						//		DstOfficeID:   remainsLastMileReportRouteInfo.DestinationOfficeID,
						//		DstOfficeName: remainsLastMileReportRouteInfo.DestinationOfficeName,
						//		Tares:         shipmentsCloseRouteInfoTare,
						//	}
						//}

						routeData.RemainsBarcodes = remainsBarcodes

						if r.isShipmentCloseReport {
							r.shipmentCloseReportDataRenderQueue = append(r.shipmentCloseReportDataRenderQueue, &reports.ShipmentCloseReportData{
								RouteID:              routeID,
								ShipmentID:           shipment.ShipmentID,
								Parking:              parking,
								DriverName:           shipment.DriverName,
								TotalRemainsBarcodes: remainsBarcodes,
								Date:                 shipment.CreateDt,
								TimeCreate:           shipment.CreateDt,
								TimeClose:            shipment.CloseDt,
								RemainsTaresInfo:     shipmentsCloseRoutesInfo,
							})
						}
					}
				}
			}

			r.isNeedProcessShipments = false
		}

		//// Delay for loading
		time.Sleep(100 * time.Millisecond)
	}
	// if there are fewer routes than expected, shorten them so that there are no zero elements
	r.routesData = r.routesData[:indexRouteData]

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// rating
	if r.cacheJobsScheduling != nil && len(r.cacheJobsScheduling.Route) > 0 {
		jobsSchedulingRoutes := r.cacheJobsScheduling.Route
		for i := 0; i < len(jobsSchedulingRoutes); i++ {
			route := jobsSchedulingRoutes[i]
			if route == nil {
				continue
			}
			// pointers data objects are already in the data slice, so we take the object from there
			routeData := r.cacheRoutesData[route.RouteID]
			if routeData == nil {
				continue
			}
			if route.Rating != nil {
				routeData.Rating = route.Rating.OverallRating
			}
		}
	}

	// sort
	if r.isSort && len(r.routesData) > 0 {
		r.sort()
	}

	return nil
}

func (r *GeneralRoutesReporter) isValidRoute(route *wb_models.Route) bool {
	if route == nil || len(route.Suppliers) == 0 {
		return false
	}
	_, ok := r.suppliers[route.Suppliers[0].ID]
	if !ok {
		return false
	}
	_, ok = r.skipRoutes[route.CarID]
	if ok {
		return false
	}
	return true
}

func (r *GeneralRoutesReporter) sort() {
	sortType := SortType(r.sortColumn)
	sort.Slice(r.routesData, func(i, j int) bool {
		switch sortType {
		case IDSortType:
			if r.isSortAscending {
				return r.routesData[i].ID < r.routesData[j].ID
			}
			return r.routesData[i].ID > r.routesData[j].ID
		case ParkingSortType:
			if r.isSortAscending {
				return r.routesData[i].Parking < r.routesData[j].Parking
			}
			return r.routesData[i].Parking > r.routesData[j].Parking
		case TaresSortType:
			if r.isSortAscending {
				return r.routesData[i].CountTares < r.routesData[j].CountTares
			}
			return r.routesData[i].CountTares > r.routesData[j].CountTares
		case BarcodesSortType:
			if r.isSortAscending {
				return r.routesData[i].Barcodes < r.routesData[j].Barcodes
			}
			return r.routesData[i].Barcodes > r.routesData[j].Barcodes
		case RatingSortType:
			if r.isSortAscending {
				return r.routesData[i].Rating < r.routesData[j].Rating
			}
			return r.routesData[i].Rating > r.routesData[j].Rating
		default:
			return r.routesData[i].ID < r.routesData[j].ID
		}
	})

}

func (r *GeneralRoutesReporter) loadRemainsLastMileReport(ctx context.Context) error {
	var err error
	now := time.Now()

	if r.cacheRemainsLastMileReport != nil && now.Before(r.timeLoadRemainsLastMileReport) {
		return nil
	}

	var remainsLastMileReport *wb_models.RemainsLastMileReport
	err = retryAction(ctx, "GeneralRoutesReporter.loadRemainsLastMileReport", 3, 1*time.Second, func() error {
		remainsLastMileReport, err = r.services.WBLogisticService.GetRemainsLastMileReportByOfficeID(ctx, r.config.Logistic().Office().ID())
		return err
	})
	if err != nil {
		return err
	}

	if remainsLastMileReport == nil || remainsLastMileReport.Routes == nil {
		return errors.New("GeneralRoutesReporter.loadRemainsLastMileReport()", "remains last miles report is empty")
	}

	r.cacheRemainsLastMileReport = remainsLastMileReport
	r.timeLoadRemainsLastMileReport = now.Add(r.ttlDurationLoadRemainsLastMileReport)
	r.isNeedRenderReport = true

	return nil
}

func (r *GeneralRoutesReporter) loadRemainsLastMileReportsRouteInfo(ctx context.Context, routeID int) error {
	r.cacheRemainsLastMileReportsRouteInfo[routeID] = nil

	var err error
	var remainsLastMileReportsRouteInfo []*wb_models.RemainsLastMileReportsRouteInfo
	err = retryAction(ctx, "GeneralRoutesReporter.loadRemainsLastMileReportRouteInfo", 3, 1*time.Second, func() error {
		remainsLastMileReportsRouteInfo, err = r.services.WBLogisticService.GetRemainsLastMileReportsRouteInfo(ctx, routeID)
		return err
	})
	if err != nil {
		return err
	}

	if remainsLastMileReportsRouteInfo == nil || len(remainsLastMileReportsRouteInfo) == 0 {
		return errors.New("GeneralRoutesReporter.loadRemainsLastMileReport()", "remains last miles reports route info is empty")
	}

	r.cacheRemainsLastMileReportsRouteInfo[routeID] = remainsLastMileReportsRouteInfo
	r.isNeedRenderReport = true

	return nil
}

func (r *GeneralRoutesReporter) loadJobsScheduling(ctx context.Context) error {
	var err error
	now := time.Now()

	if r.cacheJobsScheduling != nil && now.Before(r.timeLoadJobsScheduling) {
		return nil
	}

	var jobsScheduling *wb_models.JobsScheduling
	err = retryAction(ctx, "GeneralRoutesReporter.loadJobsScheduling", 3, 1*time.Second, func() error {
		jobsScheduling, err = r.services.WBLogisticService.GetJobsScheduling(ctx)
		return err
	})
	if err != nil {
		return err
	}

	if jobsScheduling == nil || jobsScheduling.Route == nil {
		return errors.New("GeneralRoutesReporter.loadJobsScheduling()", "jobs scheduling is empty")
	}

	// checking for rating data, because the service may return zero values
	if len(jobsScheduling.Route) > 2 {
		if (jobsScheduling.Route[0].Rating == nil || jobsScheduling.Route[0].Rating.OverallRating == 0) ||
			(jobsScheduling.Route[1].Rating == nil || jobsScheduling.Route[1].Rating.OverallRating == 0) {
			return errors.New("GeneralRoutesReporter.loadJobsScheduling()", "jobs scheduling rating is empty")
		}
	}

	r.cacheJobsScheduling = jobsScheduling
	r.timeLoadJobsScheduling = now.Add(r.ttlDurationLoadJobsScheduling)
	r.isNeedRenderReport = true

	return nil
}

func (r *GeneralRoutesReporter) loadShipments(ctx context.Context, routeID, supplierID int, destinationAddressName string) error {
	cache, ok := r.cacheShipments[routeID]
	if !ok {
		cache = &cacheWithTime[[]*wb_models.Shipment]{}
		r.cacheShipments[routeID] = cache
	}

	now := time.Now()
	if now.Before(cache.time) {
		return nil
	}

	cache.value = nil

	var err error
	var shipments []*wb_models.Shipment
	err = retryAction(ctx, "GeneralRoutesReporter.loadShipments", 3, 1*time.Second, func() error {
		shipments, _, err = r.services.WBLogisticService.GetShipments(ctx, &models.WBLogisticGetShipmentsParamsRequest{
			DataStart:       now.AddDate(0, 0, -2),
			DataEnd:         now,
			SrcOfficeID:     r.config.Logistic().Office().ID(),
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
		return err
	}

	if shipments == nil {
		return errors.New("GeneralRoutesReporter.loadShipments()", "shipments is empty")
	}

	// Some routes may not have shipments, so we use an empty object as a stub
	cache.value = shipments
	cache.time = now.Add(r.ttlDurationLoadShipments)
	r.isNeedRenderReport = true
	r.isNeedProcessShipments = true
	return nil
}

func (r *GeneralRoutesReporter) sendGoogleSheets(data [][]interface{}) error {
	err := r.services.GoogleSheetsService.ClearValues(r.spreadsheetID, r.sheetName, "A:Z")
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.sendGoogleSheets()", "failed clear google sheets")
	}

	err = r.services.GoogleSheetsService.UpdateValues(r.spreadsheetID, r.sheetName, r.sheetPosition, data, false)
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.sendGoogleSheets()", "failed update google sheets")
	}

	return nil
}

func (r *GeneralRoutesReporter) sendTelegramBot(data string) error {
	err := r.services.TelegramBotService.SendMessage(r.config.Telegram().ChatID(), data, "HTML")
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.sendTelegramBot()", "failed send telegram bot")
	}
	return nil
}
