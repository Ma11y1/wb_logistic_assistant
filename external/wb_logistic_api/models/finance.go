package models

type WaySheetFinanceDetails struct {
	Currency        string                                  `json:"currency"`
	WaySheetID      string                                  `json:"way_sheet_id"`
	PointsWithPrice []*WaySheetFinanceDetailsPointWithPrice `json:"points_with_price"`
	ReturnPrice     *WaySheetFinanceDetailsReturnPrice      `json:"return_price"`
	TotalPrice      *WaySheetFinanceDetailsTotalPrice       `json:"total_price"`
}

type WaySheetFinanceDetailsPointWithPrice struct {
	BonusSum  *WaySheetFinanceDetailsBonusSum  `json:"bonus_sum"`
	FineSum   *WaySheetFinanceDetailsFineSum   `json:"fine_sum"`
	OfficeSum *WaySheetFinanceDetailsOfficeSum `json:"office_sum"`
	OfficeID  string                           `json:"office_id"`
}

type WaySheetFinanceDetailsBonusSum struct {
	Value string `json:"value"`
}

type WaySheetFinanceDetailsFineSum struct {
	Value string `json:"value"`
}

type WaySheetFinanceDetailsOfficeSum struct {
	Value string `json:"value"`
}

type WaySheetFinanceDetailsReturnPrice struct {
	Value string `json:"value"`
}

type WaySheetFinanceDetailsTotalPrice struct {
	Value string `json:"value"`
}
