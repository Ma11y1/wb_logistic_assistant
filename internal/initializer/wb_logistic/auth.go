package wb_logistic

import (
	"context"
	"wb_logistic_assistant/external/wb_logistic_api"
	"wb_logistic_assistant/external/wb_logistic_api/session"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
)

const attemptsQuestions = 3

func (i *Initializer) authSession(client *wb_logistic_api.Client) (*session.Session, error) {
	logger.Log(logger.INFO, "Initializer.authSession()", "start auth wb logistic session")

	var login string
	var err error
	for attempts := 0; attempts < attemptsQuestions; attempts++ {
		login, err = i.prompter.PromptWBLogisticRequestAuthLogin()
		if err == nil {
			break
		}
		logger.Logf(logger.WARN, "Initializer.authSession()", "invalid login: '%s' entered: %v", login, err)
	}

	authCodeData, err := client.RequestAuthCode(context.Background(), login)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSession()", "failed to request auth code")
	}

	var code int
	for attempts := 0; attempts < attemptsQuestions; attempts++ {
		code, err = i.prompter.PromptWBLogisticRequestAuthCode(authCodeData.AuthMethod, authCodeData.Ttl)
		if err == nil {
			break
		}
		logger.Logf(logger.WARN, "Initializer.authSession()", "invalid code: '%d' entered: %v", code, err)
		i.prompter.PromptWBLogisticInvalidAuthCode()
	}

	accessToken, err := client.ExchangeAuthCode(context.Background(), code, authCodeData.Sticker)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSession()", "failed to exchange auth code")
	}

	sessionToken, userInfo, err := client.GetSessionToken(context.Background(), login, accessToken.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSession()", "failed to get session token")
	}

	s := session.NewSessionFromToken(login, accessToken, sessionToken, userInfo)
	if !s.IsAuth() || s.SessionTokenExpired() || s.AccessTokenExpired() {
		return nil, errors.New("Initializer.authSession()", "failed to auth session: session was not authorized, invalid data received")
	}

	s.Emitter().On(session.EventUpdateAccessToken, i.updateAccessTokenHandler)
	s.Emitter().On(session.EventUpdateSessionToken, i.updateSessionTokenHandler)
	s.Emitter().On(session.EventUpdateUserInfo, i.updateUserInfoHandler)

	i.SetLoginStorage(login)
	err = i.SetAccessTokenStorage(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSession()", "failed to save access token to storage")
	}
	err = i.SetSessionTokenStorage(sessionToken)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSession()", "failed to save session token to storage")
	}

	if err = userInfo.Validate(); err != nil {
		return nil, errors.Wrap(err, "Initializer.authSession()", "invalid user info. account does not exist or does not have access to this service")
	}

	err = i.SetUserInfoStorage(userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.authSession()", "failed to save user info to storage")
	}

	logger.Logf(logger.INFO, "Initializer.authSession()", "successful initialize WB logistic client session id:'%d' login:'%s'", s.UserInfo().ID, login)
	return s, nil
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
