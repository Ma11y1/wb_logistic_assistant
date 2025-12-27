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

type generalRoutesRouteData struct {
	routeID            int
	changeBarcodes     int
	prevBarcodes       int
	shipmentID         int
	isShipmentOpen     bool
	lastOpenWaySheet   *wb_models.WaySheet
	bestClosedWaySheet *wb_models.WaySheet
	prevClosedWaySheet *wb_models.WaySheet
}

type GeneralRoutesReporter struct {
	config      *config.Config
	storage     storage.Storage
	services    *services.Container
	prompter    prompters.GeneralRoutesReporterPrompter
	reportSheet *reports.GeneralRoutesSheetReport

	rendererGS report_renderers.ReportRenderer[[][]interface{}]
	isRenderGS bool

	officeID   int
	suppliers  map[int]struct{} // supplier id -> struct{}
	skipRoutes map[int]struct{} // route id -> struct{}

	intervalResetChangeBarcodes time.Duration
	intervalUpdateRating        time.Duration
	intervalUpdateShipments     time.Duration
	intervalUpdateWaySheets     time.Duration
	prevResetChangeBarcodes     time.Time
	prevUpdateRating            time.Time
	prevUpdateShipments         time.Time
	prevUpdateWaySheets         time.Time

	spreadsheetID string
	sheetName     string
	sheetPosition string

	chReportMetaData        *reports.GeneralRoutesReportMetaData
	chReportDataList        []*reports.GeneralRoutesReportData
	chReportData            map[int]*reports.GeneralRoutesReportData // report id -> ReportData
	chRouteData             map[int]*generalRoutesRouteData          // report id -> RoutesData
	chRemainsLastMileReport *wb_models.RemainsLastMileReport
}

func NewGeneralRoutesReporter(config *config.Config, storage storage.Storage, services *services.Container, prompter prompters.GeneralRoutesReporterPrompter) *GeneralRoutesReporter {
	return &GeneralRoutesReporter{
		config:   config,
		storage:  storage,
		services: services,
		prompter: prompter,

		reportSheet: reports.NewGeneralRoutesSheetReport(
			config.Reports().GeneralRoutes().IsSort(),
			config.Reports().GeneralRoutes().SortColumn(),
			config.Reports().GeneralRoutes().IsAscending(),
		),

		rendererGS: &report_renderers.GoogleSheetsRenderer{},
		isRenderGS: config.Reports().GeneralRoutes().IsRenderGoogleSheets(),

		officeID:                    config.Logistic().Office().ID(),
		suppliers:                   config.Logistic().Office().SuppliersMap(),
		skipRoutes:                  config.Logistic().Office().SkipRoutesMap(),
		intervalResetChangeBarcodes: config.Reports().GeneralRoutes().IntervalResetChangeBarcodes(),
		intervalUpdateRating:        config.Reports().GeneralRoutes().IntervalUpdateRating(),
		intervalUpdateShipments:     config.Reports().GeneralRoutes().IntervalUpdateShipments(),
		intervalUpdateWaySheets:     config.Reports().GeneralRoutes().IntervalUpdateWaySheets(),

		spreadsheetID: config.GoogleSheets().ReportSheets().GeneralRoutes().SpreadsheetID(),
		sheetName:     config.GoogleSheets().ReportSheets().GeneralRoutes().SheetName(),
		sheetPosition: "A1",

		chReportMetaData: &reports.GeneralRoutesReportMetaData{},
		chReportDataList: make([]*reports.GeneralRoutesReportData, 0, 10),
		chReportData:     map[int]*reports.GeneralRoutesReportData{},
		chRouteData:      map[int]*generalRoutesRouteData{},
	}
}

func (r *GeneralRoutesReporter) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "GeneralRoutesReporter.Run()", "task was preliminarily completed")
	}

	r.prompter.PromptStart()

	timeStart := time.Now()

	remainsReport, err := r.loadRemainsLastMileReport(ctx)
	if err != nil {
		r.prompter.PromptError("Failed load remains last miles report")
		return errors.Wrap(err, "GeneralRoutesReporter.Run()", "failed load remains last miles report")
	}

	r.chReportDataList = r.chReportDataList[:0]
	isResetChangeBarcodes := timeStart.After(r.prevResetChangeBarcodes.Add(r.intervalResetChangeBarcodes))
	isUpdateShipments := timeStart.After(r.prevUpdateShipments.Add(r.intervalUpdateShipments))

	for _, route := range remainsReport.Routes {
		if !r.isValidRoute(route) {
			continue
		}
		time.Sleep(100 * time.Millisecond)

		routeID := route.CarID
		routeData := r.chRouteData[routeID]
		reportData := r.chReportData[routeID]

		if routeData == nil {
			routeData = &generalRoutesRouteData{routeID: routeID}
			r.chRouteData[routeID] = routeData
		}
		if reportData == nil {
			reportData = &reports.GeneralRoutesReportData{RouteID: routeID}
			r.chReportData[routeID] = reportData
		}
		r.chReportDataList = append(r.chReportDataList, reportData)

		reportData.Tares = route.CountTares
		reportData.VolumeLiters = float32(route.VolumeMlByContent) / 1000 // Ml to litters
		if route.NormativeInLiters > 0 {
			reportData.VolumeNormativeLitersPercent = reportData.VolumeLiters / route.NormativeInLiters * 100
		} else {
			reportData.VolumeNormativeLitersPercent = 0
		}
		reportData.VolumeNormativeLiters = route.NormativeInLiters

		// Barcodes
		if isResetChangeBarcodes {
			r.chReportMetaData.TimeUpdateChangeBarcodes = time.Now()
			routeData.changeBarcodes = 0
		} else {
			// it's worth using 'else' so that the first pass is not the current value of the barcodes, but is immediately reset
			routeData.changeBarcodes += route.CountShk - routeData.prevBarcodes
		}
		reportData.Barcodes = route.CountShk
		reportData.ChangeBarcodes = routeData.changeBarcodes
		routeData.prevBarcodes = route.CountShk

		// First Parking and Sp name. Next updates to the parking and sp name will be when the shipment closes
		if reportData.Parking == 0 {
			routeInfo, err := r.loadRemainsLastMileReportInfo(ctx, routeID)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed load route info for route %d", routeID))
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed load route info for route %d: %v", routeID, err)
				continue
			}
			parking, _, err := r.loadParkingAndRemainsBarcodesByRouteInfo(ctx, routeInfo)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed load parking for route %d", routeID))
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed load parking and sp name for route %d: %v", routeID, err)
			} else {
				reportData.Parking = parking
			}
		}

		// Shipments
		if isUpdateShipments {
			routeInfo, err := r.loadRemainsLastMileReportInfo(ctx, routeID)
			if err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed load route info for route %d", routeID))
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed load route info for route %d: %v", routeID, err)
				continue
			}

			if err = r.processShipments(ctx, route, routeInfo); err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed process shipments for route %d", routeID))
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed process shipments for route %d: %v", routeID, err)
			}
		}
	}

	if isResetChangeBarcodes {
		r.prevResetChangeBarcodes = timeStart
	}
	if isUpdateShipments {
		logger.Log(logger.INFO, "GeneralRoutesReporter.Run()", "shipments updated")
		r.prompter.PromptUpdateShipments()
		r.prevUpdateShipments = timeStart
	}

	// Way sheets
	if timeStart.After(r.prevUpdateWaySheets.Add(r.intervalUpdateWaySheets)) {
		if err = r.processWaySheets(ctx, remainsReport); err != nil {
			r.prompter.PromptError("Failed process way sheets")
			logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed process way sheets: %v", err)
		} else {
			r.prompter.PromptUpdateWaySheets()
			logger.Log(logger.INFO, "GeneralRoutesReporter.Run()", "way sheets updated")
			r.prevUpdateWaySheets = timeStart
		}
	}

	// Rating
	if timeStart.After(r.prevUpdateRating.Add(r.intervalUpdateRating)) {
		if err = r.processRating(ctx); err != nil {
			r.prompter.PromptError("Failed process rating")
			logger.Logf(logger.ERROR, "GeneralRoutesReporter.Run()", "failed processing rating: %v", err)
		} else {
			r.prompter.PromptUpdateRating()
			logger.Log(logger.INFO, "GeneralRoutesReporter.Run()", "rating updated")
			r.prevUpdateRating = timeStart
		}
	}

	logger.Logf(logger.INFO, "GeneralRoutesReporter.Run()", "update routes %d", len(remainsReport.Routes))
	r.prompter.PromptUpdateRoutes(len(remainsReport.Routes))

	r.chReportMetaData.Update = time.Now()

	if err = r.sendReport(ctx, r.chReportMetaData, r.chReportDataList); err != nil {
		r.prompter.PromptError("Failed send report")
		return errors.Wrap(err, "GeneralRoutesReporter.Run()", "failed send report")
	}

	r.prompter.PromptFinish(time.Since(timeStart))
	return nil
}

func (r *GeneralRoutesReporter) processRating(ctx context.Context) error {
	logger.Log(logger.INFO, "GeneralRoutesReporter.processShipments()", "start process rating")

	jobsScheduling, err := r.loadJobsScheduling(ctx)
	if err != nil {
		return err
	}

	for _, route := range jobsScheduling.Route {
		if route == nil || route.Rating == nil || route.SrcOfficeId != r.officeID {
			continue
		}

		reportData := r.chReportData[route.RouteID]
		if reportData == nil {
			logger.Logf(logger.WARN, "GeneralRoutesReporter.processRating()", "there is no report data for the route %d", route.RouteID)
			continue
		}

		reportData.Rating = float32(route.Rating.OverallRating)
	}
	r.chReportMetaData.TimeUpdateRating = time.Now()
	return nil
}

func (r *GeneralRoutesReporter) processShipments(ctx context.Context, route *wb_models.Route, routesInfo []*wb_models.RemainsLastMileReportsRouteInfo) error {
	if len(routesInfo) == 0 {
		return errors.Newf("GeneralRoutesReporter.processShipments()", "failed process shipments: destination office addresses could not be found, route %d", route.CarID)
	}

	officeName := ""
	for _, info := range routesInfo {
		if info != nil && info.DestinationOfficeName != "" {
			officeName = info.DestinationOfficeName
			break
		}
	}

	if officeName == "" {
		return errors.Newf("GeneralRoutesReporter.processShipments()", "failed process shipments: destination office addresses could not be found, route %d", route.CarID)
	}

	shipments, err := r.loadShipments(ctx, route.CarID, officeName)
	if err != nil || len(shipments) == 0 {
		return errors.Wrap(err, "GeneralRoutesReporter.processShipments()", "failed load shipments")
	}

	r.chReportMetaData.TimeUpdateShipments = time.Now()

	routeData := r.chRouteData[route.CarID]
	reportData := r.chReportData[route.CarID]
	if routeData == nil || reportData == nil {
		return errors.Newf("GeneralRoutesReporter.processShipments()", "failed to retrieve data cache, it's empty")
	}

	shipment := r.findLastShipment(shipments)
	if shipment == nil {
		return errors.Newf("GeneralRoutesReporter.processShipments()", "failed find last shipment, it's empty")
	}

	currentShipmentID := shipment.ShipmentID
	prevShipmentID := routeData.shipmentID

	// First process, if there was no data at all previously
	if prevShipmentID == 0 {
		routeData.shipmentID = currentShipmentID
		reportData.ShipmentID = currentShipmentID
		reportData.ShipmentCreateDate = shipment.CreateDt
		reportData.ShipmentCloseDate = shipment.CloseDt
		routeData.isShipmentOpen = shipment.CloseDt.IsZero()
		return nil
	}

	// Subsequent process

	if currentShipmentID != prevShipmentID {
		parking, _, err := r.loadParkingAndRemainsBarcodesByRouteInfo(ctx, routesInfo)
		if err != nil {
			return errors.Wrapf(err, "GeneralRoutesReporter.processShipments()", "failed load shipment remains barcodes for shipment %d", shipment.ShipmentID)
		}
		routeData.shipmentID = currentShipmentID
		routeData.isShipmentOpen = shipment.CloseDt.IsZero()

		reportData.ShipmentID = currentShipmentID
		reportData.ShipmentCreateDate = shipment.CreateDt
		reportData.ShipmentCloseDate = shipment.CloseDt
		reportData.Parking = parking
		logger.Logf(logger.INFO, "GeneralRoutesReporter.processShipments()", "open new shipment %d", currentShipmentID)
		return nil
	}

	if !shipment.CloseDt.IsZero() && routeData.isShipmentOpen {
		parking, barcodes, err := r.loadParkingAndRemainsBarcodesByRouteInfo(ctx, routesInfo)
		if err != nil {
			return errors.Wrapf(err, "GeneralRoutesReporter.processShipments()", "failed load shipment remains barcodes for shipment %d", shipment.ShipmentID)
		}
		routeData.isShipmentOpen = shipment.CloseDt.IsZero()

		reportData.ShipmentCreateDate = shipment.CreateDt
		reportData.ShipmentCloseDate = shipment.CloseDt
		reportData.Parking = parking
		reportData.RemainsBarcodes = barcodes
		r.chReportMetaData.TimeRemainsBarcodes = time.Now()
		r.prompter.PromptCloseShipment(currentShipmentID, barcodes)
		logger.Logf(logger.INFO, "GeneralRoutesReporter.processShipments()", "close shipment %d", currentShipmentID)
	}

	return nil
}

func (r *GeneralRoutesReporter) processWaySheets(ctx context.Context, remainsReport *wb_models.RemainsLastMileReport) error {
	logger.Log(logger.INFO, "GeneralRoutesReporter.processWaySheets()", "start process way sheets")

	err := r.findLastAndPrevWaySheets(ctx)
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.processWaySheets()", "")
	}

	for _, route := range remainsReport.Routes {
		time.Sleep(100 * time.Millisecond)

		routeData := r.chRouteData[route.CarID]
		reportData := r.chReportData[route.CarID]
		if reportData == nil || routeData == nil {
			continue
		}

		var lastWaySheet, prevWaySheet *wb_models.WaySheet
		if routeData.lastOpenWaySheet != nil {
			lastWaySheet = routeData.lastOpenWaySheet
			prevWaySheet = routeData.bestClosedWaySheet
		} else {
			lastWaySheet = routeData.bestClosedWaySheet
			prevWaySheet = routeData.prevClosedWaySheet
		}

		if lastWaySheet != nil {
			if err = r.processWaySheet(ctx, reportData, atoiSafe(lastWaySheet.WaySheetID), true); err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed process last way sheet %s", lastWaySheet.WaySheetID))
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.processWaySheets()", "failed process last way sheet %s: %v", lastWaySheet.WaySheetID, err)
			}
		}

		if prevWaySheet != nil {
			if err = r.processWaySheet(ctx, reportData, atoiSafe(prevWaySheet.WaySheetID), false); err != nil {
				r.prompter.PromptError(fmt.Sprintf("Failed process previously way sheet %s", prevWaySheet.WaySheetID))
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.processWaySheets()", "failed process previously way sheet %s: %v", prevWaySheet.WaySheetID, err)
			}
		}

		if lastWaySheet != nil && prevWaySheet != nil {
			reportData.WaySheetsInterval = lastWaySheet.OpenDt.Sub(prevWaySheet.OpenDt)
		}

		routeData.lastOpenWaySheet = nil
		routeData.bestClosedWaySheet = nil
		routeData.prevClosedWaySheet = nil
	}

	r.chReportMetaData.TimeUpdateWaySheets = time.Now()
	return nil
}

func (r *GeneralRoutesReporter) findLastAndPrevWaySheets(ctx context.Context) error {
	waySheets, err := r.loadWaySheets(ctx)
	if err != nil {
		return errors.Wrap(err, "GeneralRoutesReporter.findLastAndPrevWaySheets()", "failed load way sheets")
	}

	for _, waySheet := range waySheets {
		waySheetRouteID := atoiSafe(waySheet.RouteCarID)

		routeData := r.chRouteData[waySheetRouteID]
		if routeData == nil {
			continue
		}

		if waySheet.CloseDt.IsZero() {
			if routeData.lastOpenWaySheet == nil || waySheet.OpenDt.After(routeData.lastOpenWaySheet.OpenDt) {
				routeData.lastOpenWaySheet = waySheet
			}
			continue
		}

		if routeData.bestClosedWaySheet == nil {
			routeData.bestClosedWaySheet = waySheet
			continue
		}

		if waySheet.CloseDt.After(routeData.bestClosedWaySheet.CloseDt) {
			routeData.prevClosedWaySheet = routeData.bestClosedWaySheet
			routeData.bestClosedWaySheet = waySheet
			continue
		}

		if routeData.prevClosedWaySheet == nil || waySheet.CloseDt.After(routeData.prevClosedWaySheet.CloseDt) {
			routeData.prevClosedWaySheet = waySheet
		}
	}

	return nil
}

func (r *GeneralRoutesReporter) processWaySheet(ctx context.Context, reportData *reports.GeneralRoutesReportData, waySheetID int, isLast bool) error {
	if isLast {
		reportData.WaySheetID = waySheetID
	} else {
		reportData.PrevWaySheetID = waySheetID
	}

	info, err := r.loadWaySheetInfo(ctx, waySheetID)
	if err != nil {
		return errors.Wrapf(err, "GeneralRoutesReporter.processWaySheet()", "failed load last way sheet %d info: %v", waySheetID, err)
	}

	offices := map[int]int{} // office id -> counter delivery tares
	counterCurrentReturnedTares := 0
	counterTotalReturnedTares := 0
	counterDeliveryOffices := 0
	var dateLastOperation time.Time
	// we go through all the tares and count how many tares were delivered to each address
	for _, tare := range info.Tares {
		if tare == nil {
			continue
		}

		if tare.IsReturn {
			counterTotalReturnedTares++
			// when the returned container is brought to the warehouse, a date appears
			if !tare.DtArrival.IsZero() {
				counterCurrentReturnedTares++
			}
		}

		// tares may have a warehouse destination address, they are returns and do not need to be taken into account
		if tare.DstOfficeID == "" || tare.DstOfficeID == info.SrcOffice.ID {
			// for some reason, the ID of some tares may be empty
			continue
		}

		officeID := atoiSafe(tare.DstOfficeID)
		// we always note that the address is present.
		v, ok := offices[officeID]
		if !ok {
			offices[officeID] = 0
		}

		// if the container was delivered, then time appears and not zero
		if !tare.DtArrival.IsZero() {
			// if the tare was closed in one of the offices, then we assume that this address was closed
			if v == 0 {
				counterDeliveryOffices++
			}
			if tare.DtArrival.After(dateLastOperation) {
				dateLastOperation = tare.DtArrival
			}
			offices[officeID]++
		}
	}

	if isLast {
		reportData.WaySheetTotalAddresses = len(offices)
		reportData.WaySheetCurrentAddresses = counterDeliveryOffices
		reportData.WaySheetDateLastOperation = dateLastOperation
		reportData.WaySheetCurrentReturnedTares = counterCurrentReturnedTares
		reportData.WaySheetTotalReturnedTares = counterTotalReturnedTares
	} else {
		reportData.PrevWaySheetTotalAddresses = len(offices)
		reportData.PrevWaySheetCurrentAddresses = counterDeliveryOffices
		reportData.PrevWaySheetDateLastOperation = dateLastOperation
		reportData.PrevWaySheetCurrentReturnedTares = counterCurrentReturnedTares
		reportData.PrevWaySheetTotalReturnedTares = counterTotalReturnedTares
	}

	return nil
}

func (r *GeneralRoutesReporter) isValidRoute(route *wb_models.Route) bool {
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

func (r *GeneralRoutesReporter) findLastShipment(shipments []*wb_models.Shipment) *wb_models.Shipment {
	if len(shipments) == 0 {
		return nil
	}
	lastShipment := shipments[0]
	for _, shipment := range shipments {
		if shipment.CloseDt.IsZero() {
			lastShipment = shipment
			break
		}
		if shipment.CloseDt.After(lastShipment.CloseDt) {
			lastShipment = shipment
		}
	}
	return lastShipment
}

func (r *GeneralRoutesReporter) loadRemainsLastMileReport(ctx context.Context) (remainsLastMileReport *wb_models.RemainsLastMileReport, err error) {
	err = retryAction(ctx, "GeneralRoutesReporter.loadRemainsLastMileReport", 3, 1*time.Second, func() error {
		remainsLastMileReport, err = r.services.WBLogisticService.GetRemainsLastMileReportByOfficeID(ctx, r.officeID)
		return err
	})
	if err != nil {
		return nil, err
	}
	if remainsLastMileReport != nil {
		r.chRemainsLastMileReport = remainsLastMileReport
	}
	return remainsLastMileReport, nil
}

func (r *GeneralRoutesReporter) loadRemainsLastMileReportInfo(ctx context.Context, routeID int) (remainsLastMileReportsRouteInfo []*wb_models.RemainsLastMileReportsRouteInfo, err error) {
	err = retryAction(ctx, "GeneralRoutesReporter.loadRemainsLastMileReportInfo", 3, 1*time.Second, func() error {
		remainsLastMileReportsRouteInfo, err = r.services.WBLogisticService.GetRemainsLastMileReportsRouteInfo(ctx, routeID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "GeneralRoutesReporter.loadRemainsLastMileReportInfo()", "failed load remains last miles reports route %d info", routeID)
	}
	if remainsLastMileReportsRouteInfo == nil {
		return nil, errors.Wrapf(err, "GeneralRoutesReporter.loadRemainsLastMileReportInfo()", "failed load remains last miles reports route %d info: reports is empty", routeID)
	}
	return remainsLastMileReportsRouteInfo, nil
}

func (r *GeneralRoutesReporter) loadJobsScheduling(ctx context.Context) (jobsScheduling *wb_models.JobsScheduling, err error) {
	err = retryAction(ctx, "GeneralRoutesReporter.loadJobsScheduling", 3, 1*time.Second, func() error {
		jobsScheduling, err = r.services.WBLogisticService.GetJobsScheduling(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return jobsScheduling, nil
}

func (r *GeneralRoutesReporter) loadShipments(ctx context.Context, routeID int, destinationAddressName string) (shipments []*wb_models.Shipment, err error) {
	now := time.Now()
	for supplierID := range r.suppliers {
		var res []*wb_models.Shipment
		err = retryAction(ctx, "GeneralRoutesReporter.loadShipments", 3, 1*time.Second, func() error {
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
			logger.Logf(logger.ERROR, "GeneralRoutesReporter.loadShipments()", "failed load shipments on route %d for supplier %d: %v", routeID, supplierID, err)
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
		return nil, errors.Wrapf(err, "GeneralRoutesReporter.loadShipments()", "failed load shipments on route %d", routeID)
	}
	return shipments, nil
}

func (r *GeneralRoutesReporter) loadParkingAndRemainsBarcodesByRouteInfo(ctx context.Context, info []*wb_models.RemainsLastMileReportsRouteInfo) (int, int, error) {
	if info == nil || len(info) == 0 {
		return 0, 0, nil
	}
	offices := make([]int, 0, len(info))
	for _, item := range info {
		if item == nil {
			continue
		}
		offices = append(offices, item.DestinationOfficeID)
	}

	remainsTares, err := r.loadRemainsTares(ctx, offices)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "GeneralRoutesReporter.loadParkingAndRemainsBarcodesByRouteInfo()", "failed load remains tares")
	}

	if len(remainsTares) == 0 {
		return 0, 0, nil
	}

	barcodes := 0
	for _, tare := range remainsTares {
		if tare == nil {
			continue
		}
		barcodes += tare.CountBarcodes
	}

	_, parking := SpNameToGateParking(remainsTares[0].SpName)

	return parking, barcodes, nil
}

func (r *GeneralRoutesReporter) loadRemainsTares(ctx context.Context, dstOfficeIDs []int) (tares []*wb_models.TareForOffice, err error) {
	err = retryAction(ctx, "GeneralRoutesReporter.loadRemainsTares", 3, 1*time.Second, func() error {
		// isDrive = false is default value
		tares, err = r.services.WBLogisticService.GetTaresForOffices(ctx, r.officeID, dstOfficeIDs, false)
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "GeneralRoutesReporter.loadRemainsTares()", "failed load tares for offices")
	}
	return tares, nil
}

func (r *GeneralRoutesReporter) loadWaySheets(ctx context.Context) (waySheets []*wb_models.WaySheet, err error) {
	now := time.Now().In(time.UTC)
	for supplierID := range r.suppliers {
		var page *wb_models.WaySheetsPage
		err = retryAction(ctx, "GeneralRoutesReporter.loadWaySheets", 3, 1*time.Second, func() error {
			page, err = r.services.WBLogisticService.GetWaySheets(ctx, &models.WBLogisticGetWaySheetsParamsRequest{
				DateOpen:    now.AddDate(0, 0, -2),
				DateClose:   now,
				SupplierID:  supplierID,
				SrcOfficeID: r.officeID,
				Offset:      0,
				Limit:       500,
				WayTypeID:   0,
			})
			return err
		})
		if page == nil || err != nil {
			r.prompter.PromptError(fmt.Sprintf("Failed load way sheets for supplier %d", supplierID))
			logger.Logf(logger.ERROR, "GeneralRoutesReporter.loadWaySheets()", "failed load way sheets for supplier %d: %v", supplierID, err)
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
		return nil, errors.Wrap(err, "GeneralRoutesReporter.loadWaySheets()", "failed load way sheets")
	}
	return waySheets, nil
}

func (r *GeneralRoutesReporter) loadWaySheetInfo(ctx context.Context, waySheetID int) (waySheetInfo *wb_models.WaySheetInfo, err error) {
	err = retryAction(ctx, "GeneralRoutesReporter.loadWaySheetInfo", 3, 1*time.Second, func() error {
		waySheetInfo, err = r.services.WBLogisticService.GetWaySheetInfo(ctx, waySheetID)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "GeneralRoutesReporter.loadWaySheetInfo()", "failed load way sheet %d info", waySheetID)
	}

	if waySheetInfo == nil {
		return nil, errors.Newf("GeneralRoutesReporter.loadWaySheetInfo()", "way sheet %d info returned empty value without error", waySheetID)
	}

	return waySheetInfo, nil
}

func (r *GeneralRoutesReporter) sendReport(ctx context.Context, meta *reports.GeneralRoutesReportMetaData, reportData []*reports.GeneralRoutesReportData) error {
	report, err := r.reportSheet.Render(meta, reportData)
	if err != nil {
		return errors.Wrapf(err, "GeneralRoutesReporter.sendReport()", "failed render report")
	}

	if r.isRenderGS {
		data, err := r.rendererGS.Render(report)
		if err != nil {
			r.prompter.PromptError("Failed to render report Google Sheet")
			logger.Logf(logger.ERROR, "GeneralRoutesReporter.sendReport()", "failed render report for Google Sheets: %v", err)
		} else {
			err = r.sendGoogleSheets(ctx, data)
			if err != nil {
				r.prompter.PromptError("Failed to send report Google Sheet")
				logger.Logf(logger.ERROR, "GeneralRoutesReporter.sendReport()", "failed send report to Google Sheets: %v", err)
			} else {
				r.prompter.PromptSendReport()
				logger.Logf(logger.INFO, "GeneralRoutesReporter.sendReport()", "send report to Google Sheets")
			}
		}
	}

	return nil
}

func (r *GeneralRoutesReporter) sendGoogleSheets(ctx context.Context, data [][]interface{}) error {
	err := retryAction(ctx, "GeneralRoutesReporter.sendGoogleSheets()", 3, 1*time.Second, func() error {
		err := r.services.GoogleSheetsService.ClearValues(r.spreadsheetID, r.sheetName, "A:Z")
		if err != nil {
			return errors.Wrapf(err, "GeneralRoutesReporter.sendGoogleSheets()", "failed clear sheet %s, page %s", r.spreadsheetID, r.sheetName)
		}
		return r.services.GoogleSheetsService.UpdateValues(r.spreadsheetID, r.sheetName, r.sheetPosition, data, false)
	})
	if err != nil {
		return errors.Wrapf(err, "GeneralRoutesReporter.sendGoogleSheets()", "failed update sheet %s, page %s to position %s", r.spreadsheetID, r.sheetName, r.sheetPosition)
	}
	return nil
}
