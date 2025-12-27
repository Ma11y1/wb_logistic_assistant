package response

import (
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

// AuthCodeResponse Returned when the user submits a phone number to receive an access code
//
// URL: https://drive.wb.ru/user-management/api/v1/public/registration/code
type AuthCodeResponse struct {
	Error *errors.AuthAPIError `json:"error"` // errors.ErrorTypeInvalidRequestBody
	Meta  *models.AuthMeta     `json:"meta"`
	Data  *models.AuthCode     `json:"data"`
}

// AuthResponse Returned when the user enters a passcode. Contains an access token that can be exchanged for a data access token
//
// URL: https://drive.wb.ru/user-management/api/v1/public/registration/auth
type AuthResponse struct {
	Error *errors.AuthAPIError    `json:"error"` // errors.ErrorTypeInvalidRequestBody or errors.ErrorTypeInvalidAuthCode  or errors.ErrorTypeEmptyAuthCode
	Meta  *models.AuthMeta        `json:"meta"`
	Data  *models.AuthAccessToken `json:"data"`
}

// AuthMergeResponse Returned in exchange for the access token received earlier. This token is used to access data
//
// URL: https://drive.wb.ru/user-management/api/v1/public/token/merge
type AuthMergeResponse struct {
	Error *errors.AuthAPIError     `json:"error"` // errors.ErrorTypeInvalidRequestBody or errors.ErrorTypeErrorMerge
	Data  *models.AuthSessionToken `json:"data"`
}
