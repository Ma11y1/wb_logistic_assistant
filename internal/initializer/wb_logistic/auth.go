package wb_logistic

import (
	"context"
	"encoding/json"
	"time"
	"wb_logistic_assistant/external/wb_logistic_api"
	wb_models "wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/external/wb_logistic_api/session"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
)

const (
	RepeatCodeRequest = 1
	ManualSessionAuth = 2
	Exit              = 3
)

func (i *Initializer) AuthSession(client *wb_logistic_api.Client) (*session.Session, error) {
	logger.Log(logger.INFO, "Initializer.AuthSession()", "start auth wb logistic session")

	login, err := i.getLogin()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.AuthSession()", "")
	}

	var accessToken *wb_models.AuthAccessToken
	var timeWaiting time.Time
mainLoop:
	for {
		resCodeData, err := client.RequestAuthCode(context.Background(), login)
		if err != nil {
			return nil, errors.Wrap(err, "Initializer.AuthSession()", "failed to request auth code")
		}

		code := 0
		timeWaiting = time.Now().Add(time.Duration(resCodeData.Ttl) * time.Second)
		for {
			timeRemains := int(timeWaiting.Sub(time.Now()).Seconds())
			if timeRemains < 0 {
				timeRemains = 0
			}
			code = i.prompter.PromptWBLogisticRequestAuthCode(resCodeData.AuthMethod, timeRemains)

			switch code {
			case RepeatCodeRequest:
				{
					// if the TTL time has passed, a new code is requested
					if time.Now().After(timeWaiting) {
						continue mainLoop
					}
					continue
				}
			case ManualSessionAuth:
				{
					accessToken, err = i.getManualAccessToken()
					if err != nil {
						i.prompter.PromptWBLogisticInvalidAccessTokenData()
						logger.Logf(logger.ERROR, "Initializer.AuthSession()", "failed to manual get access token: %v", err)
						continue
					}
					break mainLoop
				}
			case Exit:
				return nil, errors.New("Initializer.AuthSession()", "failed request auth code, action is canceled")
			default:
				accessToken, err = client.ExchangeAuthCode(context.Background(), code, resCodeData.Sticker)
				if err != nil {
					logger.Logf(logger.ERROR, "Initializer.AuthSession()", "failed to exchange auth code: %v", err)
					i.prompter.PromptWBLogisticAuthFailed()
					continue
				}
				break mainLoop
			}
		}
	}

	sessionToken, userInfo, err := client.GetSessionToken(context.Background(), login, accessToken.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.AuthSession()", "failed to get session token")
	}

	if err = userInfo.Validate(); err != nil {
		return nil, errors.Wrap(err, "Initializer.AuthSession()", "invalid user info. account does not exist or does not have access to this service")
	}

	s := session.NewSessionFromToken(login, accessToken, sessionToken, userInfo)
	if !s.IsAuth() || s.SessionTokenExpired() || s.AccessTokenExpired() {
		return nil, errors.New("Initializer.AuthSession()", "failed to auth session: session was not authorized, invalid data received")
	}

	s.Emitter().On(session.EventUpdateAccessToken, i.updateAccessTokenHandler)
	s.Emitter().On(session.EventUpdateSessionToken, i.updateSessionTokenHandler)
	s.Emitter().On(session.EventUpdateUserInfo, i.updateUserInfoHandler)

	i.SetLoginStorage(login)
	err = i.SetAccessTokenStorage(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.AuthSession()", "failed to save access token to storage")
	}
	err = i.SetSessionTokenStorage(sessionToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.AuthSession()", "failed to save session token to storage")
	}

	err = i.SetUserInfoStorage(userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.AuthSession()", "failed to save user info to storage")
	}

	err = i.UpdateStorage()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.AuthSession()", "")
	}

	logger.Logf(logger.INFO, "Initializer.AuthSession()", "successful initialize WB logistic client session id:'%d' login:'%s'", s.UserInfo().ID, login)
	return s, nil
}

func (i *Initializer) getLogin() (string, error) {
	for attempt := 0; attempt < 3; attempt++ {
		login := i.prompter.PromptWBLogisticRequestAuthLogin()
		if len(login) > 10 {
			if login[0] == '+' {
				login = login[1:]
			}
			return login, nil
		}

		logger.Logf(logger.WARN, "Initializer.getLogin()", "invalid login: '%s' entered", login)
		i.prompter.PromptWBLogisticInvalidAuthLogin()

		if attempt >= 2 {
			return "", errors.New("Initializer.getLogin()", "failed request auth login")
		}
	}
	return "", nil
}

func (i *Initializer) getManualAccessToken() (*wb_models.AuthAccessToken, error) {
	token := i.prompter.PromptWBLogisticRequestAccessTokenData()
	if token == "" {
		return nil, errors.New("Initializer.getManualAccessToken()", "token is empty")
	}

	accessToken := &wb_models.AuthAccessToken{}
	err := json.Unmarshal([]byte(token), accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.getManualAccessToken()", "invalid format data access token")
	}

	if accessToken.AccessToken == "" || accessToken.RefreshToken == "" || accessToken.ExpiresIn < 0 {
		return nil, errors.New("Initializer.getManualAccessToken()", "invalid data access token")
	}

	return accessToken, nil
}

func (i *Initializer) authSessionStorage(client *wb_logistic_api.Client) (*session.Session, error) {
	logger.Log(logger.INFO, "Initializer.authSessionStorage()", "start auth wb logistic session using storage")

	login := i.GetLoginStorage()
	if login == "" {
		return nil, errors.New("Initializer.authSessionStorage()", "no login storage found")
	}
	accessToken, err := i.GetAccessTokenStorage()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSessionStorage()", "failed get access token from storage")
	}

	sessionToken, err := i.GetSessionTokenStorage()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSessionStorage()", "failed get session token from storage")
	}
	userInfo, err := i.GetUserInfoStorage()
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSessionStorage()", "failed get user info from storage")
	}

	s := session.NewSessionFromToken(login, accessToken, sessionToken, userInfo)
	if !s.IsAuth() {
		return nil, errors.New("Initializer.authSessionStorage()", "session was not auth")
	}
	if s.AccessTokenExpired() {
		return nil, errors.New("Initializer.authSessionStorage()", "access token expired")
	}
	if s.SessionTokenExpired() {
		sessionToken, userInfo, err = client.GetSessionToken(context.Background(), login, accessToken.AccessToken)
		if err != nil {
			return nil, errors.Wrap(err, "Initializer.authSessionStorage()", "failed refresh session token")
		}
		if err = userInfo.Validate(); err != nil {
			return nil, errors.Wrap(err, "Initializer.authSessionStorage()", "invalid user info. account does not exist or does not have access to this service")
		}
	}

	s.Emitter().On(session.EventUpdateAccessToken, i.updateAccessTokenHandler)
	s.Emitter().On(session.EventUpdateSessionToken, i.updateSessionTokenHandler)
	s.Emitter().On(session.EventUpdateUserInfo, i.updateUserInfoHandler)

	err = i.SetAccessTokenStorage(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSessionStorage()", "failed to save access token to storage")
	}
	err = i.SetSessionTokenStorage(sessionToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSessionStorage()", "failed to save session token to storage")
	}
	err = i.SetUserInfoStorage(userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSessionStorage()", "failed to save user info to storage")
	}

	logger.Logf(logger.INFO, "Initializer.authSessionStorage()", "finish auth wb logistic session id:'%d' phone number:'%s' using storage", s.UserInfo().ID, s.UserInfo().UserDetails.PhoneNumber)
	return s, nil
}

func (i *Initializer) updateAccessTokenHandler(s *session.Session) {
	err := i.SetAccessTokenStorage(s.AccessToken())
	if err != nil {
		logger.Log(logger.ERROR, "Initializer.updateAccessTokenHandler()", "failed to save access token to storage")
	}
}

func (i *Initializer) updateSessionTokenHandler(s *session.Session) {
	err := i.SetSessionTokenStorage(s.SessionToken())
	if err != nil {
		logger.Log(logger.ERROR, "Initializer.updateSessionTokenHandler()", "failed to save session token to storage")
	}
}

func (i *Initializer) updateUserInfoHandler(s *session.Session) {
	err := i.SetUserInfoStorage(s.UserInfo())
	if err != nil {
		logger.Log(logger.ERROR, "Initializer.updateUserInfoHandler()", "failed to save user info to storage")
	}
}
