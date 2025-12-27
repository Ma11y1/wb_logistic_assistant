package models

import (
	"encoding/json"
	"sync"
	"time"
)

//// Main

type WBLogisticModel struct {
	mtx               sync.RWMutex
	Login             string                       `json:"login"`
	AccessToken       *WBLogisticAccessTokenModel  `json:"access_token"`
	MergedAccessToken *WBLogisticSessionTokenModel `json:"merged_token"`
	UserInfo          *WBLogisticUserInfoModel     `json:"user_info"`
}

func (m *WBLogisticModel) GetLogin() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.Login
}

func (m *WBLogisticModel) GetAccessToken() *WBLogisticAccessTokenModel {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.AccessToken
}

func (m *WBLogisticModel) SetAccessToken(token *WBLogisticAccessTokenModel) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.AccessToken = token
}

func (m *WBLogisticModel) GetMergedAccessToken() *WBLogisticSessionTokenModel {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.MergedAccessToken
}

func (m *WBLogisticModel) SetLogin(login string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.Login = login
}

func (m *WBLogisticModel) SetMergedAccessToken(token *WBLogisticSessionTokenModel) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.MergedAccessToken = token
}

func (m *WBLogisticModel) GetUserInfo() *WBLogisticUserInfoModel {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.UserInfo
}

func (m *WBLogisticModel) SetUserInfo(info *WBLogisticUserInfoModel) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.UserInfo = info
}

//// Access tokens

type WBLogisticAccessTokenModel struct {
	mtx    sync.RWMutex
	hidden hiddenWBLogisticAccessTokenModel
}

type hiddenWBLogisticAccessTokenModel struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (m *WBLogisticAccessTokenModel) SetAll(tokenType, accessToken, refreshToken string, expiresIn int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.TokenType = tokenType
	m.hidden.AccessToken = accessToken
	m.hidden.RefreshToken = refreshToken
	m.hidden.ExpiresIn = expiresIn
}

func (m *WBLogisticAccessTokenModel) GetTokenType() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.TokenType
}
func (m *WBLogisticAccessTokenModel) SetTokenType(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.TokenType = value
}

func (m *WBLogisticAccessTokenModel) GetAccessToken() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.AccessToken
}
func (m *WBLogisticAccessTokenModel) SetAccessToken(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.AccessToken = value
}

func (m *WBLogisticAccessTokenModel) GetRefreshToken() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.RefreshToken
}
func (m *WBLogisticAccessTokenModel) SetRefreshToken(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.RefreshToken = value
}

func (m *WBLogisticAccessTokenModel) GetExpiresIn() int64 {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ExpiresIn
}
func (m *WBLogisticAccessTokenModel) SetExpiresIn(value int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ExpiresIn = value
}

func (m *WBLogisticAccessTokenModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.hidden)
}
func (m *WBLogisticAccessTokenModel) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.hidden)
}

//// Session token

type WBLogisticSessionTokenModel struct {
	mtx    sync.RWMutex
	hidden hiddenWBLogisticSessionTokenModel
}

type hiddenWBLogisticSessionTokenModel struct {
	TokenType   string `json:"token_type"`
	Source      string `json:"source"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (m *WBLogisticSessionTokenModel) SetAll(tokenType, source, accessToken string, expiresIn int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.TokenType = tokenType
	m.hidden.Source = source
	m.hidden.AccessToken = accessToken
	m.hidden.ExpiresIn = expiresIn
}

func (m *WBLogisticSessionTokenModel) GetSource() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Source
}
func (m *WBLogisticSessionTokenModel) SetSource(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Source = value
}

func (m *WBLogisticSessionTokenModel) GetAccessToken() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.AccessToken
}
func (m *WBLogisticSessionTokenModel) SetAccessToken(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.AccessToken = value
}

func (m *WBLogisticSessionTokenModel) GetExpiresIn() int64 {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ExpiresIn
}
func (m *WBLogisticSessionTokenModel) SetExpiresIn(value int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ExpiresIn = value
}

func (m *WBLogisticSessionTokenModel) GetTokenType() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.TokenType
}
func (m *WBLogisticSessionTokenModel) SetTokenType(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.TokenType = value
}

func (m *WBLogisticSessionTokenModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.hidden)
}
func (m *WBLogisticSessionTokenModel) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.hidden)
}

//// User info

type WBLogisticUserInfoModel struct {
	mtx    sync.RWMutex
	hidden hiddenWBLogisticUserInfoModel
}

type hiddenWBLogisticUserInfoModel struct {
	ID           int                             `json:"id"`
	Verified     bool                            `json:"verified"`
	RoleIDs      []string                        `json:"role_ids"`
	Roles        []*WBLogisticUserInfoRoleModel  `json:"roles"`
	Permissions  []string                        `json:"permissions"`
	UserDetails  *WBLogisticUserInfoDetailsModel `json:"user_details"`
	DriverRoleID int                             `json:"driver_role_id"`
}

func (m *WBLogisticUserInfoModel) SetAll(id int, verified bool, roleIDs []string, roles []*WBLogisticUserInfoRoleModel, permissions []string, userDetails *WBLogisticUserInfoDetailsModel, driverRoleID int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ID = id
	m.hidden.Verified = verified
	m.hidden.RoleIDs = roleIDs
	m.hidden.Roles = roles
	m.hidden.Permissions = permissions
	m.hidden.UserDetails = userDetails
	m.hidden.DriverRoleID = driverRoleID
}

func (m *WBLogisticUserInfoModel) GetID() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.ID
}
func (m *WBLogisticUserInfoModel) SetID(value int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.ID = value
}

func (m *WBLogisticUserInfoModel) IsVerified() bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Verified
}
func (m *WBLogisticUserInfoModel) SetVerified(value bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Verified = value
}

func (m *WBLogisticUserInfoModel) GetRoleIDs() []string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.RoleIDs
}
func (m *WBLogisticUserInfoModel) SetRoleIDs(value []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.RoleIDs = value
}

func (m *WBLogisticUserInfoModel) GetRoles() []*WBLogisticUserInfoRoleModel {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Roles
}
func (m *WBLogisticUserInfoModel) SetRoles(value []*WBLogisticUserInfoRoleModel) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Roles = value
}

func (m *WBLogisticUserInfoModel) GetPermissions() []string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Permissions
}
func (m *WBLogisticUserInfoModel) SetPermissions(value []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Permissions = value
}

func (m *WBLogisticUserInfoModel) GetUserDetails() *WBLogisticUserInfoDetailsModel {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.UserDetails
}
func (m *WBLogisticUserInfoModel) SetUserDetails(value *WBLogisticUserInfoDetailsModel) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.UserDetails = value
}

func (m *WBLogisticUserInfoModel) GetDriverRoleID() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.DriverRoleID
}
func (m *WBLogisticUserInfoModel) SetDriverRoleID(value int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.DriverRoleID = value
}

func (m *WBLogisticUserInfoModel) MarshalJSON() ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return json.Marshal(m.hidden)
}
func (m *WBLogisticUserInfoModel) UnmarshalJSON(data []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return json.Unmarshal(data, &m.hidden)
}

//// User info role

type WBLogisticUserInfoRoleModel struct {
	mtx    sync.RWMutex
	hidden hiddenWBLogisticUserInfoRoleModel
}

type hiddenWBLogisticUserInfoRoleModel struct {
	UserRoleUID  string `json:"user_role_uid"`
	UserRoleName string `json:"user_role_name"`
}

func (m *WBLogisticUserInfoRoleModel) SetAll(roleUID, roleName string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.UserRoleUID = roleUID
	m.hidden.UserRoleName = roleName
}

func (m *WBLogisticUserInfoRoleModel) GetUserRoleUID() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.UserRoleUID
}

func (m *WBLogisticUserInfoRoleModel) SetUserRoleUID(id string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.UserRoleUID = id
}

func (m *WBLogisticUserInfoRoleModel) GetUserRoleName() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.UserRoleName
}

func (m *WBLogisticUserInfoRoleModel) SetUserRoleName(name string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.UserRoleName = name
}

func (m *WBLogisticUserInfoRoleModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.hidden)
}
func (m *WBLogisticUserInfoRoleModel) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.hidden)
}

//// User info details

type WBLogisticUserInfoDetailsModel struct {
	mtx    sync.RWMutex
	hidden hiddenWBLogisticUserInfoDetailsModel
}
type hiddenWBLogisticUserInfoDetailsModel struct {
	Name                 string `json:"name"`
	PhoneNumber          string `json:"phone_number"`
	SupplierID           int    `json:"supplier_id"`
	FreelancerEmployeeID int    `json:"freelancer_employee_id"`
	VatID                int    `json:"vat_id"`
	VatName              string `json:"vat_name"`
	Telegram             string `json:"telegram"`
}

func (m *WBLogisticUserInfoDetailsModel) SetAll(name, phoneNumber string, supplierID, FreelancerEmployeeID, vatID int, vatName, telegram string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Name = name
	m.hidden.PhoneNumber = phoneNumber
	m.hidden.SupplierID = supplierID
	m.hidden.FreelancerEmployeeID = FreelancerEmployeeID
	m.hidden.VatID = vatID
	m.hidden.VatName = vatName
	m.hidden.Telegram = telegram
}

func (m *WBLogisticUserInfoDetailsModel) GetName() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Name
}
func (m *WBLogisticUserInfoDetailsModel) SetName(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Name = value
}

func (m *WBLogisticUserInfoDetailsModel) GetPhoneNumber() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.PhoneNumber
}
func (m *WBLogisticUserInfoDetailsModel) SetPhoneNumber(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.PhoneNumber = value
}

func (m *WBLogisticUserInfoDetailsModel) GetSupplierID() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.SupplierID
}
func (m *WBLogisticUserInfoDetailsModel) SetSupplierID(value int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.SupplierID = value
}

func (m *WBLogisticUserInfoDetailsModel) GetFreelancerEmployeeID() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.FreelancerEmployeeID
}
func (m *WBLogisticUserInfoDetailsModel) SetFreelancerEmployeeID(value int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.FreelancerEmployeeID = value
}

func (m *WBLogisticUserInfoDetailsModel) GetVatID() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.VatID
}
func (m *WBLogisticUserInfoDetailsModel) SetVatID(value int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.VatID = value
}

func (m *WBLogisticUserInfoDetailsModel) GetVatName() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.VatName
}
func (m *WBLogisticUserInfoDetailsModel) SetVatName(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.VatName = value
}

func (m *WBLogisticUserInfoDetailsModel) GetTelegram() string {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.hidden.Telegram
}
func (m *WBLogisticUserInfoDetailsModel) SetTelegram(value string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.hidden.Telegram = value
}

func (m *WBLogisticUserInfoDetailsModel) MarshalJSON() ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return json.Marshal(m.hidden)
}
func (m *WBLogisticUserInfoDetailsModel) UnmarshalJSON(data []byte) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return json.Unmarshal(data, &m.hidden)
}

//// Params

type WBLogisticTTlParams struct {
	UserInfo                        time.Duration
	RemainsLastMileReports          time.Duration
	RemainsLastMileReportsRouteInfo time.Duration
	JobsScheduling                  time.Duration
	ShipmentInfo                    time.Duration
	ShipmentTransfers               time.Duration
	WaySheetInfo                    time.Duration
	WaySheetFinanceDetails          time.Duration
}

type WBLogisticGetShipmentsParamsRequest struct {
	DataStart                time.Time   `json:"dt_start"`
	DataEnd                  time.Time   `json:"dt_end"`
	SrcOfficeID              int         `json:"src_office_id"`
	PageIndex                int         `json:"page_index"`
	Limit                    int         `json:"limit"`
	SupplierID               int         `json:"supplier_id"`
	Direction                int         `json:"direction"`
	Sorter                   string      `json:"sorter"` // Types: "id", "supplier", "driver", "route_car_id", "updated_at"
	FilterShipmentID         int         `json:"shipment_id"`
	FilterShipmentType       string      `json:"shipment_type"` // Types: "last-mile" or "truck"
	FilterVehicleNumberPlate int         `json:"vehicle_number_plate"`
	FilterDstOffice          interface{} `json:"dst_office"` // id or name
	FilterShowOnlyOpen       bool        `json:"show_only_open"`
}

type WBLogisticGetWaySheetsParamsRequest struct {
	DateClose          time.Time `json:"date_close"`
	DateOpen           time.Time `json:"date_open"`
	SupplierID         int       `json:"supplier_id"`
	SrcOfficeID        int       `json:"src_office_id"`
	RouteCarID         int       `json:"routecar_id"`
	Limit              int       `json:"limit"`
	Offset             int       `json:"offset"`
	WayTypeID          int       `json:"way_type_id"`
	VehicleNumberPlate string    `json:"vehicle_number_plate"`
}
