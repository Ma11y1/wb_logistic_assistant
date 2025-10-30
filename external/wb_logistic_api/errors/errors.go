package errors

import "fmt"

type ErrorType string

const (
	ErrorTypeNone               ErrorType = "um:0"
	ErrorTypeInvalidRequestBody ErrorType = "um:1"
	ErrorTypeErrorMerge         ErrorType = "um:2" // Returned when exchanging an access token for a data token fails
	ErrorTypeNeedCaptcha        ErrorType = "um:wb_3"
	ErrorTypeAuthRequestLimit   ErrorType = "um:wb_4"  // Waiting before resending. Returns the timeout for the next request
	ErrorTypeInvalidAuthCode    ErrorType = "um:wb_6"  // Returned when an incorrect authorization code received by phone number is sent
	ErrorTypeEmptyAuthCode      ErrorType = "um:wb_30" // Returned when a lot of time has passed since the authorization code was received and sent to the server
)

type AuthAPIError struct {
	Err     string    `json:"error"`
	Code    ErrorType `json:"code"`
	Message string    `json:"message"`
}

func (e *AuthAPIError) Error() string {
	return fmt.Sprintf("WB logistic auth API error [%s]: %s: %s ", e.Code, e.Err, e.Message)
}

type APIError struct {
	Err  string `json:"error"`
	Code int    `json:"code"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("WB logistic API error [%d]: %s ", e.Code, e.Err)
}
