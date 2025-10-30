package prompters

import "time"

type InitializeAppPrompter interface {
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

	PromptOLLogisticAuthStart()
	PromptOLLogisticQuestionAuthNewUser() bool
	PromptOLLogisticRequestAuthData() (string, string, error)
	PromptOLLogisticInvalidAuthData()
	PromptOLLogisticAuthFailed()
	PromptOLLogisticAuthStorageFailed()
	PromptOLLogisticAuthSuccessful(login, username string)

	PromptWBLogisticAuthStart()
	PromptWBLogisticQuestionAuthNewUser() bool
	PromptWBLogisticRequestAuthLogin() (string, error)
	PromptWBLogisticRequestAuthCode(method string, time int) (int, error)
	PromptWBLogisticInvalidAuthCode()
	PromptWBLogisticAuthFailed()
	PromptWBLogisticAuthStorageFailed()
	PromptWBLogisticAuthSuccessful(login, username string)

	PromptTelegramBotAuthStart()
	PromptTelegramBotQuestionAuthNewBot() bool
	PromptTelegramBotRequestToken() (string, error)
	PromptTelegramBotInitStorageFailed()
	PromptTelegramBotInitFailed()
	PromptTelegramBotAuthSuccessful(name string)
}

type GeneralRoutesReporterPrompter interface {
	PromptError(message string)
	PromptRender()
}

type ShipmentCloseReporterPrompter interface {
	PromptError(message string)
	PromptRender(duration time.Duration)
}
