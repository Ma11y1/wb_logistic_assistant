package response

import (
	"encoding/json"
	"wb_logistic_assistant/external/wb_logistic_api/errors"
	"wb_logistic_assistant/external/wb_logistic_api/models"
)

type GetShipmentsResponse struct {
	Meta  *models.ShipmentsMeta `json:"meta"`
	Data  []*models.Shipment    `json:"data"`
	Error *errors.APIError      `json:"error"`
}

func (r *GetShipmentsResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Meta  *models.ShipmentsMeta `json:"meta"`
		Data  []*models.Shipment    `json:"data"`
		Error *errors.APIError      `json:"error"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}

type GetShipmentInfoResponse struct {
	Data  *models.ShipmentInfo `json:"data"`
	Error *errors.APIError     `json:"error"`
}

func (r *GetShipmentInfoResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Data  *models.ShipmentInfo `json:"data"`
		Error *errors.APIError     `json:"error"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}

type GetShipmentTransfersResponse struct {
	Data  *models.ShipmentTransfers `json:"data"`
	Error *errors.APIError          `json:"error"`
}

func (r *GetShipmentTransfersResponse) UnmarshalJSON(data []byte) error {
	var apiErr errors.APIError
	if err := json.Unmarshal(data, &apiErr); err == nil {
		if apiErr.Code != 0 || apiErr.Err != "" {
			r.Error = &apiErr
			return nil
		}
	}

	var temp struct {
		Data  *models.ShipmentTransfers `json:"data"`
		Error *errors.APIError          `json:"error"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*r = temp
	return nil
}
