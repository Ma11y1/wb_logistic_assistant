package response

import (
	"encoding/json"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

// GetRemainsLastMileReportsResponse Response containing routing reports for goods shipment
//
// URL: https://logistics.wb.ru/reports-service/api/v1/last-mile
type GetRemainsLastMileReportsResponse struct {
	Error *errors.APIError                `json:"error"`
	Data  []*models.RemainsLastMileReport `json:"data"`
}

func (r *GetRemainsLastMileReportsResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Error *errors.APIError                `json:"error"`
		Data  []*models.RemainsLastMileReport `json:"data"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}

// GetRemainsLastMileReportsRouteInfoResponse Response containing report info by report id
//
// URL: https://logistics.wb.ru/reports-service/api/v1/last-mile/[report id]
type GetRemainsLastMileReportsRouteInfoResponse struct {
	Error *errors.APIError                          `json:"error"`
	Data  []*models.RemainsLastMileReportsRouteInfo `json:"data"`
}

func (r *GetRemainsLastMileReportsRouteInfoResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Error *errors.APIError                          `json:"error"`
		Data  []*models.RemainsLastMileReportsRouteInfo `json:"data"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}
