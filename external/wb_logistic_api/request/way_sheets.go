package request

import (
	"context"
	"fmt"
	"time"
	"wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// GetWaySheetsRequest Get list of way sheets
//
// URL: https://drive.wb.ru/client-gateway/api/waysheets/v1/waysheets
type GetWaySheetsRequest struct {
	BaseRequest
}

func NewGetWaySheetsRequestRequest(client transport.HTTPClient, token string) *GetWaySheetsRequest {
	r := &GetWaySheetsRequest{BaseRequest: *NewRequestToken(client, "https://drive.wb.ru/client-gateway/api/waysheets/v1/waysheets", token)}
	r.Limit(10)
	r.Offset(0)
	r.WayTypeID(0)
	return r
}

func (r *GetWaySheetsRequest) Do(ctx context.Context) (response response.GetWaySheetsResponse, err error) {
	err = r.PostUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetWaySheetsRequest.Do: %s", err)
	}
	return
}

func (r *GetWaySheetsRequest) DateOpen(t time.Time) *GetWaySheetsRequest {
	r.parameters.Set("date_open", t.Format("2006-01-02T15:04:05.000Z"))
	return r
}

func (r *GetWaySheetsRequest) DateClose(t time.Time) *GetWaySheetsRequest {
	r.parameters.Set("date_close", t.Format("2006-01-02T15:04:05.000Z"))
	return r
}

func (r *GetWaySheetsRequest) Limit(limit int) *GetWaySheetsRequest {
	r.parameters.Set("limit", limit)
	return r
}

func (r *GetWaySheetsRequest) Offset(offset int) *GetWaySheetsRequest {
	r.parameters.Set("offset", offset)
	return r
}

func (r *GetWaySheetsRequest) SrcOfficeID(id int) *GetWaySheetsRequest {
	r.parameters.Set("src_office_id", id)
	return r
}

func (r *GetWaySheetsRequest) SupplierID(id int) *GetWaySheetsRequest {
	r.parameters.Set("supplier_id", id)
	return r
}

func (r *GetWaySheetsRequest) WayTypeID(id int) *GetWaySheetsRequest {
	r.parameters.Set("way_type_id", id)
	return r
}

func (r *GetWaySheetsRequest) RouteCarID(id int) *GetWaySheetsRequest {
	r.parameters.Set("routecar_id", id)
	return r
}

func (r *GetWaySheetsRequest) VehicleNumberPlate(number string) *GetWaySheetsRequest {
	r.parameters.Set("vehicle_number_plate", number)
	return r
}

func (r *GetWaySheetsRequest) FromParams(params *models.GetWaySheetsParamsRequest) *GetWaySheetsRequest {
	if params == nil {
		return r
	}
	r.DateOpen(params.DateOpen)
	r.DateClose(params.DateClose)
	if params.Limit > 0 {
		r.Limit(params.Limit)
	}
	if params.Offset > 0 {
		r.Offset(params.Offset)
	}
	if params.SrcOfficeID > 0 {
		r.SrcOfficeID(params.SrcOfficeID)
	}
	if params.SupplierID > 0 {
		r.SupplierID(params.SupplierID)
	}
	if params.WayTypeID > 0 {
		r.WayTypeID(params.WayTypeID)
	}
	if params.RouteCarID > 0 {
		r.RouteCarID(params.RouteCarID)
	}
	if params.VehicleNumberPlate != "" {
		r.VehicleNumberPlate(params.VehicleNumberPlate)
	}
	return r
}

func (r *GetWaySheetsRequest) ClearParams() {
	r.parameters.Remove("date_close")
	r.parameters.Remove("date_close")
	r.parameters.Remove("limit")
	r.parameters.Remove("offset")
	r.parameters.Remove("src_office_id")
	r.parameters.Remove("supplier_id")
	r.parameters.Remove("way_type_id")
	r.parameters.Remove("routecar_id")
	r.parameters.Remove("vehicle_number_plate")
}

// GetWaySheetInfoRequest
//
// URL: https://drive.wb.ru/client-gateway/api/waysheets/v1/waysheets/[WAYSHEET_ID]
type GetWaySheetInfoRequest struct {
	originURL string
	BaseRequest
}

func NewGetWaySheetInfoRequest(client transport.HTTPClient, token string) *GetWaySheetInfoRequest {
	url := "https://drive.wb.ru/client-gateway/api/waysheets/v1/waysheets/"
	return &GetWaySheetInfoRequest{
		originURL:   url,
		BaseRequest: *NewRequestToken(client, url, token),
	}
}

func (r *GetWaySheetInfoRequest) Do(ctx context.Context) (response response.GetWaySheetInfoResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetWaySheetInfoRequest.Do: %s", err)
	}
	return
}

func (r *GetWaySheetInfoRequest) WaySheetID(id int) *GetWaySheetInfoRequest {
	r.url = fmt.Sprintf("%s%d", r.originURL, id)
	return r
}
