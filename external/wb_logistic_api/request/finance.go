package request

import (
	"context"
	"fmt"
	"strconv"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// GetWaySheetFinanceDetailsRequest Obtaining financial information from a way sheet
//
// URL: https://drive.wb.ru/client-gateway/api/finance/credeber/v1/payment/details
type GetWaySheetFinanceDetailsRequest struct {
	BaseRequest
}

func NewGetWaySheetFinanceDetailsRequest(client transport.HTTPClient, token string) *GetWaySheetFinanceDetailsRequest {
	r := &GetWaySheetFinanceDetailsRequest{BaseRequest: *NewRequestToken(client, "https://drive.wb.ru/client-gateway/api/finance/credeber/v1/payment/details", token)}
	return r
}

func (r *GetWaySheetFinanceDetailsRequest) Do(ctx context.Context) (response response.GetWaySheetFinanceDetailsResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetWaySheetFinanceDetailsRequest.Do: %s", err)
	}
	return
}

func (r *GetWaySheetFinanceDetailsRequest) WaySheetID(id int) *GetWaySheetFinanceDetailsRequest {
	r.queryParameters.Set("waysheet", strconv.Itoa(id))
	return r
}
