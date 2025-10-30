package storage

import (
	"encoding/json"
	"sync"
	"wb_logistic_assistant/internal/models"
)

type FileConfigStore struct {
	mtx          sync.RWMutex
	googleSheets *models.GoogleSheetsModel
	wbLogistic   *models.WBLogisticModel
	telegramBot  *models.TelegramBot
	data         map[string][]byte
}

type fileConfigStore struct {
	GoogleSheets *models.GoogleSheetsModel `json:"google_sheets" xml:"google_sheets"`
	WBLogistic   *models.WBLogisticModel   `json:"wb_logistic" xml:"wb_logistic"`
	TelegramBot  *models.TelegramBot       `json:"telegram_bot" xml:"telegram_bot"`
	Data         map[string][]byte         `json:"data" xml:"data"`
}

func NewFileConfigStore() *FileConfigStore {
	return &FileConfigStore{
		googleSheets: &models.GoogleSheetsModel{},
		wbLogistic:   &models.WBLogisticModel{},
		telegramBot:  &models.TelegramBot{},
		data:         make(map[string][]byte),
	}
}

//// Google sheets

func (c *FileConfigStore) GetGoogleSheetsOAuthCredentials() *models.GoogleSheetsOAuthCredentialsModel {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.googleSheets.OAuthCredentials
}
func (c *FileConfigStore) GetGoogleSheetsServiceCredentials() *models.GoogleSheetsServiceCredentialsModel {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.googleSheets.ServiceCredentials
}
func (c *FileConfigStore) GetGoogleSheetsOAuthTokenModel() *models.GoogleSheetsOAuthTokenModel {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.googleSheets.OAuthToken
}

func (c *FileConfigStore) SetGoogleSheetsOAuthCredentials(credentials *models.GoogleSheetsOAuthCredentialsModel) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.googleSheets.OAuthCredentials = credentials
}
func (c *FileConfigStore) SetGoogleSheetsServiceCredentials(credentials *models.GoogleSheetsServiceCredentialsModel) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.googleSheets.ServiceCredentials = credentials
}
func (c *FileConfigStore) SetGoogleSheetsOAuthToken(token *models.GoogleSheetsOAuthTokenModel) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.googleSheets.OAuthToken = token
}

//// WB logistic

func (c *FileConfigStore) GetWBLogisticLogin() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.wbLogistic.Login
}
func (c *FileConfigStore) GetWBLogisticAccessToken() *models.WBLogisticAccessTokenModel {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.wbLogistic.AccessToken
}
func (c *FileConfigStore) GetWBLogisticSessionToken() *models.WBLogisticSessionTokenModel {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.wbLogistic.MergedAccessToken
}
func (c *FileConfigStore) GetWBLogisticUserInfo() *models.WBLogisticUserInfoModel {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.wbLogistic.UserInfo
}

func (c *FileConfigStore) SetWBLogisticLogin(login string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.wbLogistic.Login = login
}
func (c *FileConfigStore) SetWBLogisticAccessToken(token *models.WBLogisticAccessTokenModel) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.wbLogistic.AccessToken = token
}
func (c *FileConfigStore) SetWBLogisticSessionToken(token *models.WBLogisticSessionTokenModel) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.wbLogistic.MergedAccessToken = token
}
func (c *FileConfigStore) SetWBLogisticUserInfo(userInfo *models.WBLogisticUserInfoModel) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.wbLogistic.UserInfo = userInfo
}

//// Telegram bot

func (c *FileConfigStore) GetTelegramBotToken() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.telegramBot.Token
}

func (c *FileConfigStore) SetTelegramBotToken(token string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.telegramBot.Token = token
}

//// Map

func (c *FileConfigStore) Set(name string, data string) {
	if name == "" {
		return
	}
	c.mtx.Lock()
	c.data[name] = []byte(data)
	c.mtx.Unlock()
}

func (c *FileConfigStore) SetBytes(name string, data []byte) {
	if name == "" || data == nil {
		return
	}
	c.mtx.Lock()
	copied := make([]byte, len(data))
	copy(copied, data)
	c.data[name] = copied
	c.mtx.Unlock()
}

func (c *FileConfigStore) Get(name string) string {
	if name == "" {
		return ""
	}
	c.mtx.RLock()
	data := c.data[name]
	c.mtx.RUnlock()
	return string(data)
}

func (c *FileConfigStore) Remove(name string) {
	if name == "" {
		return
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if val, ok := c.data[name]; ok {
		for i := range val {
			val[i] = 0
		}
		delete(c.data, name)
	}
}

func (c *FileConfigStore) Has(name string) bool {
	if name == "" {
		return false
	}
	c.mtx.RLock()
	_, ok := c.data[name]
	c.mtx.RUnlock()
	return ok
}

func (c *FileConfigStore) Clear() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for k, v := range c.data {
		for i := range v {
			v[i] = 0
		}
		delete(c.data, k)
	}
	c.googleSheets = &models.GoogleSheetsModel{}
	c.wbLogistic = &models.WBLogisticModel{}
	c.telegramBot = &models.TelegramBot{}
}

func (c *FileConfigStore) MarshalJSON() ([]byte, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return json.Marshal(&fileConfigStore{
		GoogleSheets: c.googleSheets,
		WBLogistic:   c.wbLogistic,
		TelegramBot:  c.telegramBot,
		Data:         c.data,
	})
}

func (c *FileConfigStore) UnmarshalJSON(data []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	temp := &fileConfigStore{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	c.googleSheets = temp.GoogleSheets
	c.wbLogistic = temp.WBLogistic
	c.telegramBot = temp.TelegramBot
	c.data = temp.Data
	return nil
}
