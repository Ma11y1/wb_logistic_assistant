package models

import (
	"encoding/json"
	"sync"
	"time"
)

//// Main

type OLLogisticModel struct {
	AccessToken *OLLogisticAccessTokenModel `json:"access_token"`
	AuthData    *OLLogisticAuthDataModel    `json:"login_data"`
	UserInfo    *OLLogisticUserInfoModel    `json:"user_info"`
}

//// Login data

type OLLogisticAuthDataModel struct {
	mtx    sync.RWMutex
	hidden hiddenOLLogisticAuthDataModel
}

type hiddenOLLogisticAuthDataModel struct {
	Login    string `json:"login"`
	Password []byte `json:"password"`
}

func (m *OLLogisticAuthDataModel) SetLogin(login string) {
	m.mtx.Lock()
	m.hidden.Login = login
	m.mtx.Unlock()
}
func (m *OLLogisticAuthDataModel) SetPassword(password []byte) {
	m.mtx.Lock()
	for i := 0; i < len(m.hidden.Password); i++ {
		m.hidden.Password[i] = '0'
	}
	m.hidden.Password = password
	m.mtx.Unlock()
}

func (m *OLLogisticAuthDataModel) GetLogin() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Login
}
func (m *OLLogisticAuthDataModel) GetPassword() []byte {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Password
}

func (m *OLLogisticAuthDataModel) Clear() {
	m.hidden.Login = ""
	for i := 0; i < len(m.hidden.Password); i++ {
		m.hidden.Password[i] = '0'
	}
}

func (m *OLLogisticAuthDataModel) MarshalJSON() ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return json.Marshal(m.hidden)
}
func (m *OLLogisticAuthDataModel) UnmarshalJSON(data []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return json.Unmarshal(data, &m.hidden)
}

//// Access token

type OLLogisticAccessTokenModel struct {
	mtx    sync.RWMutex
	hidden hiddenOLLogisticAccessTokenModel
}

type hiddenOLLogisticAccessTokenModel struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	TsExpiresIn  int64  `json:"ts_expires_in"`
}

func (m *OLLogisticAccessTokenModel) GetTokenType() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.TokenType
}
func (m *OLLogisticAccessTokenModel) GetAccessToken() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.AccessToken
}
func (m *OLLogisticAccessTokenModel) GetRefreshToken() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.RefreshToken
}
func (m *OLLogisticAccessTokenModel) GetScope() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Scope
}
func (m *OLLogisticAccessTokenModel) GetExpiresIn() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ExpiresIn
}
func (m *OLLogisticAccessTokenModel) GetTsExpiresIn() int64 {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.TsExpiresIn
}

func (m *OLLogisticAccessTokenModel) SetAll(tokenType, accessToken, refreshToken, scope string, expiresIn int, tsExpiresIn int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.hidden.TokenType = tokenType
	m.hidden.AccessToken = accessToken
	m.hidden.RefreshToken = refreshToken
	m.hidden.Scope = scope
	m.hidden.ExpiresIn = expiresIn
	m.hidden.TsExpiresIn = tsExpiresIn
}

func (m *OLLogisticAccessTokenModel) SetTokenType(tokenType string) {
	m.mtx.Lock()
	m.hidden.TokenType = tokenType
	m.mtx.Unlock()
}
func (m *OLLogisticAccessTokenModel) SetAccessToken(token string) {
	m.mtx.Lock()
	m.hidden.AccessToken = token
	m.mtx.Unlock()
}
func (m *OLLogisticAccessTokenModel) SetRefreshToken(token string) {
	m.mtx.Lock()
	m.hidden.RefreshToken = token
	m.mtx.Unlock()
}
func (m *OLLogisticAccessTokenModel) SetScope(scope string) {
	m.mtx.Lock()
	m.hidden.Scope = scope
	m.mtx.Unlock()
}
func (m *OLLogisticAccessTokenModel) SetExpiresIn(seconds int) {
	m.mtx.Lock()
	m.hidden.ExpiresIn = seconds
	m.mtx.Unlock()
}
func (m *OLLogisticAccessTokenModel) SetTsExpiresIn(seconds int64) {
	m.mtx.Lock()
	m.hidden.TsExpiresIn = seconds
	m.mtx.Unlock()
}

func (m *OLLogisticAccessTokenModel) MarshalJSON() ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return json.Marshal(m.hidden)
}

func (m *OLLogisticAccessTokenModel) UnmarshalJSON(data []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return json.Unmarshal(data, &m.hidden)
}

//// User info

type OLLogisticUserInfoModel struct {
	mtx    sync.RWMutex
	hidden hiddenOLLogisticUserInfoModel
}

type hiddenOLLogisticUserInfoModel struct {
	ID                    int                       `json:"id"`
	Name                  string                    `json:"name"`
	Login                 string                    `json:"login"`
	PhoneNumber           string                    `json:"phoneNumber"`
	FreelancerEmployeeID  interface{}               `json:"freelancerEmployeeId"`
	SupplierID            int                       `json:"supplierId"`
	SupplierName          string                    `json:"supplierName"`
	RoleID                int                       `json:"roleId"`
	IsBoss                bool                      `json:"isBoss"`
	IsDriver              bool                      `json:"isDriver"`
	IsDeleted             bool                      `json:"isDeleted"`
	IsVerified            bool                      `json:"isVerified"`
	IsNewFreelancer       bool                      `json:"isNewFreelancer"`
	SupplierVatID         interface{}               `json:"supplierVatId"`
	SupplierBankAccountID interface{}               `json:"supplierBankAccountId"`
	SupplierContractDt    interface{}               `json:"supplierContractDt"`
	Roles                 []OLLogisticUserRoleModel `json:"roles"`
}

func (m *OLLogisticUserInfoModel) GetID() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ID
}
func (m *OLLogisticUserInfoModel) GetName() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Name
}
func (m *OLLogisticUserInfoModel) GetLogin() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Login
}
func (m *OLLogisticUserInfoModel) GetPhoneNumber() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.PhoneNumber
}
func (m *OLLogisticUserInfoModel) GetFreelancerEmployeeID() interface{} {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.FreelancerEmployeeID
}
func (m *OLLogisticUserInfoModel) GetSupplierID() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.SupplierID
}
func (m *OLLogisticUserInfoModel) GetSupplierName() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.SupplierName
}
func (m *OLLogisticUserInfoModel) GetRoleID() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.RoleID
}
func (m *OLLogisticUserInfoModel) IsBoss() bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.IsBoss
}
func (m *OLLogisticUserInfoModel) IsDriver() bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.IsDriver
}
func (m *OLLogisticUserInfoModel) IsDeleted() bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.IsDeleted
}
func (m *OLLogisticUserInfoModel) IsVerified() bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.IsVerified
}
func (m *OLLogisticUserInfoModel) IsNewFreelancer() bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.IsNewFreelancer
}
func (m *OLLogisticUserInfoModel) GetSupplierVatID() interface{} {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.SupplierVatID
}
func (m *OLLogisticUserInfoModel) GetSupplierBankAccountID() interface{} {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.SupplierBankAccountID
}
func (m *OLLogisticUserInfoModel) GetSupplierContractDt() interface{} {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.SupplierContractDt
}
func (m *OLLogisticUserInfoModel) GetRoles() []OLLogisticUserRoleModel {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Roles
}

func (m *OLLogisticUserInfoModel) SetID(id int) {
	m.mtx.Lock()
	m.hidden.ID = id
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetName(name string) {
	m.mtx.Lock()
	m.hidden.Name = name
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetLogin(login string) {
	m.mtx.Lock()
	m.hidden.Login = login
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetPhoneNumber(phone string) {
	m.mtx.Lock()
	m.hidden.PhoneNumber = phone
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetFreelancerEmployeeID(id interface{}) {
	m.mtx.Lock()
	m.hidden.FreelancerEmployeeID = id
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetSupplierID(id int) {
	m.mtx.Lock()
	m.hidden.SupplierID = id
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetSupplierName(name string) {
	m.mtx.Lock()
	m.hidden.SupplierName = name
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetRoleID(id int) {
	m.mtx.Lock()
	m.hidden.RoleID = id
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetIsBoss(b bool) {
	m.mtx.Lock()
	m.hidden.IsBoss = b
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetIsDriver(b bool) {
	m.mtx.Lock()
	m.hidden.IsDriver = b
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetIsDeleted(b bool) {
	m.mtx.Lock()
	m.hidden.IsDeleted = b
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetIsVerified(b bool) {
	m.mtx.Lock()
	m.hidden.IsVerified = b
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetIsNewFreelancer(b bool) {
	m.mtx.Lock()
	m.hidden.IsNewFreelancer = b
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetSupplierVatID(id interface{}) {
	m.mtx.Lock()
	m.hidden.SupplierVatID = id
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetSupplierBankAccountID(id interface{}) {
	m.mtx.Lock()
	m.hidden.SupplierBankAccountID = id
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetSupplierContractDt(dt interface{}) {
	m.mtx.Lock()
	m.hidden.SupplierContractDt = dt
	m.mtx.Unlock()
}
func (m *OLLogisticUserInfoModel) SetRoles(roles []OLLogisticUserRoleModel) {
	m.mtx.Lock()
	m.hidden.Roles = roles
	m.mtx.Unlock()
}

func (m *OLLogisticUserInfoModel) SetAll(
	id int,
	name, login, phoneNumber string,
	freelancerEmployeeID interface{},
	supplierID int, supplierName string,
	roleID int,
	isBoss, isDriver, isDeleted, isVerified, isNewFreelancer bool,
	supplierVatID, supplierBankAccountID, supplierContractDt interface{},
	roles []OLLogisticUserRoleModel,
) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.hidden.ID = id
	m.hidden.Name = name
	m.hidden.Login = login
	m.hidden.PhoneNumber = phoneNumber
	m.hidden.FreelancerEmployeeID = freelancerEmployeeID
	m.hidden.SupplierID = supplierID
	m.hidden.SupplierName = supplierName
	m.hidden.RoleID = roleID
	m.hidden.IsBoss = isBoss
	m.hidden.IsDriver = isDriver
	m.hidden.IsDeleted = isDeleted
	m.hidden.IsVerified = isVerified
	m.hidden.IsNewFreelancer = isNewFreelancer
	m.hidden.SupplierVatID = supplierVatID
	m.hidden.SupplierBankAccountID = supplierBankAccountID
	m.hidden.SupplierContractDt = supplierContractDt
	m.hidden.Roles = roles
}

func (m *OLLogisticUserInfoModel) MarshalJSON() ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return json.Marshal(m.hidden)
}

func (m *OLLogisticUserInfoModel) UnmarshalJSON(data []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return json.Unmarshal(data, &m.hidden)
}

type OLLogisticUserRoleModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//// Service params

type OLLogisticTTlParams struct {
	UserInfo       time.Duration
	Balance        time.Duration
	WaySheet       time.Duration
	TransportTypes time.Duration
}

type OLLogisticGetWaySheetsParamsRequest struct {
	From                                time.Time `json:"from"`
	To                                  time.Time `json:"to"`
	OfficeID                            int       `json:"officeId,omitempty"`
	PaymentTypeID                       int       `json:"paymentTypeId,omitempty"`
	RouteID                             int       `json:"routeId,omitempty"`
	SupplierID                          int       `json:"supplierId,omitempty"`
	IsNewMotivation                     bool      `json:"isNewMotivation"`
	IsNotConfirmedWithTransportRequests bool      `json:"isNotConfirmedWithTransportRequests,omitempty"`
}
