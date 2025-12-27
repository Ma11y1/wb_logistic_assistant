package request

import (
	"context"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// GetAssociationOfficesInfoByNameRequest Getting information about existing offices by name or part of name
//
// URL: https://logistics.wb.ru/routes-netcore-service/api/v1/office
type GetAssociationOfficesInfoByNameRequest struct {
	BaseRequest
}

func NewGetAssociationOfficesInfoByNameRequest(client transport.HTTPClient, token string) *GetAssociationOfficesInfoByNameRequest {
	r := &GetAssociationOfficesInfoByNameRequest{BaseRequest: *NewRequestToken(client, "https://logistics.wb.ru/routes-netcore-service/api/v1/office", token)}
	r.IsDc(true)
	return r
}

func (r *GetAssociationOfficesInfoByNameRequest) Do(ctx context.Context) (response response.GetAssociationOfficesInfoByNameResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetAssociationOfficesInfoByNameRequest.Do: %s", err)
	}
	return
}

func (r *GetAssociationOfficesInfoByNameRequest) Param(param string) *GetAssociationOfficesInfoByNameRequest {
	r.queryParameters.Set("param", param)
	return r
}

func (r *GetAssociationOfficesInfoByNameRequest) IsDc(v bool) *GetAssociationOfficesInfoByNameRequest {
	r.queryParameters.Set("isDc", v)
	return r
}
