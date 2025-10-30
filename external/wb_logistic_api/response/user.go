package response

import (
	"encoding/json"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

// UserGetInfoResponse Obtained as a result of a request to obtain account data
//
// URL: https://logistics.wb.ru/api/v1/public/iam/accounts/(ID)
type UserGetInfoResponse struct {
	Error *errors.APIError `json:"error"`
	Data  *models.UserInfo `json:"data"`
}

func (r *UserGetInfoResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Error *errors.APIError `json:"error"`
		Data  *models.UserInfo `json:"data"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}
