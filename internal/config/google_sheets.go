package config

import "encoding/json"

type GoogleSheets struct {
	client       *GoogleSheetsClient       // ro
	reportSheets *GoogleSheetsReportSheets // ro
}

type googleSheets struct {
	Client       *GoogleSheetsClient       `json:"client"`
	ReportSheets *GoogleSheetsReportSheets `json:"report_sheets"`
}

func newGoogleSheets() *GoogleSheets {
	return &GoogleSheets{
		client:       newGoogleSheetsClient(),       // default
		reportSheets: newGoogleSheetsReportSheets(), // default
	}
}

func (s *GoogleSheets) Client() *GoogleSheetsClient             { return s.client }
func (s *GoogleSheets) ReportSheets() *GoogleSheetsReportSheets { return s.reportSheets }

func (s *GoogleSheets) UnmarshalJSON(b []byte) error {
	temp := &googleSheets{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	s.client = temp.Client
	s.reportSheets = temp.ReportSheets
	return nil
}

func (s *GoogleSheets) MarshalJSON() ([]byte, error) {
	return json.Marshal(&googleSheets{
		Client:       s.client,
		ReportSheets: s.reportSheets,
	})
}

type GoogleSheetsClient struct {
	secondsWaitServer  int    // ro
	oauthCredentials   string // ro
	serviceCredentials string // ro
	isOAuth            bool   // ro
}

type googleSheetsClient struct {
	SecondsWaitServer  int    `json:"seconds_wait_server"`
	OauthCredentials   string `json:"oauth_credentials"`
	ServiceCredentials string `json:"service_credentials"`
	IsOauth            bool   `json:"oauth"`
}

func newGoogleSheetsClient() *GoogleSheetsClient {
	return &GoogleSheetsClient{
		secondsWaitServer:  60,                           // default
		oauthCredentials:   "./oauth_credentials.json",   // default
		serviceCredentials: "./service_credentials.json", // default
		isOAuth:            false,                        // default
	}
}

func (s *GoogleSheetsClient) SecondsWaitServer() int     { return s.secondsWaitServer }
func (s *GoogleSheetsClient) OAuthCredentials() string   { return s.oauthCredentials }
func (s *GoogleSheetsClient) ServiceCredentials() string { return s.serviceCredentials }
func (s *GoogleSheetsClient) IsOAuth() bool              { return s.isOAuth }

func (s *GoogleSheetsClient) UnmarshalJSON(b []byte) error {
	temp := &googleSheetsClient{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	s.secondsWaitServer = temp.SecondsWaitServer
	s.oauthCredentials = temp.OauthCredentials
	s.serviceCredentials = temp.ServiceCredentials
	s.isOAuth = temp.IsOauth
	return nil
}

func (s *GoogleSheetsClient) MarshalJSON() ([]byte, error) {
	return json.Marshal(&googleSheetsClient{
		SecondsWaitServer:  s.secondsWaitServer,
		OauthCredentials:   s.oauthCredentials,
		ServiceCredentials: s.serviceCredentials,
		IsOauth:            s.isOAuth,
	})
}

type GoogleSheetsReportSheets struct {
	generalRoutes *GoogleSheetsReportSheet // ro
	shipmentClose *GoogleSheetsReportSheet // ro
}

type googleSheetsSheetsData struct {
	GeneralRoutes *GoogleSheetsReportSheet `json:"general_routes"`
	ShipmentClose *GoogleSheetsReportSheet `json:"shipment_close"`
}

func newGoogleSheetsReportSheets() *GoogleSheetsReportSheets {
	return &GoogleSheetsReportSheets{
		generalRoutes: newGoogleSheetsReportSheet(), // default
		shipmentClose: newGoogleSheetsReportSheet(), // default
	}
}

func (s *GoogleSheetsReportSheets) GeneralRoutes() *GoogleSheetsReportSheet { return s.generalRoutes }
func (s *GoogleSheetsReportSheets) ShipmentClose() *GoogleSheetsReportSheet { return s.shipmentClose }

func (s *GoogleSheetsReportSheets) UnmarshalJSON(b []byte) error {
	temp := &googleSheetsSheetsData{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	s.generalRoutes = temp.GeneralRoutes
	s.shipmentClose = temp.ShipmentClose
	return nil
}

func (s *GoogleSheetsReportSheets) MarshalJSON() ([]byte, error) {
	return json.Marshal(&googleSheetsSheetsData{
		GeneralRoutes: s.generalRoutes,
		ShipmentClose: s.shipmentClose,
	})
}

type GoogleSheetsReportSheet struct {
	spreadsheetID string // ro
	sheetName     string // ro
}

type googleSheetsGeneralRoutes struct {
	SpreadsheetID string `json:"spreadsheet_id"`
	SheetName     string `json:"sheet_name"`
}

func newGoogleSheetsReportSheet() *GoogleSheetsReportSheet {
	return &GoogleSheetsReportSheet{
		spreadsheetID: "UndefinedID", // default
		sheetName:     "Лист1",       // default
	}
}

func (s *GoogleSheetsReportSheet) SpreadsheetID() string { return s.spreadsheetID }
func (s *GoogleSheetsReportSheet) SheetName() string     { return s.sheetName }

func (s *GoogleSheetsReportSheet) UnmarshalJSON(b []byte) error {
	temp := &googleSheetsGeneralRoutes{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	s.spreadsheetID = temp.SpreadsheetID
	s.sheetName = temp.SheetName
	return nil
}

func (s *GoogleSheetsReportSheet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&googleSheetsGeneralRoutes{
		SpreadsheetID: s.spreadsheetID,
		SheetName:     s.sheetName,
	})
}
