package response

import (
	"encoding/json"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

type GetTaresForOfficesResponse struct {
	Data  []*models.TareForOffice `json:"data"`
	Error *errors.APIError        `json:"error"`
}

func (r *GetTaresForOfficesResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Data  []*models.TareForOffice `json:"data"`
		Error *errors.APIError        `json:"error"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}
