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

type shipmentCloseOpenedData struct {
	RouteID                          int
	Parking                          int
	RemainsLastMileReportsRoutesInfo []*wb_models.RemainsLastMileReportsRouteInfo
	Shipment                         *wb_models.Shipment
	PrevShipmentInfo                 *wb_models.ShipmentInfo
	SpName                           string
}

type ShipmentCloseReporter struct {
	config         *config.Config
	storage        storage.Storage
	services       *services.Container
	report         *reports.ShipmentCloseReport
	prompter       prompters.ShipmentCloseReporterPrompter
	messageQueueTG *Queue[string]

	rendererGS   report_renderers.ReportRenderer[[][]interface{}]
	rendererTG   report_renderers.ReportRenderer[[]string]
	renderModeTG report_renderers.TelegramBotRenderMode
	isRenderGS   bool
	isRenderTG   bool

	// params
	findShipmentsInterval         time.Duration
	prevTimeFindShipmentsInterval time.Time
	officeID                      int
	parking                       map[int]int      // route id -> parking
	suppliers                     map[int]struct{} // supplier id -> struct{}
	skipRoutes                    map[int]struct{} // route id -> struct{}
	paramShipmentsPageIndex       int
	paramShipmentsLimit           int
	paramShipmentsDirection       int
	paramShipmentsSorter          string

	// google sheets
	spreadsheetID string
	sheetName     string
	sheetPosition string

	// telegram
	tgChatID int64

	// cache
	cacheRemainsLastMileReport           *wb_models.RemainsLastMileReport
	cacheRemainsLastMileReportsRouteInfo map[int][]*wb_models.RemainsLastMileReportsRouteInfo // key route id
	cacheOpenedShipments                 map[int]*shipmentCloseOpenedData                     // key shipment id
	cachePrevShipments                   map[int]*wb_models.Shipment                          // key shipment id
	cachePrevShipmentsInfo               map[int]*wb_models.ShipmentInfo                      // key shipment id
}

func NewShipmentCloseReporter(config *config.Config, storage storage.Storage, service *services.Container, prompter prompters.ShipmentCloseReporterPrompter) *ShipmentCloseReporter {
	suppliers := make(map[int]struct{})
	for _, supplier := range config.Logistic().Office().Suppliers() {
		suppliers[supplier] = struct{}{}
	}
	skipRoutes := make(map[int]struct{})
	for _, route := range config.Reports().GeneralRoutes().SkipRoutes() {
		skipRoutes[route] = struct{}{}
	}
	return &ShipmentCloseReporter{
		config:                               config,
		storage:                              storage,
		services:                             service,
		report:                               reports.NewShipmentCloseReport(),
		rendererGS:                           &report_renderers.GoogleSheetsRenderer{},
		rendererTG:                           &report_renderers.TelegramBotRenderer{Mode: report_renderers.TelegramBotRenderHTML, IsTitle: false},
		renderModeTG:                         report_renderers.TelegramBotRenderHTML,
		isRenderGS:                           config.Reports().ShipmentClose().IsRenderGoogleSheets(),
		isRenderTG:                           config.Reports().ShipmentClose().IsRenderTelegramBot(),
		prompter:                             prompter,
		messageQueueTG:                       New[string](300),
		suppliers:                            suppliers,
		skipRoutes:                           skipRoutes,
		paramShipmentsPageIndex:              0,
		paramShipmentsLimit:                  3,
		paramShipmentsDirection:              -1,
		paramShipmentsSorter:                 "updated_at",
		findShipmentsInterval:                config.Reports().ShipmentClose().FindShipmentsInterval(),
		officeID:                             config.Logistic().Office().ID(),
		spreadsheetID:                        config.GoogleSheets().ReportSheets().GeneralRoutes().SpreadsheetID(),
		sheetName:                            config.GoogleSheets().ReportSheets().GeneralRoutes().SheetName(),
		sheetPosition:                        "A1",
		cacheRemainsLastMileReportsRouteInfo: map[int][]*wb_models.RemainsLastMileReportsRouteInfo{},
		cacheOpenedShipments:                 map[int]*shipmentCloseOpenedData{},
		cachePrevShipments:                   map[int]*wb_models.Shipment{},
		tgChatID:                             config.Telegram().ChatID(),
		parking:                              config.Logistic().Office().Parking(),
	}
}

func (r *ShipmentCloseReporter) Run(ctx context.Context) error {
	timeStart := time.Now()

	//err := r.sendReport(ctx, &reports.ShipmentCloseReportData{
	//	RouteID:               111111,
	//	ShipmentID:            222222,
	//	WaySheetID:            33333333,
	//	Parking:               444444,
	//	DriverName:            "info.Dr!!!iverName",
	//	VehicleNumberPlate:    "info.VehicleNumberPlate",
	//	TotalRemainsBarcodes:  5555,
	//	TotalRemainsTares:     6666,
	//	TotalTransferBarcodes: 777,
	//	TotalTransferTares:    8888,
	//	Date:                  time.Now(),
	//	TimeCreate:            time.Now(),
	//	TimeClose:             time.Now(),
	//	SpName:                "data.SpName", // It is taken from these data because they were obtained when the remains were still there and this value is guaranteed to be known.
	//	RemainsTaresInfo: []*reports.ShipmentCloseRemainsTareInfo{
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//		{11, 22, "OfficeName", time.Now(), 33},
	//	},
	//})

	//err := r.sendReport(ctx, &reports.ShipmentCloseReportData{
	//	RouteID:               111111,
	//	ShipmentID:            222222,
	//	WaySheetID:            33333333,
	//	Parking:               444444,
	//	DriverName:            "info.Dr!!!iverName",
	//	VehicleNumberPlate:    "info.VehicleNumberPlate",
	//	TotalRemainsBarcodes:  5555,
	//	TotalRemainsTares:     6666,
	//	TotalTransferBarcodes: 777,
	//	TotalTransferTares:    8888,
	//	Date:                  time.Now(),
	//	TimeCreate:            time.Now(),
	//	TimeClose:             time.Now(),
	//	SpName:                "data.SpName", // It is taken from these data because they were obtained when the remains were still there and this value is guaranteed to be known.
	//	RemainsTaresInfo: []*reports.ShipmentCloseRemainsTareInfo{
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//		{136003223891, 360089, "Гатчина проспект 25 Октября 28А", time.Now(), 33},
	//	},
	//})
	//if err != nil {
	//	panic(err)
	//}

	//os.Exit(111)
	if timeStart.After(r.prevTimeFindShipmentsInterval.Add(r.findShipmentsInterval)) {
		if err := r.findOpenedShipments(ctx); err != nil {
			r.prompter.PromptError("failed finding opened shipments")
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.Run()", "failed finding opened shipments: %v", err)
		}
		r.prevTimeFindShipmentsInterval = timeStart
	}

	if err := r.processOpenedShipments(ctx); err != nil {
		r.prompter.PromptError("failed processing opened shipments")
		return errors.Wrap(err, "ShipmentCloseReporter.Run()", "failed processing opened shipments")
	}

	if r.messageQueueTG.Len() > 0 {
		err := r.sendTelegramBot(ctx)
		if err != nil {
			return errors.Wrap(err, "ShipmentCloseReporter.Run()", "failed sending telegram bot")
		}
	}
	
	r.prompter.PromptRender(time.Since(timeStart))
	return nil
}

func (r *ShipmentCloseReporter) findOpenedShipments(ctx context.Context) error {
	remainsReport, err := r.loadRemainsLastMileReport(ctx)
	if err != nil {
		return errors.Wrap(err, "ShipmentCloseReporter.findOpenedShipments()", "failed load remains last miles report")
	}

	for _, route := range remainsReport.Routes {
		if !r.isValidRoute(route) {
			continue
		}

		routeID := route.CarID
		parking := r.parking[routeID]

		routeInfo, err := r.loadRemainsLastMileReportRouteInfo(ctx, routeID)
		if err != nil {
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.findOpenedShipments", "failed load route info for route %d: %v", routeID, err)
			continue
		}
		if len(route.Suppliers) == 0 || routeInfo[0].DestinationOfficeName == "" {
			continue
		}

		shipments, err := r.loadShipments(ctx, routeID, route.Suppliers[0].ID, routeInfo[0].DestinationOfficeName)
		if err != nil {
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.findOpenedShipments()", "failed load shipments: %v", err)
			continue
		}
		if len(shipments) == 0 {
			continue
		}

		// Checking all shipments for opening. Here information about open shipments is collected, which will then be tracked separately.
		for _, shipment := range shipments {
			// if the shipment is closed and there is nothing in the cache, then it has already been processed
			prev := r.cachePrevShipments[shipment.ShipmentID]

			// already closed - skipping
			if !shipment.CloseDt.IsZero() && (prev == nil || !prev.CloseDt.IsZero()) {
				continue
			}

			opened := r.cacheOpenedShipments[shipment.ShipmentID]
			if opened == nil {
				opened = &shipmentCloseOpenedData{}
				r.cacheOpenedShipments[shipment.ShipmentID] = opened
			}

			opened.RouteID = routeID
			opened.Parking = parking
			opened.Shipment = shipment
			opened.RemainsLastMileReportsRoutesInfo = routeInfo

			r.cachePrevShipments[shipment.ShipmentID] = shipment

			logger.Logf(logger.INFO, "ShipmentCloseReporter.Run()", "open shipment: route %d, shipment %d. Cache size: prev shipments: %d; opened shipments: %d", routeID, shipment.ShipmentID, len(r.cachePrevShipments), len(r.cacheOpenedShipments))
		}
	}
	return nil
}

func (r *ShipmentCloseReporter) processOpenedShipments(ctx context.Context) error {
	for shipmentID, data := range r.cacheOpenedShipments {
		time.Sleep(100 * time.Millisecond)

		info, err := r.loadShipmentInfo(ctx, shipmentID)
		if err != nil {
			r.prompter.PromptError(fmt.Sprintf("failed loading shipment info for shipment %d", shipmentID))
			r.removeShipmentFromCache(shipmentID)
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.processOpenedShipments", "failed load shipment info, route %d, shipment %d: %v", data.RouteID, shipmentID, err)
			continue
		}

		if data.PrevShipmentInfo == nil {
			data.PrevShipmentInfo = info
		}
		prevInfo := data.PrevShipmentInfo
		data.PrevShipmentInfo = info

		if data.SpName == "" && len(info.DestinationOfficesInfo) > 0 {
			dstOfficesIDs := make([]int, len(info.DestinationOfficesInfo))
			for i, v := range info.DestinationOfficesInfo {
				dstOfficesIDs[i] = v.DstOfficeID
			}
			data.SpName, err = r.loadSpName(ctx, dstOfficesIDs)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("failed loading sp name for shipment %d", shipmentID))
				logger.Logf(logger.ERROR, "ShipmentCloseReporter.processOpenedShipments()", "failed load sp name for shipment %d: %v", shipmentID, err)
			}
		}

		// if according to previous data shipment was open, and according to current data it is closed, then it is closed
		// or if it is not possible to send report, then current and previous shipment will be closed, so we try to send report again
		if prevInfo.CloseDt.IsZero() && !info.CloseDt.IsZero() ||
			(!prevInfo.CloseDt.IsZero() && !info.CloseDt.IsZero()) {

			dstOfficesIDs := make([]int, len(info.DestinationOfficesInfo))
			for i, v := range info.DestinationOfficesInfo {
				dstOfficesIDs[i] = v.DstOfficeID
			}

			remainsTares, err := r.loadRemainsTares(ctx, dstOfficesIDs)
			if err != nil {
				remainsTares = make([]*wb_models.TareForOffice, 0)
				r.prompter.PromptError(fmt.Sprintf("failed loading remains tares for shipment %d", shipmentID))
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

			transferBoxes, err := r.loadShipmentTransfersBoxes(ctx, shipmentID)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("failed loading shipment transfers boxes for shipment %d", shipmentID))
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

			err = r.sendReport(ctx, &reports.ShipmentCloseReportData{
				RouteID:               data.RouteID,
				ShipmentID:            shipmentID,
				WaySheetID:            info.WaySheetID,
				Parking:               data.Parking,
				DriverName:            info.DriverName,
				VehicleNumberPlate:    info.VehicleNumberPlate,
				TotalRemainsBarcodes:  totalRemainsBarcodes,
				TotalRemainsTares:     len(remainsTares),
				TotalTransferBarcodes: totalTransferBarcodes,
				TotalTransferTares:    len(transferBoxes),
				Date:                  info.CreateDt,
				TimeCreate:            info.CreateDt,
				TimeClose:             info.CloseDt,
				SpName:                data.SpName, // It is taken from these data because they were obtained when the remains were still there and this value is guaranteed to be known.
				RemainsTaresInfo:      remainsTaresInfo,
			})
			if err != nil {
				logger.Logf(logger.ERROR, "ShipmentCloseReporter.processOpenedShipments()", "failed send report: %v", err)
				continue
			}
			r.removeShipmentFromCache(shipmentID)
		}
	}
	return nil
}

func (r *ShipmentCloseReporter) isValidRoute(route *wb_models.Route) bool {
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

func (r *ShipmentCloseReporter) removeShipmentFromCache(shipmentID int) {
	delete(r.cacheOpenedShipments, shipmentID)
	delete(r.cachePrevShipments, shipmentID)
}

func (r *ShipmentCloseReporter) loadSpName(ctx context.Context, dstOfficeIDs []int) (string, error) {
	var err error
	var tares []*wb_models.TareForOffice
	err = retryAction(ctx, "ShipmentCloseReporter.loadShipmentTransfers", 3, 1*time.Second, func() error {
		tares, err = r.loadRemainsTares(ctx, dstOfficeIDs)
		return err
	})
	if err != nil {
		return "", errors.Wrap(err, "ShipmentCloseReporter.loadSpName()", "")
	}

	if len(tares) == 0 {
		return "", nil
	}
	// All tares have the same SpName value
	return tares[0].SpName, nil
}

func (r *ShipmentCloseReporter) loadRemainsLastMileReport(ctx context.Context) (*wb_models.RemainsLastMileReport, error) {
	var err error
	var remainsLastMileReport *wb_models.RemainsLastMileReport
	err = retryAction(ctx, "ShipmentCloseReporter.loadRemainsLastMileReport", 3, 1*time.Second, func() error {
		remainsLastMileReport, err = r.services.WBLogisticService.GetRemainsLastMileReportByOfficeID(ctx, r.officeID)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Just in case, for reliability
	if remainsLastMileReport == nil || remainsLastMileReport.Routes == nil {
		return nil, errors.New("ShipmentCloseReporter.loadRemainsLastMileReport()", "remains last miles report is empty")
	}

	r.cacheRemainsLastMileReport = remainsLastMileReport
	return remainsLastMileReport, nil
}

func (r *ShipmentCloseReporter) loadRemainsLastMileReportRouteInfo(ctx context.Context, routeID int) ([]*wb_models.RemainsLastMileReportsRouteInfo, error) {
	r.cacheRemainsLastMileReportsRouteInfo[routeID] = nil

	var err error
	var remainsLastMileReportsRouteInfo []*wb_models.RemainsLastMileReportsRouteInfo
	err = retryAction(ctx, "ShipmentCloseReporter.loadRemainsLastMileReportRouteInfo()", 3, 1*time.Second, func() error {
		remainsLastMileReportsRouteInfo, err = r.services.WBLogisticService.GetRemainsLastMileReportsRouteInfo(ctx, routeID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "ShipmentCloseReporter.loadRemainsLastMileReportRouteInfo()", "failed load remains last miles reports route %d info", routeID)
	}

	// Just in case, for reliability
	if remainsLastMileReportsRouteInfo == nil || len(remainsLastMileReportsRouteInfo) == 0 {
		return nil, errors.Newf("ShipmentCloseReporter.loadRemainsLastMileReportRouteInfo()", "remains last miles reports route %d info is empty", routeID)
	}

	r.cacheRemainsLastMileReportsRouteInfo[routeID] = remainsLastMileReportsRouteInfo
	return remainsLastMileReportsRouteInfo, nil
}

func (r *ShipmentCloseReporter) loadShipments(ctx context.Context, routeID, supplierID int, destinationAddressName string) ([]*wb_models.Shipment, error) {
	var err error
	var shipments []*wb_models.Shipment
	now := time.Now()
	err = retryAction(ctx, "ShipmentCloseReporter.loadShipments()", 3, 1*time.Second, func() error {
		shipments, _, err = r.services.WBLogisticService.GetShipments(ctx, &models.WBLogisticGetShipmentsParamsRequest{
			DataStart:       now.AddDate(0, 0, -2),
			DataEnd:         now,
			SrcOfficeID:     r.officeID,
			PageIndex:       r.paramShipmentsPageIndex,
			Limit:           r.paramShipmentsLimit,
			SupplierID:      supplierID,
			Direction:       r.paramShipmentsDirection,
			Sorter:          r.paramShipmentsSorter,
			FilterDstOffice: destinationAddressName,
		})
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "ShipmentCloseReporter.loadShipments()", "failed load shipments on route %d for supplier %d", routeID, supplierID)
	}

	// Just in case, for reliability
	if shipments == nil {
		return nil, errors.Newf("ShipmentCloseReporter.loadShipments()", "shipments is empty on route %d for supplier %d", routeID, supplierID)
	}

	return shipments, nil
}

func (r *ShipmentCloseReporter) loadShipmentInfo(ctx context.Context, shipmentID int) (*wb_models.ShipmentInfo, error) {
	var err error
	var shipmentInfo *wb_models.ShipmentInfo
	err = retryAction(ctx, "ShipmentCloseReporter.loadShipmentInfo()", 3, 1*time.Second, func() error {
		shipmentInfo, err = r.services.WBLogisticService.GetShipmentInfo(ctx, shipmentID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "ShipmentCloseReporter.loadShipmentInfo()", "failed load shipment %d info", shipmentID)
	}

	// Just in case, for reliability
	if shipmentInfo == nil {
		return nil, errors.Newf("ShipmentCloseReporter.loadShipmentInfo()", "shipment %d info returned empty value without error", shipmentID)
	}

	// Just in case, for reliability
	if shipmentInfo.DestinationOfficesInfo == nil {
		shipmentInfo.DestinationOfficesInfo = make([]*wb_models.ShipmentInfoDestinationOfficeInfo, 0)
	}

	return shipmentInfo, nil
}

func (r *ShipmentCloseReporter) loadRemainsTares(ctx context.Context, dstOfficeIDs []int) ([]*wb_models.TareForOffice, error) {
	var err error
	var tares []*wb_models.TareForOffice
	err = retryAction(ctx, "ShipmentCloseReporter.loadRemainsTares", 3, 1*time.Second, func() error {
		// isDrive = false is default value
		tares, err = r.services.WBLogisticService.GetTaresForOffices(ctx, r.officeID, dstOfficeIDs, false)
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "ShipmentCloseReporter.loadRemainsTares()", "failed load tares for offices")
	}

	// Just in case, for reliability
	if tares == nil || len(tares) == 0 {
		return nil, errors.New("ShipmentCloseReporter.loadRemainsTares()", "remains last miles reports route info is empty")
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
		return errors.Wrapf(err, "ShipmentCloseReporter.sendReport()", "failed render report, route id: %d shipment id: %d", reportData.RouteID, reportData.ShipmentID)
	}

	if r.isRenderGS {
		data, err := r.rendererGS.Render(report)
		if err != nil {
			r.prompter.PromptError("failed to render report Google Sheet")
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.sendReport()", "failed render report for google sheets: %v, route id: %d shipment id: %d", err, reportData.RouteID, reportData.ShipmentID)
		} else {
			err = r.sendGoogleSheets(ctx, data)
			if err != nil {
				r.prompter.PromptError("failed to send report Google Sheet")
				logger.Logf(logger.ERROR, "ShipmentCloseReporter.sendReport()", "failed send report to google sheets: %v, route id: %d shipment id: %d", err, reportData.RouteID, reportData.ShipmentID)
			}
		}
	}

	if r.isRenderTG {
		messages, err := r.rendererTG.Render(report)
		if err != nil {
			r.prompter.PromptError("failed to render report Telegram Bot")
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.sendReport()", "failed render report for telegram bot: %v, route id: %d shipment id: %d", err, reportData.RouteID, reportData.ShipmentID)
		}

		if len(messages) == 0 {
			return nil
		}

		for _, message := range messages {
			if message != "" {
				r.messageQueueTG.Push(message)
			}
		}

		err = r.sendTelegramBot(ctx)
		if err != nil {
			r.prompter.PromptError("failed to send report Telegram Bot")
			logger.Logf(logger.ERROR, "ShipmentCloseReporter.sendReport()", "failed send report to telegram bot: %v, route id: %d shipment id: %d", err, reportData.RouteID, reportData.ShipmentID)
		}
	}

	return nil
}

func (r *ShipmentCloseReporter) sendGoogleSheets(ctx context.Context, data [][]interface{}) error {
	err := r.services.GoogleSheetsService.ClearValues(r.spreadsheetID, r.sheetName, "A:Z")
	if err != nil {
		return errors.Wrapf(err, "ShipmentCloseReporter.sendGoogleSheets()", "failed clear sheet %s, page %s", r.spreadsheetID, r.sheetName)
	}

	err = retryAction(ctx, "ShipmentCloseReporter.sendGoogleSheets()", 3, 1*time.Second, func() error {
		err = r.services.GoogleSheetsService.UpdateValues(r.spreadsheetID, r.sheetName, r.sheetPosition, data, false)
		return err
	})
	if err != nil {
		return errors.Wrapf(err, "GeneralRoutesReporter.sendGoogleSheets()", "failed update sheet %s, page %s to position %s", r.spreadsheetID, r.sheetName, r.sheetPosition)
	}

	return nil
}

func (r *ShipmentCloseReporter) sendTelegramBot(ctx context.Context) error {
	//fmt.Println(data, len(data))
	//os.Exit(11)

	if r.messageQueueTG.Len() <= 0 {
		return nil
	}

	for r.messageQueueTG.Len() > 0 {
		var err error
		err = retryAction(ctx, "ShipmentCloseReporter.sendTelegramBot()", 3, 1*time.Second, func() error {
			data, ok := r.messageQueueTG.Peek()
			if !ok {
				return errors.New("ShipmentCloseReporter.sendTelegramBot()", "failed to get message from telegram message queue")
			}
			if r.renderModeTG == report_renderers.TelegramBotRenderHTML {
				err = r.services.TelegramBotService.SendMessage(r.config.Telegram().ChatID(), data, "HTML")
			} else if r.renderModeTG == report_renderers.TelegramBotRenderMarkdown {
				err = r.services.TelegramBotService.SendMessage(r.config.Telegram().ChatID(), data, "MarkdownV2")
			} else {
				err = r.services.TelegramBotService.SendMessage(r.config.Telegram().ChatID(), data, "")
			}
			return err
		})
		if err != nil {
			return errors.Wrapf(err, "ShipmentCloseReporter.sendTelegramBot()", "failed send data to chat %d", r.config.Telegram().ChatID())
		}
		r.messageQueueTG.Pop()
	}

	return nil
}
