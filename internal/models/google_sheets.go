package models

import (
	"encoding/json"
	"encoding/xml"
	"sync"
	"time"
)

//// Main

type GoogleSheetsModel struct {
	OAuthCredentials   *GoogleSheetsOAuthCredentialsModel   `json:"oauth_credentials" xml:"oauth_credentials"`
	ServiceCredentials *GoogleSheetsServiceCredentialsModel `json:"service_credentials" xml:"service_credentials"`
	OAuthToken         *GoogleSheetsOAuthTokenModel         `json:"access_token" xml:"access_token"`
}

//// OAuth credentials

type GoogleSheetsOAuthCredentialsModel struct {
	mtx    sync.RWMutex
	hidden hiddenGoogleSheetsOAuthCredentialsModel
}

type hiddenGoogleSheetsOAuthCredentialsModel struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthUri                 string   `json:"auth_uri"`
	TokenUri                string   `json:"token_uri"`
	AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
}

func (m *GoogleSheetsOAuthCredentialsModel) ClientID() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ClientID
}

func (m *GoogleSheetsOAuthCredentialsModel) ProjectID() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ProjectID
}

func (m *GoogleSheetsOAuthCredentialsModel) AuthUri() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.AuthUri
}

func (m *GoogleSheetsOAuthCredentialsModel) TokenUri() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.TokenUri
}

func (m *GoogleSheetsOAuthCredentialsModel) AuthProviderX509CertUrl() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.AuthProviderX509CertUrl
}

func (m *GoogleSheetsOAuthCredentialsModel) ClientSecret() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ClientSecret
}

func (m *GoogleSheetsOAuthCredentialsModel) RedirectUris() []string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	copyUris := make([]string, len(m.hidden.RedirectUris))
	copy(copyUris, m.hidden.RedirectUris)
	return copyUris
}

func (m *GoogleSheetsOAuthCredentialsModel) SetAll(
	clientID, projectID, authUri, tokenUri, authProviderCertUrl, clientSecret string,
	redirectUris []string,
) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ClientID = clientID
	m.hidden.ProjectID = projectID
	m.hidden.AuthUri = authUri
	m.hidden.TokenUri = tokenUri
	m.hidden.AuthProviderX509CertUrl = authProviderCertUrl
	m.hidden.ClientSecret = clientSecret
	m.hidden.RedirectUris = append([]string{}, redirectUris...)
}

func (m *GoogleSheetsOAuthCredentialsModel) SetClientID(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ClientID = v
}

func (m *GoogleSheetsOAuthCredentialsModel) SetProjectID(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ProjectID = v
}

func (m *GoogleSheetsOAuthCredentialsModel) SetAuthUri(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.AuthUri = v
}

func (m *GoogleSheetsOAuthCredentialsModel) SetTokenUri(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.TokenUri = v
}

func (m *GoogleSheetsOAuthCredentialsModel) SetAuthProviderX509CertUrl(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.AuthProviderX509CertUrl = v
}

func (m *GoogleSheetsOAuthCredentialsModel) SetClientSecret(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ClientSecret = v
}

func (m *GoogleSheetsOAuthCredentialsModel) SetRedirectUris(uris []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	copyUris := make([]string, len(uris))
	copy(copyUris, uris)
	m.hidden.RedirectUris = copyUris
}

func (m *GoogleSheetsOAuthCredentialsModel) MarshalJSON() ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return json.Marshal(m.hidden)
}

func (m *GoogleSheetsOAuthCredentialsModel) UnmarshalJSON(data []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return json.Unmarshal(data, &m.hidden)
}

//// Service credentials

type GoogleSheetsServiceCredentialsModel struct {
	mtx    sync.RWMutex
	hidden hiddenGoogleSheetsServiceCredentialsModel
}

type hiddenGoogleSheetsServiceCredentialsModel struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

func (m *GoogleSheetsServiceCredentialsModel) Type() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Type
}

func (m *GoogleSheetsServiceCredentialsModel) ProjectID() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ProjectID
}

func (m *GoogleSheetsServiceCredentialsModel) PrivateKeyID() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.PrivateKeyID
}

func (m *GoogleSheetsServiceCredentialsModel) PrivateKey() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.PrivateKey
}

func (m *GoogleSheetsServiceCredentialsModel) ClientEmail() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ClientEmail
}

func (m *GoogleSheetsServiceCredentialsModel) ClientID() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ClientID
}

func (m *GoogleSheetsServiceCredentialsModel) AuthUri() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.AuthUri
}

func (m *GoogleSheetsServiceCredentialsModel) TokenUri() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.TokenUri
}

func (m *GoogleSheetsServiceCredentialsModel) AuthProviderX509CertUrl() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.AuthProviderX509CertUrl
}

func (m *GoogleSheetsServiceCredentialsModel) ClientX509CertUrl() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ClientX509CertUrl
}

func (m *GoogleSheetsServiceCredentialsModel) UniverseDomain() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.UniverseDomain
}

func (m *GoogleSheetsServiceCredentialsModel) SetAll(
	_type, projectID, privateKeyID, privateKey, clientEmail, clientID,
	authUri, tokenUri, authProviderCertUrl, clientCertUrl, universeDomain string,
) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Type = _type
	m.hidden.ProjectID = projectID
	m.hidden.PrivateKeyID = privateKeyID
	m.hidden.PrivateKey = privateKey
	m.hidden.ClientEmail = clientEmail
	m.hidden.ClientID = clientID
	m.hidden.AuthUri = authUri
	m.hidden.TokenUri = tokenUri
	m.hidden.AuthProviderX509CertUrl = authProviderCertUrl
	m.hidden.ClientX509CertUrl = clientCertUrl
	m.hidden.UniverseDomain = universeDomain
}

func (m *GoogleSheetsServiceCredentialsModel) SetType(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Type = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetProjectID(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ProjectID = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetPrivateKeyID(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.PrivateKeyID = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetPrivateKey(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.PrivateKey = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetClientEmail(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ClientEmail = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetClientID(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ClientID = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetAuthUri(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.AuthUri = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetTokenUri(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.TokenUri = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetAuthProviderX509CertUrl(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.AuthProviderX509CertUrl = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetClientX509CertUrl(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ClientX509CertUrl = v
}

func (m *GoogleSheetsServiceCredentialsModel) SetUniverseDomain(v string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.UniverseDomain = v
}

func (m *GoogleSheetsServiceCredentialsModel) MarshalJSON() ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return json.Marshal(m.hidden)
}

func (m *GoogleSheetsServiceCredentialsModel) UnmarshalJSON(data []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return json.Unmarshal(data, &m.hidden)
}

//// Access token

type GoogleSheetsOAuthTokenModel struct {
	mtx    sync.RWMutex
	hidden hiddenGoogleSheetsAccessTokenModel
}

type hiddenGoogleSheetsAccessTokenModel struct {
	AccessToken  string    `json:"access_token" xml:"access_token"`
	TokenType    string    `json:"token_type,omitempty" xml:"token_type"`
	RefreshToken string    `json:"refresh_token,omitempty" xml:"refresh_token"`
	Expiry       time.Time `json:"expiry,omitempty" xml:"expiry"`
	ExpiresIn    int64     `json:"expires_in,omitempty" xml:"expires_in"`
}

func (m *GoogleSheetsOAuthTokenModel) TokenType() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.TokenType
}
func (m *GoogleSheetsOAuthTokenModel) AccessToken() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.AccessToken
}
func (m *GoogleSheetsOAuthTokenModel) RefreshToken() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.RefreshToken
}
func (m *GoogleSheetsOAuthTokenModel) Expiry() time.Time {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Expiry
}
func (m *GoogleSheetsOAuthTokenModel) ExpiresIn() int64 {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ExpiresIn
}

func (m *GoogleSheetsOAuthTokenModel) SetAll(
	tokenType, accessToken, refreshToken string,
	expiry time.Time, expiresIn int64,
) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.AccessToken = accessToken
	m.hidden.TokenType = tokenType
	m.hidden.RefreshToken = refreshToken
	m.hidden.Expiry = expiry
	m.hidden.ExpiresIn = expiresIn
}

func (m *GoogleSheetsOAuthTokenModel) SetTokenType(t string) {
	m.mtx.Lock()
	m.hidden.TokenType = t
	m.mtx.Unlock()
}
func (m *GoogleSheetsOAuthTokenModel) SetAccessToken(token string) {
	m.mtx.Lock()
	m.hidden.AccessToken = token
	m.mtx.Unlock()
}
func (m *GoogleSheetsOAuthTokenModel) SetRefreshToken(token string) {
	m.mtx.Lock()
	m.hidden.RefreshToken = token
	m.mtx.Unlock()
}
func (m *GoogleSheetsOAuthTokenModel) SetExpiry(t time.Time) {
	m.mtx.Lock()
	m.hidden.Expiry = t
	m.mtx.Unlock()
}
func (m *GoogleSheetsOAuthTokenModel) SetExpiresIn(t int64) {
	m.mtx.Lock()
	m.hidden.ExpiresIn = t
	m.mtx.Unlock()
}

func (m *GoogleSheetsOAuthTokenModel) MarshalJSON() ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return json.Marshal(m.hidden)
}
func (m *GoogleSheetsOAuthTokenModel) UnmarshalJSON(data []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return json.Unmarshal(data, &m.hidden)
}
func (m *GoogleSheetsOAuthTokenModel) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return e.EncodeElement(m.hidden, start)
}
func (m *GoogleSheetsOAuthTokenModel) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return d.DecodeElement(&m.hidden, &start)
}
