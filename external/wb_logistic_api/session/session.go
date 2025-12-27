package session

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

type Session struct {
	mtx          sync.RWMutex
	login        string
	accessToken  *models.AuthAccessToken
	sessionToken *models.AuthSessionToken
	userInfo     *models.UserInfo
	emitter      *Emitter
}

type session struct {
	Login       string                   `json:"login"`
	AccessToken *models.AuthAccessToken  `json:"access_token"`
	Session     *models.AuthSessionToken `json:"session_token"`
	UserInfo    *models.UserInfo         `json:"user_info"`
}

func NewSession() *Session {
	return &Session{emitter: &Emitter{}}
}

func NewSessionFromToken(
	login string,
	accessToken *models.AuthAccessToken,
	sessionToken *models.AuthSessionToken,
	userInfo *models.UserInfo,
) *Session {
	return &Session{
		login:        login,
		accessToken:  accessToken,
		sessionToken: sessionToken,
		userInfo:     userInfo,
		emitter:      NewEmitter(),
	}
}

func (s *Session) Emitter() *Emitter {
	return s.emitter
}

func (s *Session) Set(
	login string,
	accessToken *models.AuthAccessToken,
	sessionToken *models.AuthSessionToken,
	userInfo *models.UserInfo,
) error {
	if login == "" {
		return fmt.Errorf("empty login")
	}
	if accessToken == nil || sessionToken == nil {
		return fmt.Errorf("access or session token is nil")
	}
	if userInfo == nil {
		return fmt.Errorf("user info is nil")
	}
	s.mtx.Lock()
	s.login = login
	s.accessToken = accessToken
	s.sessionToken = sessionToken
	s.userInfo = userInfo
	s.mtx.Unlock()
	s.emitter.Emit(EventUpdateAccessToken, s)
	s.emitter.Emit(EventUpdateSessionToken, s)
	s.emitter.Emit(EventUpdateUserInfo, s)
	return nil
}

func (s *Session) SetAccessToken(accessToken *models.AuthAccessToken) {
	if accessToken == nil {
		return
	}
	s.mtx.Lock()
	s.accessToken = accessToken
	s.mtx.Unlock()
	s.emitter.Emit(EventUpdateAccessToken, s)
}

func (s *Session) SetSessionToken(sessionToken *models.AuthSessionToken) {
	if sessionToken == nil {
		return
	}
	s.mtx.Lock()
	s.sessionToken = sessionToken
	s.mtx.Unlock()
	s.emitter.Emit(EventUpdateSessionToken, s)
}

func (s *Session) SetUserInfo(userInfo *models.UserInfo) {
	if userInfo == nil {
		return
	}
	s.mtx.Lock()
	s.userInfo = userInfo
	s.mtx.Unlock()
	s.emitter.Emit(EventUpdateUserInfo, s)
}

func (s *Session) Login() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.login
}

// AccessToken Used to obtain session token. In most cases, it is issued indefinitely and is used to obtain a session token for various WB services
func (s *Session) AccessToken() *models.AuthAccessToken {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.accessToken
}

// SessionToken Obtained using an access token and used only for accessing the WB logistics API
func (s *Session) SessionToken() *models.AuthSessionToken {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.sessionToken
}

func (s *Session) UserInfo() *models.UserInfo {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.userInfo
}

// RefreshToken Access token refresh token, not session token
func (s *Session) RefreshToken() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if s.accessToken == nil {
		return ""
	}
	return s.accessToken.RefreshToken
}

func (s *Session) AccessTokenString() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if s.accessToken == nil {
		return ""
	}
	return s.accessToken.AccessToken
}

func (s *Session) SessionTokenString() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if s.sessionToken == nil {
		return ""
	}
	return s.sessionToken.Token.AccessToken
}

func (s *Session) AccessTokenExpiresIn() int64 {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if s.accessToken == nil {
		return 0
	}
	return s.accessToken.ExpiresIn
}

func (s *Session) SessionTokenExpiresIn() int64 {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if s.sessionToken == nil {
		return 0
	}
	return s.sessionToken.Token.ExpiresIn
}

func (s *Session) AccessTokenExpired() bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if s.accessToken == nil {
		return true
	}
	// Value returned by the server may be 0, but the token will be valid
	return s.accessToken.ExpiresIn != 0 && time.Now().Unix() >= s.accessToken.ExpiresIn
}

func (s *Session) SessionTokenExpired() bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	if s.sessionToken == nil {
		return true
	}
	return time.Now().Unix() >= s.sessionToken.Token.ExpiresIn
}

func (s *Session) Clear() {
	s.mtx.Lock()
	s.accessToken = nil
	s.sessionToken = nil
	s.userInfo = nil
	s.mtx.Unlock()
	s.emitter.Emit(EventClear, s)
}

func (s *Session) IsAuth() bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.accessToken != nil &&
		s.accessToken.AccessToken != "" &&
		s.sessionToken != nil &&
		s.userInfo != nil &&
		s.sessionToken.Token.AccessToken != ""
}

func (s *Session) MarshalJSON() ([]byte, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return json.Marshal(session{
		AccessToken: s.accessToken,
		Session:     s.sessionToken,
		UserInfo:    s.userInfo,
	})
}

func (s *Session) UnmarshalJSON(data []byte) error {
	temp := session{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	s.mtx.Lock()
	s.accessToken = temp.AccessToken
	s.sessionToken = temp.Session
	s.userInfo = temp.UserInfo
	s.mtx.Unlock()
	return nil
}
