package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
	"wb_logistic_assistant/external/google_sheets_api/models"
)

func GenerateState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

type OAuthNotifications struct {
	InvalidPath        string
	InvalidState       string
	InvalidCode        string
	FailedExchangeCode string
	Successful         string
}

func NewOAuthNotifications() *OAuthNotifications {
	return &OAuthNotifications{
		InvalidPath:        "Service redirected to wrong address! There may be internal error, you need to contact the support of the application from which auth was carried out.",
		InvalidState:       "Service returned invalid state! Body interception during auth is possible.",
		InvalidCode:        "Service returned invalid auth code! Try auth again.",
		FailedExchangeCode: "Failed to exchange auth code! Try auth in again.",
		Successful:         "Auth successful! You can close this page.",
	}
}

type OAuthAuthorizerResult struct {
	Actor *OAuthActor
	Err   error
}

type OAuthAuthorizer struct {
	mtx           sync.RWMutex
	server        *http.Server
	serverAddress string
	callbackPath  string
	chResult      chan *OAuthAuthorizerResult
	state         string
	isStarted     bool
	actor         *OAuthActor
	credentials   *models.OAuthCredentials
	scope         Scope
	notifications *OAuthNotifications
}

func NewOAuthAuthorizer(serverAddress, callbackPath string) *OAuthAuthorizer {
	if serverAddress == "" {
		serverAddress = "localhost:8080"
	}
	if callbackPath == "" {
		callbackPath = "/callback"
	} else if callbackPath[0] != '/' {
		callbackPath = "/" + callbackPath
	}

	authorizer := &OAuthAuthorizer{
		serverAddress: serverAddress,
		callbackPath:  callbackPath,
		notifications: NewOAuthNotifications(),
	}
	authorizer.server = &http.Server{
		Addr:    serverAddress,
		Handler: http.HandlerFunc(authorizer.handler),
	}
	return authorizer
}

func (c *OAuthAuthorizer) Start(actor *OAuthActor, credential *models.OAuthCredentials, scope Scope) (string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.isStarted {
		return "", fmt.Errorf("code catcher already started")
	}

	c.actor = actor
	c.credentials = credential
	c.scope = scope
	c.chResult = make(chan *OAuthAuthorizerResult, 2)
	c.isStarted = true

	redirectURI := fmt.Sprintf("http://%s%s", c.serverAddress, c.callbackPath)
	copyCredential := *credential
	copyCredential.Installed.RedirectUris = append([]string{redirectURI}, credential.Installed.RedirectUris...)

	c.state = GenerateState()

	// when creating a configuration, the redirect URI is taken from credential.installed.redirectURIs[0]
	authURL, err := c.actor.AuthCodeURL(&copyCredential, scope, c.state)
	if err != nil {
		return "", fmt.Errorf("failed to generate session URL: %w", err)
	}

	go c.startServer()

	return authURL, nil
}

func (c *OAuthAuthorizer) Wait(timeout time.Duration) (*OAuthAuthorizerResult, error) {
	select {
	case result := <-c.chResult:
		return result, nil
	case <-time.After(timeout):
		_ = c.Stop()
		return nil, fmt.Errorf("authorization timeout after %s", timeout)
	}
}

func (c *OAuthAuthorizer) startServer() {
	if err := c.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		c.pushResult(nil, fmt.Errorf("oauth catcher server error: %w", err))
	}
}

func (c *OAuthAuthorizer) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != c.callbackPath {
		http.Error(w, c.notifications.InvalidPath, http.StatusNotFound)
		c.pushResult(nil, fmt.Errorf("invalid callback path: current %s, needs %s", r.URL.Path, c.callbackPath))
		return
	}

	state := r.URL.Query().Get("state")
	if state != c.state {
		http.Error(w, c.notifications.InvalidState, http.StatusBadRequest)
		c.pushResult(nil, fmt.Errorf("invalid state: expected %s, got %s", c.state, state))
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, c.notifications.InvalidCode, http.StatusBadRequest)
		c.pushResult(nil, fmt.Errorf("missing code in request"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := c.actor.ExchangeCode(ctx, code)
	if err != nil {
		http.Error(w, c.notifications.FailedExchangeCode, http.StatusInternalServerError)
		c.pushResult(nil, fmt.Errorf("failed to exchange code: %w", err))
		return
	}

	w.Write([]byte(c.notifications.Successful))
	c.pushResult(c.actor, nil)
}

func (c *OAuthAuthorizer) pushResult(actor *OAuthActor, err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.isStarted {
		select {
		case c.chResult <- &OAuthAuthorizerResult{Actor: actor, Err: err}:
		default:
		}
	}
}

func (c *OAuthAuthorizer) Stop() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if !c.isStarted {
		return nil
	}
	c.isStarted = false
	close(c.chResult)
	c.chResult = nil
	return c.server.Close()
}

func (c *OAuthAuthorizer) SetOAuthNotifications(notifications *OAuthNotifications) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.isStarted {
		return errors.New("oauth authorizer already started")
	}
	c.notifications = notifications
	return nil
}
