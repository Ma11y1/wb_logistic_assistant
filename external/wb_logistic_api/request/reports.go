package request

import (
	"context"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// GetRemainsLastMileReportsRequest Returns reports on existing routes
//
// URL: https://logistics.wb.ru/reports-service/api/v1/last-mile
type GetRemainsLastMileReportsRequest struct {
	BaseRequest
}

func NewGetRemainsLastMileReportsRequest(client transport.HTTPClient, token string) *GetRemainsLastMileReportsRequest {
	return &GetRemainsLastMileReportsRequest{BaseRequest: *NewRequestToken(client, "https://logistics.wb.ru/reports-service/api/v1/last-mile", token)}
}

func (r *GetRemainsLastMileReportsRequest) Do(ctx context.Context) (response response.GetRemainsLastMileReportsResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetRemainsLastMileReportsRequest.Do: %s", err)
	}
	return
}

// GetRemainsLastMileReportInfoRequest Returns report info by report id
//
// URL: https://logistics.wb.ru/reports-service/api/v1/last-mile/[report id]
type GetRemainsLastMileReportInfoRequest struct {
	BaseRequest
	originURL string
}

func NewGetRemainsLastMileReportsInfoRequest(client transport.HTTPClient, token string) *GetRemainsLastMileReportInfoRequest {
	url := "https://logistics.wb.ru/reports-service/api/v1/last-mile/"
	r := &GetRemainsLastMileReportInfoRequest{
		originURL:   url,
		BaseRequest: *NewRequestToken(client, url, token),
	}
	return r
}

func (r *GetRemainsLastMileReportInfoRequest) Do(ctx context.Context) (response response.GetRemainsLastMileReportsRouteInfoResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetRemainsLastMileReportInfoRequest.Do: %s", err)
	}
	return
}

func (r *GetRemainsLastMileReportInfoRequest) RouteID(id int) *GetRemainsLastMileReportInfoRequest {
	r.url = fmt.Sprintf("%s%d", r.originURL, id)
	return r
}
