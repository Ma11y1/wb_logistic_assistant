package response

import (
	"encoding/json"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

// GetWaySheetFinanceDetailsResponse Obtaining financial information from a way sheet
//
// URL: https://drive.wb.ru/client-gateway/api/finance/credeber/v1/payment/details
type GetWaySheetFinanceDetailsResponse struct {
	Error *errors.APIError               `json:"error"`
	Data  *models.WaySheetFinanceDetails `json:"data"`
}

func (r *GetWaySheetFinanceDetailsResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Error *errors.APIError               `json:"error"`
		Data  *models.WaySheetFinanceDetails `json:"data"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}
