package prompters

import "time"

type InitializeAppPrompter interface {
	InitializeWBLogisticPrompter
	InitializeGoogleSheetsPrompter
	InitializeTelegramBotPrompter
	PromptInitFinish()
}

type InitializeWBLogisticPrompter interface {
	PromptWBLogisticAuthStart()
	PromptWBLogisticQuestionAuthNewUser() bool
	PromptWBLogisticRequestAuthLogin() string
	PromptWBLogisticInvalidAuthLogin()
	PromptWBLogisticRequestAuthCode(method string, time int) int
	PromptWBLogisticInvalidAuthCode()
	PromptWBLogisticRequestAccessTokenData() string
	PromptWBLogisticInvalidAccessTokenData()
	PromptWBLogisticAuthFailed()
	PromptWBLogisticAuthStorageFailed()
	PromptWBLogisticAuthSuccessful(login, username string)
}

type InitializeGoogleSheetsPrompter interface {
	PromptGoogleSheetsAuthStart()
	PromptGoogleSheetsQuestionAuthNewCredentials() bool
	PromptGoogleSheetsRequestAuthCodeAuto(url string, seconds int)
	PromptGoogleSheetsRequestAuthCode(url string) (string, error)
	PromptGoogleSheetsInvalidAuthCode()
	PromptGoogleSheetsReadCredentialsFailed()
	PromptGoogleSheetsAuthAutoFailed()
	PromptGoogleSheetsAuthFailed()
	PromptGoogleSheetsAuthStorageFailed()
	PromptGoogleSheetsAuthSuccessful()
}

type InitializeTelegramBotPrompter interface {
	PromptTelegramBotAuthStart()
	PromptTelegramBotQuestionAuthNewBot() bool
	PromptTelegramBotRequestToken() (string, error)
	PromptTelegramBotInitStorageFailed()
	PromptTelegramBotInitFailed()
	PromptTelegramBotAuthSuccessful(name string)
}

type GeneralRoutesReporterPrompter interface {
	PromptStart()
	PromptFinish(duration time.Duration)
	PromptUpdateRoutes(count int)
	PromptUpdateShipments()
	PromptUpdateRating()
	PromptCloseShipment(id, remainsBarcodes int)
	PromptUpdateWaySheets()
	PromptSendReport()
	PromptError(message string)
}

type ShipmentCloseReporterPrompter interface {
	PromptStart()
	PromptFinish(duration time.Duration)
	PromptShipmentOpened(routeID, shipmentID, opened int)
	PromptShipmentClose(routeID, shipmentID int)
	PromptSendReport(routeID, shipmentID, WaySheetID int)
	PromptError(message string)
}

type FinanceRoutesReporterPrompter interface {
	PromptStart()
	PromptFinish(duration time.Duration)
	PromptCountWaySheet(count int)
	PromptCloseWaySheet(routeID int, waySheetID, shipmentID string)
	PromptSendReport(routeID int, waySheetID, shipmentID string)
	PromptError(message string)
}

type FinanceDailyReporterPrompter interface {
	PromptStart(date time.Time)
	PromptFinish(duration time.Duration)
	PromptCountWaySheet(total, closed, opened int)
	PromptSendReport(routeID int)
	PromptError(message string)
}
