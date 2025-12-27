package request

import (
	"context"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// UserGetInfoRequest Exchanges authorization code for access tokens
//
// URL: https://logistics.wb.ru/api/v1/public/iam/accounts/[USER_ID]
type UserGetInfoRequest struct {
	originURL string
	BaseRequest
}

func NewUserGetInfoRequest(client transport.HTTPClient, token string) *UserGetInfoRequest {
	url := "https://logistics.wb.ru/api/v1/public/iam/accounts/"
	return &UserGetInfoRequest{
		originURL:   url,
		BaseRequest: *NewRequestToken(client, url, token),
	}
}

func (r *UserGetInfoRequest) Do(ctx context.Context) (response response.UserGetInfoResponse, err error) {
	err = r.GetUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("UserGetInfoRequest.Do: %s", err)
	}
	return
}

func (r *UserGetInfoRequest) ClientID(id int) *UserGetInfoRequest {
	r.url = fmt.Sprintf("%s%d", r.originURL, id)
	return r
}
