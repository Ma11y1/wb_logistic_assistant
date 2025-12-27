package storage

import "wb_logistic_assistant/internal/models"

const (
	rawStorageMarker     byte = 0
	encryptStorageMarker byte = 1
)

type GoogleSheetsConfigStore interface {
	GetGoogleSheetsOAuthCredentials() *models.GoogleSheetsOAuthCredentialsModel
	GetGoogleSheetsServiceCredentials() *models.GoogleSheetsServiceCredentialsModel
	GetGoogleSheetsOAuthTokenModel() *models.GoogleSheetsOAuthTokenModel
	SetGoogleSheetsOAuthCredentials(credentials *models.GoogleSheetsOAuthCredentialsModel)
	SetGoogleSheetsServiceCredentials(credentials *models.GoogleSheetsServiceCredentialsModel)
	SetGoogleSheetsOAuthToken(token *models.GoogleSheetsOAuthTokenModel)
}

type WBLogisticConfigStore interface {
	GetWBLogisticLogin() string
	GetWBLogisticAccessToken() *models.WBLogisticAccessTokenModel
	GetWBLogisticSessionToken() *models.WBLogisticSessionTokenModel
	GetWBLogisticUserInfo() *models.WBLogisticUserInfoModel
	SetWBLogisticLogin(login string)
	SetWBLogisticAccessToken(token *models.WBLogisticAccessTokenModel)
	SetWBLogisticSessionToken(token *models.WBLogisticSessionTokenModel)
	SetWBLogisticUserInfo(userInfo *models.WBLogisticUserInfoModel)
}

type TelegramBotConfigStore interface {
	GetTelegramBotToken() string
	SetTelegramBotToken(token string)
}

type ConfigStore interface {
	GoogleSheetsConfigStore
	WBLogisticConfigStore
	TelegramBotConfigStore
	Set(name string, data string)
	SetBytes(name string, data []byte)
	Get(name string) string
	Remove(name string)
	Has(name string) bool
	Clear()
}

type CacheStore interface {
	Set(name, data string)
	Get(name string) string
	Remove(name string)
	Has(name string) bool
	Clear()
}

type Storage interface {
	ConfigStore() ConfigStore
	CacheStore() CacheStore
	Load(path string) error
	Save(path string) error
	Clear()
	SetEncrypt(password []byte)
	IsEncrypted() bool
}

type storageModel struct {
	ConfigStore ConfigStore `json:"config" bson:"config" xml:"config"  yaml:"config"`
	CacheStore  CacheStore  `json:"cache" bson:"cache" xml:"cache"  yaml:"cache"`
}
