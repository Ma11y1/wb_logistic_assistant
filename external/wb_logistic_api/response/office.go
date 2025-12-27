package response

import (
	"encoding/json"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

// GetAssociationOfficesInfoByNameResponse Get list of offices by name or part of name
//
// URL: https://logistics.wb.ru/routes-netcore-service/api/v1/office
type GetAssociationOfficesInfoByNameResponse struct {
	Error *errors.APIError `json:"error"`
	Data  []*models.AssociationOfficeInfoByName
}

func (r *GetAssociationOfficesInfoByNameResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp []*models.AssociationOfficeInfoByName
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	r.Data = temp
	return nil
}
