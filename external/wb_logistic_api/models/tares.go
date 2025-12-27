package models

import "time"

type TareForOffice struct {
	ID              int       `json:"tare_id"`
	Type            string    `json:"tare_type"`
	CountBarcodes   int       `json:"count_shk"`
	DstOfficeID     int       `json:"dst_office_id"`
	DstOfficeName   string    `json:"dst_office_name"`
	LastOperationDt time.Time `json:"last_operation_dt"`
	SpName          string    `json:"sp_name"`
}

type WaySheetTare struct {
	ID              string    `json:"tare_id"`
	CountBarcodes   int       `json:"count_shk"`
	IsReturn        bool      `json:"is_return"`
	TareType        string    `json:"tare_type"`
	BoxVolume       float64   `json:"box_volume"`
	DtArrival       time.Time `json:"dt_arrival"`
	TareReturnType  string    `json:"tare_return_type"`
	OfficeArrivalID string    `json:"office_arrival_id"`
	DstOfficeID     string    `json:"dst_office_id"`
}
