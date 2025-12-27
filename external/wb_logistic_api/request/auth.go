package request

import (
	"context"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// AuthGetCodeRequest Sends a request with a phone number to receive an authorization code
//
// URL: https://drive.wb.ru/user-management/api/v1/public/registration/code
type AuthGetCodeRequest struct {
	BaseRequest
}

func NewAuthGetCodeRequest(client transport.HTTPClient) *AuthGetCodeRequest {
	req := &AuthGetCodeRequest{BaseRequest: *NewRequest(client, "https://drive.wb.ru/user-management/api/v1/public/registration/code")}
	req.header.Set("X-App-Type", "web")
	req.header.Set("X-Auth-Provider", "wb")
	return req
}

func (r *AuthGetCodeRequest) Do(ctx context.Context) (response response.AuthCodeResponse, err error) {
	err = r.PostUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("AuthGetCodeRequest.Do: %s", err)
	}
	return
}

func (r *AuthGetCodeRequest) CaptchaToken(token string) *AuthGetCodeRequest {
	r.parameters.Set("captcha_token", token)
	return r
}

// PhoneNumber Example: 79998882233
func (r *AuthGetCodeRequest) PhoneNumber(number string) *AuthGetCodeRequest {
	r.parameters.Set("phone_number", number)
	return r
}

// AuthRequest Exchanges authorization code for receive access tokens
//
// URL: https://drive.wb.ru/user-management/api/v1/public/registration/auth
type AuthRequest struct {
	BaseRequest
}

func NewAuthRequest(client transport.HTTPClient) *AuthRequest {
	req := &AuthRequest{BaseRequest: *NewRequest(client, "https://drive.wb.ru/user-management/api/v1/public/registration/auth")}
	req.header.Set("X-App-Type", "web")
	req.header.Set("X-Auth-Provider", "wb")
	req.TokenType("employee")
	return req
}

func (r *AuthRequest) Do(ctx context.Context) (response response.AuthResponse, err error) {
	err = r.PostUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("AuthGetCodeRequest.Do: %s", err)
	}
	return
}

func (r *AuthRequest) Code(code int) *AuthRequest {
	r.parameters.Set("code", code)
	return r
}

func (r *AuthRequest) Sticker(sticker string) *AuthRequest {
	r.parameters.Set("sticker", sticker)
	return r
}

func (r *AuthRequest) TokenType(t string) *AuthRequest {
	r.parameters.Set("token_type", t)
	return r
}

// AuthMergeRequest
//
// URL: https://drive.wb.ru/user-management/api/v1/public/token/merge
type AuthMergeRequest struct {
	BaseRequest
}

func NewAuthMergeRequest(client transport.HTTPClient) *AuthMergeRequest {
	req := &AuthMergeRequest{BaseRequest: *NewRequest(client, "https://drive.wb.ru/user-management/api/v1/public/token/merge")}
	return req
}

func (r *AuthMergeRequest) Do(ctx context.Context) (response response.AuthMergeResponse, err error) {
	err = r.PostUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("AuthMergeRequest.Do: %s", err)
	}
	return
}

func (r *AuthMergeRequest) Token(token string) *AuthMergeRequest {
	r.parameters.Set("authv3_token", token)
	return r
}

func (r *AuthMergeRequest) PhoneNumber(number string) *AuthMergeRequest {
	r.parameters.Set("phone", number)
	return r
}
