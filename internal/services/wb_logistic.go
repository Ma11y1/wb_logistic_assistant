package services

import (
	"context"
	"wb_logistic_assistant/external/wb_logistic_api"
	wb_models "wb_logistic_assistant/external/wb_logistic_api/models"
	wb_logistic_session "wb_logistic_assistant/external/wb_logistic_api/session"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/models"
)

type WBLogisticService interface {
	GetUserInfo(ctx context.Context) (*wb_models.UserInfo, error)
	GetRemainsLastMileReports(ctx context.Context) (wb_models.RemainsLastMileReports, error)
	GetRemainsLastMileReportByOfficeID(ctx context.Context, officeID int) (*wb_models.RemainsLastMileReport, error)
	GetRemainsLastMileReportsRouteInfo(ctx context.Context, routeID int) ([]*wb_models.RemainsLastMileReportsRouteInfo, error)
	GetJobsScheduling(ctx context.Context) (*wb_models.JobsScheduling, error)
	GetShipments(ctx context.Context, params *models.WBLogisticGetShipmentsParamsRequest) ([]*wb_models.Shipment, int, error)
	GetShipmentInfo(ctx context.Context, shipmentID int) (*wb_models.ShipmentInfo, error)
	GetShipmentTransfers(ctx context.Context, shipmentID int) (*wb_models.ShipmentTransfers, error)
	GetTaresForOffices(ctx context.Context, officeID int, dstOfficeIDs []int, isDrive bool) ([]*wb_models.TareForOffice, error)
}
type BaseWBLogisticService struct {
	client                               *wb_logistic_api.Client
	session                              *wb_logistic_session.Session
	cacheUserInfo                        *Cache[*wb_models.UserInfo]
	cacheRemainsLastMileReports          *Cache[wb_models.RemainsLastMileReports]
	cacheRemainsLastMileReport           *GenericMapCache[int, *wb_models.RemainsLastMileReport]
	cacheRemainsLastMileReportsRouteInfo *GenericMapCache[int, []*wb_models.RemainsLastMileReportsRouteInfo]
	cacheJobsScheduling                  *Cache[*wb_models.JobsScheduling]
	cacheShipmentInfo                    *GenericMapCache[int, *wb_models.ShipmentInfo]
	cacheShipmentTransfers               *GenericMapCache[int, *wb_models.ShipmentTransfers]
	cacheWaySheetInfo                    *GenericMapCache[int, *wb_models.WaySheetInfo]
}

func NewBaseWBLogisticService(client *wb_logistic_api.Client, session *wb_logistic_session.Session, ttl *models.WBLogisticTTlParams) *BaseWBLogisticService {
	var s *BaseWBLogisticService
	s = &BaseWBLogisticService{
		client:  client,
		session: session,
		cacheUserInfo: NewCache[*wb_models.UserInfo](ttl.UserInfo, func(ctx context.Context) (*wb_models.UserInfo, error) {
			res, err := client.GetUserInfo(ctx, session)
			if err != nil {
				return nil, errors.Wrap(err, "BaseWBLogisticService.GetUserInfo()", "")
			}
			return res.Data, nil
		}),
		cacheRemainsLastMileReports: NewCache[wb_models.RemainsLastMileReports](ttl.RemainsLastMileReports, func(ctx context.Context) (wb_models.RemainsLastMileReports, error) {
			res, err := client.GetRemainsLastMileReports(ctx, session)
			if err != nil {
				return nil, errors.Wrap(err, "BaseWBLogisticService.GetRemainsLastMileReports()", "")
			}
			return res.Data, nil
		}),
		cacheRemainsLastMileReport: NewGenericMapCache[int, *wb_models.RemainsLastMileReport](ttl.RemainsLastMileReports, func(ctx context.Context, officeID int) (*wb_models.RemainsLastMileReport, error) {
			reports, err := s.GetRemainsLastMileReports(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "BaseWBLogisticService.GetRemainsLastMileReportByOfficeID()", "")
			}

			for _, report := range reports {
				if report.OfficeID == officeID {
					return report, nil
				}
			}

			return nil, errors.New("BaseWBLogisticService.GetRemainsLastMileReportByOfficeID()", "route report not found")
		}),
		cacheRemainsLastMileReportsRouteInfo: NewGenericMapCache[int, []*wb_models.RemainsLastMileReportsRouteInfo](ttl.RemainsLastMileReportsRouteInfo, func(ctx context.Context, routeID int) ([]*wb_models.RemainsLastMileReportsRouteInfo, error) {
			info, err := client.GetRemainsLastMileReportsRouteInfo(ctx, session, routeID)
			if err != nil {
				return nil, errors.Wrap(err, "BaseWBLogisticService.GetRemainsLastMileReportsRouteInfo()", "")
			}
			return info.Data, nil
		}),
		cacheJobsScheduling: NewCache[*wb_models.JobsScheduling](ttl.JobsScheduling, func(ctx context.Context) (*wb_models.JobsScheduling, error) {
			res, err := s.client.GetJobsScheduling(ctx, s.session)
			if err != nil {
				return nil, errors.Wrap(err, "BaseWBLogisticService.GetJobsScheduling()", "")
			}
			if res.Data.Route == nil || len(res.Data.Route) == 0 {
				return nil, errors.New("BaseWBLogisticService.GetJobsScheduling()", "jobs scheduling routes not found")
			}
			return res.Data, nil
		}),
		cacheShipmentInfo: NewGenericMapCache[int, *wb_models.ShipmentInfo](ttl.ShipmentInfo, func(ctx context.Context, shipmentID int) (*wb_models.ShipmentInfo, error) {
			info, err := client.GetShipmentInfo(ctx, session, shipmentID)
			if err != nil {
				return nil, errors.Wrap(err, "BaseWBLogisticService.GetShipmentInfo()", "")
			}
			return info.Data, nil
		}),
		cacheShipmentTransfers: NewGenericMapCache[int, *wb_models.ShipmentTransfers](ttl.ShipmentTransfers, func(ctx context.Context, shipmentID int) (*wb_models.ShipmentTransfers, error) {
			info, err := client.GetShipmentTransfers(ctx, session, shipmentID)
			if err != nil {
				return nil, errors.Wrap(err, "BaseWBLogisticService.GetShipmentTransfers()", "")
			}
			return info.Data, nil
		}),
		cacheWaySheetInfo: NewGenericMapCache[int, *wb_models.WaySheetInfo](ttl.WaySheetInfo, func(ctx context.Context, id int) (*wb_models.WaySheetInfo, error) {
			info, err := client.GetWaySheetInfo(ctx, session, id)
			if err != nil {
				return nil, errors.Wrap(err, "BaseWBLogisticService.GetWaySheetInfo()", "")
			}
			return info.Data.WaySheet, nil
		}),
	}
	return s
}

func (s *BaseWBLogisticService) GetUserInfo(ctx context.Context) (*wb_models.UserInfo, error) {
	return s.cacheUserInfo.Get(ctx)
}

func (s *BaseWBLogisticService) GetRemainsLastMileReports(ctx context.Context) (wb_models.RemainsLastMileReports, error) {
	return s.cacheRemainsLastMileReports.Get(ctx)
}

func (s *BaseWBLogisticService) GetRemainsLastMileReportByOfficeID(ctx context.Context, officeID int) (*wb_models.RemainsLastMileReport, error) {
	return s.cacheRemainsLastMileReport.Get(ctx, officeID)
}

func (s *BaseWBLogisticService) GetRemainsLastMileReportsRouteInfo(ctx context.Context, routeID int) ([]*wb_models.RemainsLastMileReportsRouteInfo, error) {
	return s.cacheRemainsLastMileReportsRouteInfo.Get(ctx, routeID)
}

func (s *BaseWBLogisticService) GetJobsScheduling(ctx context.Context) (*wb_models.JobsScheduling, error) {
	return s.cacheJobsScheduling.Get(ctx)
}

// GetShipments Returns shipments list, shipments total count, error
func (s *BaseWBLogisticService) GetShipments(ctx context.Context, params *models.WBLogisticGetShipmentsParamsRequest) ([]*wb_models.Shipment, int, error) {
	if params == nil {
		return nil, 0, errors.New("BaseWBLogisticService.GetShipments()", "params is nil")
	}

	apiParams := &wb_models.GetShipmentParamsRequest{
		DataStart:                params.DataStart,
		DataEnd:                  params.DataEnd,
		SrcOfficeID:              params.SrcOfficeID,
		PageIndex:                params.PageIndex,
		Limit:                    params.Limit,
		SupplierID:               params.SupplierID,
		Direction:                params.Direction,
		Sorter:                   params.Sorter,
		FilterShipmentID:         params.FilterShipmentID,
		FilterShipmentType:       params.FilterShipmentType,
		FilterVehicleNumberPlate: params.FilterVehicleNumberPlate,
		FilterDstOffice:          params.FilterDstOffice,
		FilterShowOnlyOpen:       params.FilterShowOnlyOpen,
	}

	res, err := s.client.GetShipments(ctx, s.session, apiParams)
	if err != nil {
		return nil, 0, errors.Wrap(err, "BaseWBLogisticService.GetShipments()", "")
	}
	return res.Data, res.Meta.TotalCount, nil
}

func (s *BaseWBLogisticService) GetShipmentInfo(ctx context.Context, shipmentID int) (*wb_models.ShipmentInfo, error) {
	return s.cacheShipmentInfo.Get(ctx, shipmentID)
}

func (s *BaseWBLogisticService) GetShipmentTransfers(ctx context.Context, shipmentID int) (*wb_models.ShipmentTransfers, error) {
	return s.cacheShipmentTransfers.Get(ctx, shipmentID)
}

func (s *BaseWBLogisticService) GetTaresForOffices(ctx context.Context, officeID int, dstOfficeIDs []int, isDrive bool) ([]*wb_models.TareForOffice, error) {
	if officeID <= 0 {
		return nil, errors.New("BaseWBLogisticService.GetTaresForOffices()", "office id is invalid")
	}
	if len(dstOfficeIDs) == 0 {
		return nil, errors.New("BaseWBLogisticService.GetTaresForOffices()", "destination office id's is empty")
	}

	res, err := s.client.GetTaresForOffices(ctx, s.session, officeID, dstOfficeIDs, isDrive)
	if err != nil {
		return nil, errors.Wrap(err, "BaseWBLogisticService.GetTaresForOffices()", "")
	}
	return res.Data, nil
}

func (s *BaseWBLogisticService) GetWaySheets(ctx context.Context, params *models.WBLogisticGetWaySheetsParamsRequest) (*wb_models.WaySheetsPage, error) {
	if params == nil {
		return nil, errors.New("BaseWBLogisticService.GetWaySheets()", "params is nil")
	}

	apiParams := &wb_models.GetWaySheetsParamsRequest{
		DateOpen:    params.DateOpen,
		DateClose:   params.DateClose,
		SupplierID:  params.SupplierID,
		SrcOfficeID: params.SrcOfficeID,
		Limit:       params.Limit,
		Offset:      params.Offset,
		WayTypeID:   params.WayTypeID,
	}

	res, err := s.client.GetWaySheets(ctx, s.session, apiParams)
	if err != nil {
		return nil, errors.Wrap(err, "BaseWBLogisticService.GetWaySheets()", "")
	}
	return res.Data, nil
}

func (s *BaseWBLogisticService) GetWaySheetInfo(ctx context.Context, id int) (*wb_models.WaySheetInfo, error) {
	return s.cacheWaySheetInfo.Get(ctx, id)
}
