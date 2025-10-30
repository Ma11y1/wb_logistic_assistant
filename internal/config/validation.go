package config

import (
	"wb_logistic_assistant/internal/errors"
)

func validation(config *Config) error {
	if config == nil {
		return errors.Newf("Config.validation()", "config is nil")
	}
	if err := validationReports(config.reports); err != nil {
		return errors.Newf("config.validation()", "config 'reports' validation failed: %v", err)
	}
	if err := validationStorage(config.storage); err != nil {
		return errors.Newf("config.validation()", "config 'storage' validation failed: %v", err)
	}
	if err := validationLogistic(config.logistic); err != nil {
		return errors.Newf("config.validation()", "config 'logistic' validation failed: %v", err)
	}
	if err := validationGoogleSheets(config.googleSheets); err != nil {
		return errors.Newf("config.validation()", "config 'google_sheets' validation failed: %v", err)
	}
	if err := validationTelegramBot(config.telegram); err != nil {
		return errors.Newf("config.validation()", "config 'telegramBot' validation failed: %v", err)
	}
	return nil
}

func validationReports(config *Reports) error {
	if config == nil {
		return errors.New("config.validationReports()", "config is nil")
	}

	generalRoutes := config.generalRoutes
	if generalRoutes == nil {
		return errors.New("config.validationReports()", "'general_routes' is nil")
	}
	if generalRoutes.pollingInterval <= 0 {
		return errors.New("config.validationReports()", "'general_routes.polling_interval' is invalid, it must be > 0")
	}
	if generalRoutes.errRetryTaskLimit <= 0 {
		return errors.New("config.validationReports()", "'general_routes.err_retry_limit' is invalid, it must be > 0")
	}
	if generalRoutes.taskTimeout <= 0 {
		return errors.New("config.validationReports()", "'general_routes.task_timeout' is invalid, it must be > 0")
	}
	if generalRoutes.ttlResetChangeBarcodes < 0 {
		return errors.New("config.validationReports()", "'general_routes.ttl_reset_change_barcodes' is it must be >= 0")
	}
	if generalRoutes.ttlLoadRemainsLastMileReport < 0 {
		return errors.New("config.validationReports()", "'general_routes.ttl_load_remains_last_mile_report' is it must be >= 0")
	}
	if generalRoutes.ttlLoadJobsScheduling < 0 {
		return errors.New("config.validationReports()", "'general_routes.ttl_load_jobs_scheduling' is it must be >= 0")
	}
	if generalRoutes.ttlLoadShipments < 0 {
		return errors.New("config.validationReports()", "'general_routes.ttl_load_shipments' is it must be >= 0")
	}
	if generalRoutes.sortColumn < 0 {
		return errors.New("config.validationReports()", "'general_routes.sort_column' is invalid")
	}
	if generalRoutes.skipRoutes == nil {
		return errors.New("config.validationReports()", "'general_routes.skip_routes' is nil")
	}

	shipmentsClose := config.shipmentClose
	if shipmentsClose == nil {
		return errors.New("config.validationReports()", "'shipment_close' is nil")
	}
	if shipmentsClose.pollingInterval <= 0 {
		return errors.New("config.validationReports()", "'shipment_close.polling_interval' is invalid, it must be > 0")
	}
	if shipmentsClose.errRetryTaskLimit <= 0 {
		return errors.New("config.validationReports()", "'shipment_close.err_retry_limit' is invalid, it must be > 0")
	}
	if shipmentsClose.taskTimeout <= 0 {
		return errors.New("config.validationReports()", "'shipment_close.task_timeout' is invalid, it must be > 0")
	}
	if shipmentsClose.skipRoutes == nil {
		return errors.New("config.validationReports()", "'shipment_close.skip_routes' is nil")
	}

	return nil
}

func validationStorage(config *Storage) error {
	if config == nil {
		return errors.New("config.validationStorage()", "config is nil")
	}
	if config.path == "" {
		return errors.New("config.validationStorage()", "'path' is empty")
	}
	return nil
}

func validationLogistic(config *Logistic) error {
	if config == nil {
		return errors.New("config.validationLogistic()", "config is nil")
	}

	if config.wbClient == nil {
		return errors.New("config.validationLogistic()", "'wb_client' is nil")
	}
	if config.wbClient.userAgent == "" {
		return errors.New("config.validationLogistic()", "'wb_client.user_agent' is empty")
	}
	if config.wbClient.secUserAgent == "" {
		return errors.New("config.validationLogistic()", "'wb_client.sec_user_agent' is empty")
	}
	if config.wbClient.platform == "" {
		return errors.New("config.validationLogistic()", "'wb_client.platform' is empty")
	}

	if config.office == nil {
		return errors.New("config.validationLogistic()", "'office' is nil")
	}
	if config.office.id <= 0 {
		return errors.New("config.validationLogistic()", "'office.id' is invalid, it must be > 0")
	}
	if config.office.suppliers == nil || len(config.office.suppliers) == 0 {
		return errors.New("config.validationLogistic()", "'office.suppliers' is empty")
	}
	if config.office.parking == nil {
		return errors.New("config.validationLogistic()", "'office.parking' is empty")
	}

	cacheTTL := config.cacheTTL
	if cacheTTL == nil {
		return errors.New("config.validationLogistic()", "'ttl' is nil")
	}
	if cacheTTL.userInfo < 0 {
		return errors.New("config.validationLogistic()", "'ttl.wb_user_info' is invalid, it must be >= 0")
	}
	if cacheTTL.remainsLastMileReports < 0 {
		return errors.New("config.validationLogistic()", "'ttl.wb_remains_last_mile_reports' is invalid, it must be >= 0")
	}
	if cacheTTL.remainsLastMileReportsRouteInfo < 0 {
		return errors.New("config.validationLogistic()", "'ttl.wb_remains_last_mile_reports_route_info' is invalid, it must be >= 0")
	}
	if cacheTTL.jobsScheduling < 0 {
		return errors.New("config.validationLogistic()", "'ttl.wb_jobs_scheduling' is invalid, it must be >= 0")
	}
	if cacheTTL.shipmentInfo < 0 {
		return errors.New("config.validationLogistic()", "'ttl.shipment_info' is invalid, it must be >= 0")
	}
	if cacheTTL.shipmentTransfers < 0 {
		return errors.New("config.validationLogistic()", "'ttl.shipment_transfers' is invalid, it must be >= 0")
	}
	if cacheTTL.waySheetInfo < 0 {
		return errors.New("config.validationLogistic()", "'ttl.way_sheet_info' is invalid, it must be >= 0")
	}

	return nil
}

func validationGoogleSheets(config *GoogleSheets) error {
	if config == nil {
		return errors.New("config.validationGoogleSheets()", "config is empty")
	}

	client := config.client
	if client == nil {
		return errors.New("config.validationGoogleSheets()", "'client' is empty")
	}
	if client.isOAuth && client.oauthCredentials == "" {
		return errors.New("config.validationGoogleSheets()", "'client.oauth_credentials' is empty")
	}
	if client.isOAuth && client.secondsWaitServer <= 0 {
		return errors.New("config.validationGoogleSheets()", "'client.seconds_wait_server' is == 0 or < 0, it must be > 0")
	}
	if !client.isOAuth && client.serviceCredentials == "" {
		return errors.New("config.validationGoogleSheets()", "'client.service_credentials' is empty")
	}

	if config.reportSheets == nil {
		return errors.New("config.validationGoogleSheets()", "'report_sheets' is nil")
	}

	generalRoutesReport := config.reportSheets.generalRoutes
	if generalRoutesReport == nil {
		return errors.New("config.validationGoogleSheets()", "'report_sheets.general_routes' is nil")
	}
	if generalRoutesReport.spreadsheetID == "" {
		return errors.New("config.validationGoogleSheets()", "'report_sheets.general_routes.spreadsheet_id' is empty")
	}
	if generalRoutesReport.sheetName == "" {
		return errors.New("config.validationGoogleSheets()", "'report_sheets.general_routes.sheet_name' is empty")
	}

	shipmentClose := config.reportSheets.shipmentClose
	if shipmentClose == nil {
		return errors.New("config.validationGoogleSheets()", "'report_sheets.shipment_close' is nil")
	}
	if shipmentClose.spreadsheetID == "" {
		return errors.New("config.validationGoogleSheets()", "'report_sheets.shipment_close.spreadsheet_id' is empty")
	}
	if shipmentClose.sheetName == "" {
		return errors.New("config.validationGoogleSheets()", "'report_sheets.shipment_close.sheet_name' is empty")
	}

	return nil
}

func validationTelegramBot(config *TelegramBot) error {
	if config == nil {
		return errors.New("config.validationTelegramBot()", "config is nil")
	}
	if config.chatID == 0 {
		return errors.New("config.validationTelegramBot()", "'chat_id' is it not must be 0")
	}
	return nil
}
