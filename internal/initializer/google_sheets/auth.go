package google_sheets

import (
	"context"
	"time"
	"wb_logistic_assistant/external/google_sheets_api/auth"
	google_sheets_models "wb_logistic_assistant/external/google_sheets_api/models"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/utils"

	"golang.org/x/oauth2"
)

const (
	serverAddr      = "localhost:8080"
	callbackPath    = "/code"
	maxAuthAttempts = 3
)

func (i *Initializer) authOAuthSession() (*auth.OAuthActor, error) {
	logger.Log(logger.INFO, "Initializer.authOAuthSession()", "start auth google sheets oauth session")

	if !i.prompter.PromptGoogleSheetsQuestionAuthNewCredentials() {
		actor, err := i.authOAuthSessionStorage()
		if err == nil {
			return actor, nil
		}
		logger.Logf(logger.WARN, "Initializer.authOAuthSession()", "failed to auth google sheets oauth session using storage: %v", err)
		i.prompter.PromptGoogleSheetsAuthStorageFailed()
	} else {
		logger.Log(logger.INFO, "Initializer.authOAuthSession()", "auth google sheets oauth session with new credentials without using storage")
	}

	credentials, err := i.GetOAuthCredentialsFile()
	if err != nil {
		i.prompter.PromptGoogleSheetsReadCredentialsFailed()
		return nil, errors.Wrap(err, "Initializer.authOAuthSession()", "failed to get oauth credentials")
	}

	actor, err := i.authOAuthSessionAuto(credentials)
	if err != nil {
		logger.Logf(logger.WARN, "Initializer.authOAuthSession()", "failed to auth google sheets oauth session: %v", err)
		i.prompter.PromptGoogleSheetsAuthAutoFailed()

		actor, err = i.authOAuthSessionManual(credentials)
		if err != nil {
			return nil, errors.Wrap(err, "Initializer.authOAuthSession()", "failed to auth google sheets oauth session")
		}
	}

	if actor == nil || !actor.IsAuth() {
		return nil, errors.New("Initializer.authOAuthSession()", "oauth session is not auth")
	}

	err = i.SetOAuthCredentialsStorage(credentials)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSession()", "failed to set oauth credentials to storage")
	}

	err = i.SetOAuthTokenStorage(actor.Token())
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSession()", "failed to set oauth access token to storage")
	}

	logger.Log(logger.INFO, "Initializer.authOAuthSession()", "finish auth google sheets oauth session")
	return actor, nil
}

func (i *Initializer) authOAuthSessionStorage() (*auth.OAuthActor, error) {
	logger.Log(logger.INFO, "Initializer.authOAuthSessionStorage()", "start auth google sheets oauth session using storage")

	credentials, err := i.GetOAuthCredentialsStorage()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSessionStorage()", "failed to get oauth credentials")
	}

	token, err := i.GetOAuthTokenStorage()
	if err != nil {
		logger.Log(logger.WARN, "Initializer.authOAuthSessionStorage()", "failed to get google sheets oauth access token from storage")
	}

	actor := &auth.OAuthActor{}
	if token != nil {
		err = actor.AuthTokenSource(
			context.Background(),
			token,
			credentials,
			auth.SheetsAllScope,
			&auth.TokenSource{
				ErrorHandler:        i.errorRefreshOAuthTokenHandler,
				RefreshTokenHandler: i.refreshOAuthTokenHandler,
			},
		)
		if err != nil {
			logger.Logf(logger.WARN, "Initializer.authOAuthSessionStorage()", "failed to auth google sheets oauth session using access token from storage: %v", err)
		}
	}

	if !actor.IsAuth() {
		actor, err = i.authOAuthSessionAuto(credentials)
		if err != nil {
			logger.Logf(logger.WARN, "Initializer.authOAuthSessionStorage()", "failed to auth google sheets oauth session: %v", err)
			i.prompter.PromptGoogleSheetsAuthAutoFailed()

			actor, err = i.authOAuthSessionManual(credentials)
			if err != nil {
				return nil, errors.Wrap(err, "Initializer.authOAuthSessionStorage()", "failed to auth google sheets oauth session")
			}
		}
	}

	if !actor.IsAuth() {
		return nil, errors.New("Initializer.authOAuthSessionStorage()", "oauth session was not auth")
	}

	err = i.SetOAuthTokenStorage(actor.Token())
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSessionStorage()", "failed to set oauth access token to storage")
	}

	return actor, nil
}

func (i *Initializer) authOAuthSessionAuto(credentials *google_sheets_models.OAuthCredentials) (*auth.OAuthActor, error) {
	logger.Log(logger.INFO, "Initializer.authOAuthSessionAuto()", "start auto auth google sheets oauth session")
	actor := &auth.OAuthActor{}

	authorizer := auth.NewOAuthAuthorizer(serverAddr, callbackPath)
	defer authorizer.Stop()

	authorizerNotifications := auth.NewOAuthNotifications()
	authorizerNotifications.Successful = "Авторизация прошла успешно! Можно закрывать данную вкладку."
	err := authorizer.SetOAuthNotifications(authorizerNotifications)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSessionAuto()", "failed to set oauth authorizer notifications")
	}

	url, err := authorizer.Start(actor, credentials, auth.SheetsAllScope)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSessionAuto()", "failed to start oauth authorizer")
	}

	secondsWaitServer := i.config.GoogleSheets().Client().SecondsWaitServer()
	i.prompter.PromptGoogleSheetsRequestAuthCodeAuto(url, secondsWaitServer)
	utils.OpenBrowser(url)

	result, err := authorizer.Wait(time.Duration(secondsWaitServer) * time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSessionAuto()", "failed to authorize oauth session")
	}
	if result.Err != nil {
		return nil, errors.Wrap(result.Err, "Initializer.authOAuthSessionAuto()", "failed to authorize oauth session")
	}
	if result.Actor == nil {
		return nil, errors.New("Initializer.authOAuthSessionAuto()", "failed to authorize oauth session: no session found")
	}

	return result.Actor, nil
}

func (i *Initializer) authOAuthSessionManual(credentials *google_sheets_models.OAuthCredentials) (*auth.OAuthActor, error) {
	logger.Log(logger.INFO, "Initializer.authOAuthSessionManual()", "start manual auth google sheets oauth session")

	actor := &auth.OAuthActor{}

	url, err := actor.AuthCodeURL(credentials, auth.SheetsAllScope, auth.GenerateState())
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSessionManual()", "failed to get auth code url")
	}
	utils.OpenBrowser(url)

	code := ""
	for attempts := 0; attempts < maxAuthAttempts; attempts++ {
		code, err = i.prompter.PromptGoogleSheetsRequestAuthCode(url)
		if err == nil {
			break
		}
		logger.Logf(logger.WARN, "Initializer.authOAuthSessionManual()", "invalid code: '%s' entered: %v", code, err)
		i.prompter.PromptGoogleSheetsInvalidAuthCode()
	}

	err = actor.ExchangeCodeTokenSource(context.Background(), code, &auth.TokenSource{
		ErrorHandler:        i.errorRefreshOAuthTokenHandler,
		RefreshTokenHandler: i.refreshOAuthTokenHandler,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authOAuthSessionManual()", "failed to get auth code token source")
	}

	return actor, nil
}

func (i *Initializer) authServiceSession() (*auth.ServiceActor, error) {
	logger.Log(logger.INFO, "Initializer.authServiceSession()", "start auth service session Google Sheets")

	var credentials *google_sheets_models.ServiceCredentials
	var err error
	if i.prompter.PromptGoogleSheetsQuestionAuthNewCredentials() {
		logger.Log(logger.INFO, "Initializer.authServiceSession()", "auth google sheets service session with new credentials without using storage")

		credentials, err = i.GetServiceCredentialsFile()
		if err != nil {
			i.prompter.PromptGoogleSheetsReadCredentialsFailed()
			return nil, errors.Wrap(err, "Initializer.authServiceSession()", "failed to get service credentials from storage")
		}
	} else {
		credentials, err = i.GetServiceCredentialsStorage()
		if err != nil {
			logger.Logf(logger.WARN, "Initializer.authServiceSession()", "failed to get service credentials from storage: %v", err)
			i.prompter.PromptGoogleSheetsAuthStorageFailed()

			credentials, err = i.GetServiceCredentialsFile()
			if err != nil {
				i.prompter.PromptGoogleSheetsReadCredentialsFailed()
				return nil, errors.Wrap(err, "Initializer.authServiceSession()", "failed to get service credentials from storage")
			}
		}
	}

	service := &auth.ServiceActor{}
	err = service.Auth(context.Background(), credentials, auth.SheetsAllScope)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authServiceSession()", "failed to auth service actor")
	}

	if !service.IsAuth() || !service.IsValidToken() {
		return nil, errors.New("Initializer.authServiceSession()", "failed to auth service actor, it is not auth or valid token")
	}

	err = i.SetServiceCredentialsStorage(credentials)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authServiceSession()", "failed to save service credentials to storage")
	}

	logger.Log(logger.INFO, "Initializer.authServiceSession()", "finish auth service session Google Sheets")
	return service, nil
}

func (i *Initializer) errorRefreshOAuthTokenHandler(err error) {
	logger.Logf(logger.ERROR, "Initializer.errorRefreshOAuthTokenHandler()", "failed to refresh google sheets oauth token: %v", err)
}

func (i *Initializer) refreshOAuthTokenHandler(token *oauth2.Token) {
	if token == nil {
		logger.Log(logger.ERROR, "Initializer.refreshOAuthTokenHandler()", "invalid google sheets oauth token provided")
		return
	}
	err := i.SetOAuthTokenStorage(token)
	if err != nil {
		logger.Logf(logger.ERROR, "Initializer.refreshOAuthTokenHandler()", "failed to set google sheets oauth token to storage: %v", err)
		return
	}
	logger.Log(logger.INFO, "Initializer.refreshOAuthTokenHandler()", "google sheets oauth access token refreshed")
}
