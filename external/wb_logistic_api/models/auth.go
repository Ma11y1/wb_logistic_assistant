package models

import (
	"errors"
	wb_logistic_errors "wb_logistic_assistant/external/wb_logistic_api/errors"
)

type AuthMeta struct {
	Code    wb_logistic_errors.ErrorType `json:"code"` // error code
	Message string                       `json:"message"`
}

type AuthCode struct {
	AuthMethod string `json:"auth_method"`
	Sticker    string `json:"sticker"`
	Ttl        int    `json:"ttl"` // seconds between requests
}

// AuthAccessToken Primary access token and refresh token. This access token is used to obtain the user's data token AuthSessionToken
type AuthAccessToken struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (t *AuthAccessToken) Validate() error {
	if t.TokenType == "" {
		return errors.New("token_type is empty")
	}
	if t.AccessToken == "" {
		return errors.New("access_token is empty")
	}
	if t.RefreshToken == "" {
		return errors.New("refresh_token is empty")
	}
	return nil
}

// AuthAccessTokenJWTDecode Object encoded into access token which is found in AuthAccessToken.AccessToken. Resulting from decoding the JWT access token
type AuthAccessTokenJWTDecode struct {
	IAt                int    `json:"iat"`     // token release time in format Unix
	Version            int    `json:"version"` // token version
	User               string `json:"user"`
	ShardKey           string `json:"shard_key"`
	ClientID           string `json:"client_id"`
	SessionID          string `json:"session_id"`
	UserRegistrationDt int    `json:"user_registration_dt"`
	ValidationKey      string `json:"validation_key"`
}

// AuthRefreshTokenJWTDecode Object encoded into refresh token which is found in AuthAccessToken.RefreshToken. Resulting from decoding the JWT refresh token
type AuthRefreshTokenJWTDecode struct {
	IAt                int    `json:"iat"`     // token release time in format Unix
	Version            int    `json:"version"` // token version
	User               string `json:"user"`
	ShardKey           string `json:"shard_key"`
	ClientID           string `json:"client_id"`
	SessionID          string `json:"session_id"`
	UserRegistrationDt int    `json:"user_registration_dt"`
}

// AuthSessionToken Merged token. User data token, used to retrieve data
type AuthSessionToken struct {
	Source string               `json:"source"`
	Token  AuthSessionTokenData `json:"token"`
}

func (t *AuthSessionToken) Validate() error {
	if t.Source == "" {
		return errors.New("source is empty")
	}
	err := t.Token.Validate()
	if err != nil {
		return err
	}
	return nil
}

type AuthSessionTokenData struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (t *AuthSessionTokenData) Validate() error {
	if t.TokenType == "" {
		return errors.New("token_type is empty")
	}
	if t.AccessToken == "" {
		return errors.New("access_token is empty")
	}
	if t.ExpiresIn == 0 {
		return errors.New("expires_in is zero")
	}
	return nil

}

// AuthSessionTokenJWTDecode Merged token. Resulting from decoding the JWT AuthSessionToken.Token.OAuthToken
type AuthSessionTokenJWTDecode struct {
	Exp          int                            `json:"exp"`
	ID           int                            `json:"id"`
	FreelancerID int                            `json:"freelancer_id"` // User ID
	EmployeeID   int                            `json:"employee_id"`
	IsFreelancer int                            `json:"is_freelancer"`
	SessionID    string                         `json:"session_id"`
	Roles        []string                       `json:"roles"`
	Extra        *AuthMergedTokenJWTDecodeExtra `json:"extra"`
}

type AuthMergedTokenJWTDecodeExtra struct {
	ID         string   `json:"wbID"`
	Name       string   `json:"name"`
	Phone      string   `json:"phone"`
	Position   string   `json:"position"`
	PositionID int      `json:"position_id"`
	Company    string   `json:"company"`
	Resources  []string `json:"resources"`
	Admin      []string `json:"admin"`
}
