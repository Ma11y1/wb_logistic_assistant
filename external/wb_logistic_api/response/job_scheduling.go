package response

import (
	"encoding/json"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

// GetJobsSchedulingResponse
//
// URL: https://logistics.wb.ru/transport-planning-service/api/v1/planning/last-mile
type GetJobsSchedulingResponse struct {
	Error *errors.APIError       `json:"error"`
	Data  *models.JobsScheduling `json:"data"`
}

func (r *GetJobsSchedulingResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Error *errors.APIError       `json:"error"`
		Data  *models.JobsScheduling `json:"data"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}
