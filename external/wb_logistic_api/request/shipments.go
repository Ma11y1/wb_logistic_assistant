package request

import (
	"context"
	"fmt"
	"time"
	"wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// GetShipmentsRequest Returns shipments
//
// URL: https://logistics.wb.ru/shipments-service/api/v1/shipments
type GetShipmentsRequest struct {
	BaseRequest
}

func NewGetShipmentsRequest(client transport.HTTPClient, token string) *GetShipmentsRequest {
	r := &GetShipmentsRequest{BaseRequest: *NewRequestToken(client, "https://logistics.wb.ru/shipments-service/api/v1/shipments", token)}
	r.queryParameters.Set("dt_start", time.Now().Format("2006-01-02"))
	r.queryParameters.Set("dt_end", time.Now().Format("2006-01-02"))
	r.queryParameters.Set("page_index", 0)
	r.queryParameters.Set("limit", 50)
	r.queryParameters.Set("direction", -1)
	r.queryParameters.Set("sorter", "updated_at")
	return r
}

func (r *GetShipmentsRequest) Do(ctx context.Context) (response response.GetShipmentsResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetShipmentsRequest.Do: %s", err)
	}
	return
}

func (r *GetShipmentsRequest) DateStart(date time.Time) *GetShipmentsRequest {
	r.queryParameters.Set("dt_start", date.Format("2006-01-02"))
	return r
}

func (r *GetShipmentsRequest) DateEnd(date time.Time) *GetShipmentsRequest {
	r.queryParameters.Set("dt_end", date.Format("2006-01-02"))
	return r
}

func (r *GetShipmentsRequest) OfficeID(id int) *GetShipmentsRequest {
	r.queryParameters.Set("src_office_id", id)
	return r
}

func (r *GetShipmentsRequest) PageIndex(index int) *GetShipmentsRequest {
	r.queryParameters.Set("page_index", index)
	return r
}

func (r *GetShipmentsRequest) Limit(limit int) *GetShipmentsRequest {
	r.queryParameters.Set("limit", limit)
	return r
}

func (r *GetShipmentsRequest) SupplierID(id int) *GetShipmentsRequest {
	r.queryParameters.Set("supplier_id", id)
	return r
}

// Direction Sets the sort direction. Types: 1 or -1
func (r *GetShipmentsRequest) Direction(direction int) *GetShipmentsRequest {
	r.queryParameters.Set("direction", direction)
	return r
}

// Sorter Specifies the sort type. Types: "id", "supplier", "driver", "route_car_id", "updated_at"
func (r *GetShipmentsRequest) Sorter(t string) *GetShipmentsRequest {
	r.queryParameters.Set("sorter", t)
	return r
}

func (r *GetShipmentsRequest) ClearFilter() *GetShipmentsRequest {
	r.queryParameters.Remove("shipment_id")
	r.queryParameters.Remove("shipment_type")
	r.queryParameters.Remove("vehicle_number_plate")
	r.queryParameters.Remove("dst_office")
	r.queryParameters.Remove("show_only_open")
	return r
}

// FilterShipmentID Filter by shipment id
func (r *GetShipmentsRequest) FilterShipmentID(id int) *GetShipmentsRequest {
	r.queryParameters.Set("shipment_id", id)
	return r
}

// FilterShipmentType Filter by shipment type. Types: "last-mile" or "truck"
func (r *GetShipmentsRequest) FilterShipmentType(t string) *GetShipmentsRequest {
	r.queryParameters.Set("shipment_type", t)
	return r
}

// FilterVehicleNumberPlate Filter by Vehicle number plate
func (r *GetShipmentsRequest) FilterVehicleNumberPlate(number int) *GetShipmentsRequest {
	r.queryParameters.Set("vehicle_number_plate", number)
	return r
}

// FilterDstOffice Filter by destination office
func (r *GetShipmentsRequest) FilterDstOffice(v interface{}) *GetShipmentsRequest {
	r.queryParameters.Set("dst_office", v)
	return r
}

// FilterShowOnlyOpen Filter by show only open shipments
func (r *GetShipmentsRequest) FilterShowOnlyOpen(v bool) *GetShipmentsRequest {
	if v {
		r.queryParameters.Set("show_only_open", v)
	} else {
		r.queryParameters.Remove("show_only_open")
	}
	return r
}

func (r *GetShipmentsRequest) FromParams(params *models.GetShipmentParamsRequest) *GetShipmentsRequest {
	if params == nil {
		return r
	}
	r.DateStart(params.DataStart)
	r.DateEnd(params.DataEnd)
	if params.SrcOfficeID > 0 {
		r.OfficeID(params.SrcOfficeID)
	}
	if params.PageIndex >= 0 {
		r.PageIndex(params.PageIndex)
	}
	if params.Limit > 0 {
		r.Limit(params.Limit)
	}
	if params.SupplierID > 0 {
		r.SupplierID(params.SupplierID)
	}
	if params.Direction == -1 || params.Direction == 1 {
		r.Direction(params.Direction)
	}
	if params.Sorter != "" {
		r.Sorter(params.Sorter)
	}
	if params.FilterShipmentID > 0 {
		r.FilterShipmentID(params.FilterShipmentID)
	}
	if params.FilterShipmentType != "" {
		r.FilterShipmentType(params.FilterShipmentType)
	}
	if params.FilterVehicleNumberPlate > 0 {
		r.FilterVehicleNumberPlate(params.FilterVehicleNumberPlate)
	}
	if params.FilterDstOffice != "" {
		r.FilterDstOffice(params.FilterDstOffice)
	}
	if params.FilterShowOnlyOpen {
		r.FilterShowOnlyOpen(params.FilterShowOnlyOpen)
	}
	return r
}

// GetShipmentInfoRequest Return shipment info by shipment id
//
// URL: https://logistics.wb.ru/shipments-service/api/v1/shipments/[SHIPMENT_ID]/info
type GetShipmentInfoRequest struct {
	originURL string
	BaseRequest
}

func NewGetShipmentInfoRequest(client transport.HTTPClient, token string) *GetShipmentInfoRequest {
	originURL := "https://logistics.wb.ru/shipments-service/api/v1/shipments/"
	r := &GetShipmentInfoRequest{BaseRequest: *NewRequestToken(client, originURL, token)}
	r.originURL = originURL
	return r
}

func (r *GetShipmentInfoRequest) Do(ctx context.Context) (response response.GetShipmentInfoResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetShipmentInfoRequest.Do: %s", err)
	}
	return
}

func (r *GetShipmentInfoRequest) ShipmentID(id int) *GetShipmentInfoRequest {
	r.url = fmt.Sprintf("%s%d/info", r.originURL, id)
	return r
}

// GetShipmentTransfersRequest Return shipment info by shipment id
//
// URL: https://logistics.wb.ru/shipments-service/api/v1/shipments/[SHIPMENT_ID]/transfers

type GetShipmentTransfersRequest struct {
	originURL string
	BaseRequest
}

func NewGetShipmentTransfersRequest(client transport.HTTPClient, token string) *GetShipmentTransfersRequest {
	originURL := "https://logistics.wb.ru/shipments-service/api/v1/shipments/"
	r := &GetShipmentTransfersRequest{BaseRequest: *NewRequestToken(client, originURL, token)}
	r.originURL = originURL
	return r
}

func (r *GetShipmentTransfersRequest) Do(ctx context.Context) (response response.GetShipmentTransfersResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetShipmentTransfersRequest.Do: %s", err)
	}
	return
}

func (r *GetShipmentTransfersRequest) ShipmentID(id int) *GetShipmentTransfersRequest {
	r.url = fmt.Sprintf("%s%d/transfers", r.originURL, id)
	return r
}
