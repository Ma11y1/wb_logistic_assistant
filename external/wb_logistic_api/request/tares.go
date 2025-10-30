package request

import (
	"context"
	"fmt"
	"wb_logistic_assistant/external/wb_logistic_api/response"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

// GetTaresForOffices
//
// URL: https://logistics.wb.ru/tares/api/v2/public/tares/tares-for-offices
type GetTaresForOffices struct {
	BaseRequest
}

func NewGetTaresForOffices(client transport.HTTPClient, token string) *GetTaresForOffices {
	r := &GetTaresForOffices{BaseRequest: *NewRequestToken(client, "https://logistics.wb.ru/tares/api/v2/public/tares/tares-for-offices", token)}
	r.parameters.Set("is_drive", false)
	return r
}

func (r *GetTaresForOffices) Do(ctx context.Context) (response response.GetTaresForOfficesResponse, err error) {
	err = r.PostUnmarshal(ctx, &response)
	if err != nil {
		err = fmt.Errorf("GetTaresForOffices.Do: %s", err)
	}
	return
}

func (r *GetTaresForOffices) DestinationOfficeIDs(id []int) *GetTaresForOffices {
	r.parameters.Set("dst_office_ids", id)
	return r
}

func (r *GetTaresForOffices) IsDrive(value bool) *GetTaresForOffices {
	r.parameters.Set("is_drive", value)
	return r
}

func (r *GetTaresForOffices) SourceOfficeID(id int) *GetTaresForOffices {
	r.parameters.Set("src_office_id", id)
	return r
}
