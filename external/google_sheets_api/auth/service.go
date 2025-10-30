package auth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"sync"
	"wb_logistic_assistant/external/google_sheets_api/models"
)

type ServiceActor struct {
	mtx        sync.RWMutex
	httpClient *http.Client
	service    *sheets.Service
}

func (*ServiceActor) Type() ActorType {
	return ServiceActorType
}

func (a *ServiceActor) Auth(ctx context.Context, credential *models.ServiceCredentials, scope Scope) error {
	if err := credential.Validate(); err != nil {
		return fmt.Errorf("invalid credential: %w", err)
	}

	credentialJSON, err := credential.ToJSON()
	if err != nil {
		return fmt.Errorf("error marshalling credential: %w", err)
	}

	config, err := google.JWTConfigFromJSON(credentialJSON, string(scope))
	if err != nil {
		return fmt.Errorf("error create config: %w", err)
	}

	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.httpClient = config.Client(ctx)
	a.service, err = sheets.NewService(ctx, option.WithHTTPClient(a.httpClient))
	if err != nil {
		return fmt.Errorf("failed to create Google Sheets service: %w", err)
	}

	return nil
}

func (a *ServiceActor) Service() *sheets.Service {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	return a.service
}

func (a *ServiceActor) Spreadsheets() *sheets.SpreadsheetsService {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	if a.service == nil {
		return nil
	}
	return a.service.Spreadsheets
}

// Token returns nil since ServiceAccount does not use oauth2.Token.
func (a *ServiceActor) Token() *oauth2.Token {
	return nil
}

func (a *ServiceActor) IsAuth() bool {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	return a.service != nil
}

func (a *ServiceActor) IsValidToken() bool {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	return a.service != nil
}

func (a *ServiceActor) Close() error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	if a.httpClient != nil {
		if transport, ok := a.httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
	a.httpClient = nil
	a.service = nil
	return nil
}
