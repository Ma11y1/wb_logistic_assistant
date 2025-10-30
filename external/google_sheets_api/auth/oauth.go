package auth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"strings"
	"sync"
	models2 "wb_logistic_assistant/external/google_sheets_api/models"
)

//// Test applications receive a refresh token for 7 days only

type OAuthActor struct {
	mtx        sync.RWMutex
	httpClient *http.Client
	service    *sheets.Service
	token      *oauth2.Token
	config     *oauth2.Config
}

func (*OAuthActor) Type() ActorType {
	return OAuthActorType
}

// AuthCodeURL To start authorization, you need to call this method
func (a *OAuthActor) AuthCodeURL(credential *models2.OAuthCredentials, scope Scope, state string) (string, error) {
	if err := credential.Validate(); err != nil {
		return "", fmt.Errorf("invalid credential: %w", err)
	}

	credentialJSON, err := credential.ToJSON()
	if err != nil {
		return "", fmt.Errorf("error marshalling credential: %w", err)
	}

	config, err := google.ConfigFromJSON(credentialJSON, string(scope))
	if err != nil {
		return "", fmt.Errorf("failed to create config from JSON: %w", err)
	}

	a.mtx.Lock()
	a.config = config
	a.mtx.Unlock()

	return config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompters", "consent"),
	), nil
}

func (a *OAuthActor) ExchangeCode(ctx context.Context, code string) error {
	return a.ExchangeCodeTokenSource(ctx, code, &TokenSource{})
}

func (a *OAuthActor) ExchangeCodeTokenSource(ctx context.Context, code string, tokenSource *TokenSource) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.config == nil {
		return fmt.Errorf("authorization process has not started, you must first request an authorization code")
	}
	if code == "" {
		return fmt.Errorf("authorization code is empty")
	}
	if tokenSource == nil {
		return fmt.Errorf("token source is nil")
	}

	newToken, err := a.config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange session code for token: %w", err)
	}

	if newToken.RefreshToken == "" && a.token != nil {
		newToken.RefreshToken = a.token.RefreshToken
	}
	a.token = newToken

	baseSource := a.config.TokenSource(ctx, a.token)
	tokenSource.TokenSource = oauth2.ReuseTokenSource(a.token, baseSource)
	a.httpClient = oauth2.NewClient(ctx, tokenSource)

	a.service, err = sheets.NewService(ctx, option.WithHTTPClient(a.httpClient))
	if err != nil {
		return fmt.Errorf("failed to create Google Sheets service: %w", err)
	}

	return nil
}

func (a *OAuthActor) AuthToken(ctx context.Context, token *models2.OAuthToken, credential *models2.OAuthCredentials, scope Scope) error {
	return a.AuthTokenSource(ctx, token, credential, scope, &TokenSource{})
}

func (a *OAuthActor) AuthTokenSource(ctx context.Context, token *models2.OAuthToken, credential *models2.OAuthCredentials, scope Scope, tokenSource *TokenSource) error {
	if err := credential.Validate(); err != nil {
		return fmt.Errorf("invalid credential: %v", err)
	}
	if token == nil {
		return fmt.Errorf("invalid token")
	}

	credentialJSON, err := credential.ToJSON()
	if err != nil {
		return fmt.Errorf("error marshalling credential: %w", err)
	}

	config, err := google.ConfigFromJSON(credentialJSON, string(scope))
	if err != nil {
		return fmt.Errorf("failed to parse credential: %w", err)
	}

	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.config = config
	a.token = (*oauth2.Token)(token)

	if !a.token.Valid() {
		err = a.refreshToken(ctx, config, tokenSource)
		if err != nil {
			return fmt.Errorf("failed to refresh invalid access token: %w", err)
		}
	}

	baseSource := config.TokenSource(ctx, a.token)
	tokenSource.TokenSource = oauth2.ReuseTokenSource(a.token, baseSource)
	a.httpClient = oauth2.NewClient(ctx, tokenSource)

	a.service, err = sheets.NewService(ctx, option.WithHTTPClient(a.httpClient))
	if err != nil {
		return fmt.Errorf("failed to create Google Sheets service: %w", err)
	}

	return nil
}

func (a *OAuthActor) RefreshToken(ctx context.Context, credential *models2.OAuthCredentials, scope Scope) error {
	return a.RefreshTokenTokenSource(ctx, credential, scope, &TokenSource{})
}

func (a *OAuthActor) RefreshTokenTokenSource(ctx context.Context, credential *models2.OAuthCredentials, scope Scope, tokenSource *TokenSource) error {
	if err := credential.Validate(); err != nil {
		return fmt.Errorf("invalid credential: %v", err)
	}

	credentialJSON, err := credential.ToJSON()
	if err != nil {
		return fmt.Errorf("error marshalling credential: %w", err)
	}

	config, err := google.ConfigFromJSON(credentialJSON, string(scope))
	if err != nil {
		return fmt.Errorf("failed to parse credential: %w", err)
	}

	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.config = config

	return a.refreshToken(ctx, config, tokenSource)
}

func (a *OAuthActor) refreshToken(ctx context.Context, config *oauth2.Config, tokenSource *TokenSource) error {
	if a.token == nil {
		return fmt.Errorf("no token to refresh")
	}

	baseSource := config.TokenSource(ctx, a.token)
	tokenSource.TokenSource = oauth2.ReuseTokenSource(a.token, baseSource)

	newToken, err := tokenSource.Token()
	if err != nil {
		if strings.Contains(err.Error(), "invalid_grant") {
			return fmt.Errorf("refresh token invalid, re-authentication required")
		}
		return fmt.Errorf("failed to refresh access token: %w", err)
	}

	if newToken.RefreshToken == "" && a.token != nil {
		newToken.RefreshToken = a.token.RefreshToken
	}
	a.token = newToken

	return nil
}

func (a *OAuthActor) Service() *sheets.Service {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	return a.service
}

func (a *OAuthActor) Spreadsheets() *sheets.SpreadsheetsService {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	if a.service == nil {
		return nil
	}
	return a.service.Spreadsheets
}

func (a *OAuthActor) Token() *oauth2.Token {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	return a.token
}

func (a *OAuthActor) IsAuth() bool {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	return a.service != nil
}

func (a *OAuthActor) IsValidToken() bool {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	return a.token != nil && a.token.Valid()
}

func (a *OAuthActor) Close() error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	if a.httpClient != nil {
		if transport, ok := a.httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
	a.httpClient = nil
	a.service = nil
	a.token = nil
	a.config = nil
	return nil
}
