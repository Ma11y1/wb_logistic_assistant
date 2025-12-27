package request

import (
	"context"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// GetAssociationRoutesInfoByNameRequest Getting information about existing routes by name or part of the name
//
// URL: https://logistics.wb.ru/routes-netcore-service/api/v1/route/by-param/false/true/ // straight params
type GetAssociationRoutesInfoByNameRequest struct {
	originURL string
	BaseRequest
}

func NewGetAssociationRoutesInfoByNameRequest(client transport.HTTPClient, token string) *GetAssociationRoutesInfoByNameRequest {
	originURL := "https://logistics.wb.ru/routes-netcore-service/api/v1/route/by-param/false/true/"
	r := &GetAssociationRoutesInfoByNameRequest{BaseRequest: *NewRequestToken(client, originURL, token)}
	r.originURL = originURL
	return r
}

func (r *GetAssociationRoutesInfoByNameRequest) Do(ctx context.Context) (response response.GetAssociationRoutesInfoByNameResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetAssociationRoutesInfoByNameRequest.Do: %s", err)
	}
	return
}

func (r *GetAssociationRoutesInfoByNameRequest) Param(param string) *GetAssociationRoutesInfoByNameRequest {
	r.SetURL(fmt.Sprintf("%s%s", r.originURL, param))
	return r
}
