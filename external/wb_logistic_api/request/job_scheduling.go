package request

import (
	"context"
	"fmt"
	"strconv"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// GetJobSchedulingRequest Retrieving a User's Scheduled Jobs
//
// URL: https://logistics.wb.ru/transport-planning-service/api/v1/planning/last-mile
type GetJobSchedulingRequest struct {
	BaseRequest
}

func NewGetJobsSchedulingRequest(client transport.HTTPClient, token string) *GetJobSchedulingRequest {
	r := &GetJobSchedulingRequest{BaseRequest: *NewRequestToken(client, "https://logistics.wb.ru/transport-planning-service/api/v1/planning/last-mile", token)}
	return r
}

func (r *GetJobSchedulingRequest) Do(ctx context.Context) (response response.GetJobsSchedulingResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetJobSchedulingRequest.Do: %s", err)
	}
	return
}

func (r *GetJobSchedulingRequest) SupplierID(id int) *GetJobSchedulingRequest {
	r.queryParameters.Set("supplier_id", strconv.Itoa(id))
	return r
}
