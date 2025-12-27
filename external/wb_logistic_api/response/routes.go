package response

import (
	"encoding/json"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

// GetAssociationRoutesInfoByNameResponse Get list of routes by name or part of name
//
// URL: https://logistics.wb.ru/routes-netcore-service/api/v1/route/
type GetAssociationRoutesInfoByNameResponse struct {
	Error *errors.APIError `json:"error"`
	Data  []*models.AssociationRouteInfoByName
}

func (r *GetAssociationRoutesInfoByNameResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp []*models.AssociationRouteInfoByName
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	r.Data = temp
	return nil
}
