package errors

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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
	Code    string `json:"code"`
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("WB logistic API error [%s]: %s %s", e.Code, e.Err, e.Message)
}

func (e *APIError) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Err     string      `json:"error"`
		Code    interface{} `json:"code"`
		Message string      `json:"message"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	e.Err = tmp.Err
	e.Message = tmp.Message

	switch v := tmp.Code.(type) {
	case string:
		e.Code = v
	case int:
		e.Code = strconv.Itoa(v)
	case float64:
		e.Code = strconv.Itoa(int(v))
	case nil:
		e.Code = ""
	default:
		return fmt.Errorf("unexpected type for code: %T", v)
	}

	return nil
}
