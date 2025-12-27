package session

import (
	"context"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/external/wb_logistic_api/request"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

type AuthService struct {
	client transport.HTTPClient
}

func NewAuthService(client transport.HTTPClient) *AuthService {
	return &AuthService{client: client}
}

func (s *AuthService) RequestAuthCode(ctx context.Context, login string) (*models.AuthCode, error) {
	login, err := s.normalizeLogin(login)
	if err != nil {
		return nil, err
	}

	req := request.NewAuthGetCodeRequest(s.client).
		PhoneNumber(login)

	res, err := req.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get code by login '%s': %s", login, err)
	}
	meta := res.Meta
	apiError := res.Error

	if meta != nil && meta.Code != errors.ErrorTypeNone && (meta.Code != "" || meta.Message != "") {
		return nil, fmt.Errorf("failed get code by login '%s': [%s] %s", login, meta.Code, meta.Message)
	}

	if apiError != nil && apiError.Code != errors.ErrorTypeNone && (apiError.Code != "" || apiError.Message != "") {
		return nil, fmt.Errorf("failed get code by login'%s': [%s] %s", login, apiError.Code, apiError.Message)
	}

	return res.Data, err
}

// ExchangeCode Exchanges code for access token
func (s *AuthService) ExchangeCode(ctx context.Context, code int, sticker string) (*models.AuthAccessToken, error) {
	if code <= 0 {
		return nil, fmt.Errorf("invalid code %d", code)
	}
	if sticker == "" {
		return nil, fmt.Errorf("sticker is empty")
	}

	req := request.NewAuthRequest(s.client).
		Code(code).
		Sticker(sticker)

	res, err := req.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed exchange code by code '%d' and sticker '%s': %s", code, sticker, err)
	}

	meta := res.Meta
	apiError := res.Error
	if meta != nil && meta.Code != errors.ErrorTypeNone && (meta.Code != "" || meta.Message != "") {
		return nil, fmt.Errorf("failed exchange code by code '%d' and sticker '%s': [%s] %s", code, sticker, meta.Code, meta.Message)
	}

	if apiError != nil && apiError.Code != errors.ErrorTypeNone && (apiError.Code != "" || apiError.Message != "") {
		return nil, fmt.Errorf("failed exchange code by code '%d' and sticker '%s': [%s] %s", code, sticker, apiError.Code, apiError.Message)
	}

	return res.Data, nil
}

func (s *AuthService) GetSessionToken(ctx context.Context, login, accessToken string) (*models.AuthSessionToken, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is empty")
	}
	login, err := s.normalizeLogin(login)
	if err != nil {
		return nil, err
	}

	req := request.NewAuthMergeRequest(s.client).
		PhoneNumber(login).
		Token(accessToken)

	res, err := req.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed merge token by login '%s': %s", login, err)
	}

	apiError := res.Error
	if apiError != nil && apiError.Code != errors.ErrorTypeNone && (apiError.Code != "" || apiError.Message != "") {
		return nil, fmt.Errorf("failed merge token by login '%s': [%s] %s", login, apiError.Code, apiError.Message)
	}

	return res.Data, nil
}

// RefreshAccessToken TODO It is not known how to update tokens, there is no documentation!
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*models.AuthAccessToken, error) {
	return nil, fmt.Errorf("AuthService.RefreshAccessToken(): not implemented yet")
}

// RefreshAccessToken TODO It is not known how to update tokens, there is no documentation!
func (s *AuthService) RefreshSession(ctx context.Context, refreshToken string) (*models.AuthAccessToken, error) {
	return nil, fmt.Errorf("AuthService.RefreshAccessToken(): not implemented yet")
}

func (s *AuthService) GetUserInfo(ctx context.Context, sessionToken string) (*models.UserInfo, error) {
	if sessionToken == "" {
		return nil, fmt.Errorf("session token is empty")
	}

	decodedToken := &models.AuthSessionTokenJWTDecode{}
	err := decodeJWT(sessionToken, decodedToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT session token: %w", err)
	}

	req := request.NewUserGetInfoRequest(s.client, sessionToken).
		ClientID(decodedToken.FreelancerID)

	res, err := req.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get user info by user ID %d: %w", decodedToken.FreelancerID, err)
	}
	if res.Error != nil && (res.Error.Code > 0 || res.Error.Err != "") {
		return nil, fmt.Errorf("failed to get user info by user ID %d: %d %s", decodedToken.FreelancerID, res.Error.Code, res.Error.Error())
	}

	return res.Data, nil
}

func (s *AuthService) normalizeLogin(login string) (string, error) {
	if len(login) < 10 {
		return "", fmt.Errorf("login '%s' is too short, len: %d", login, len(login))
	}
	if login[0] == '+' {
		login = login[1:]
	}
	return login, nil
}
