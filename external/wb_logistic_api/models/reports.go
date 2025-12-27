package models

import "time"

type RemainsLastMileReports []*RemainsLastMileReport

type RemainsLastMileReport struct {
	OfficeID      int      `json:"office_id"`
	OfficeName    string   `json:"office_name"`
	CountBarcodes int      `json:"count_shk"`
	CountTares    int      `json:"count_tares"`
	TotalVolumeMl int      `json:"total_volume_ml"`
	Routes        []*Route `json:"routes"`
}

type Route struct {
	Name               string               `json:"route_name"`
	CarID              int                  `json:"route_car_id"`
	CarType            string               `json:"route_car_type"`
	Distance           float32              `json:"distance"`
	CountTares         int                  `json:"count_tares"`
	CountShk           int                  `json:"count_shk"`
	ShkLastHours       int                  `json:"shk_last_hours"`
	PlanCountDeparture int                  `json:"plan_count_departure"`
	Parking            []int                `json:"parking"`
	Suppliers          []*RouteSupplierInfo `json:"suppliers"`
	NormativeInLiters  float32              `json:"normative_liters"`
	VolumeMlByContent  int                  `json:"volume_ml_by_content"`
}

type RouteSupplierInfo struct {
	ID             int    `json:"supplier_id"`
	Name           string `json:"supplier_name"`
	CountDeparture int    `json:"count_departure"`
}

type RemainsLastMileReportsRouteInfo struct {
	CountBarcodes         int                                    `json:"count_shk"`
	CountTare             int                                    `json:"count_tare"`
	DestinationOfficeID   int                                    `json:"dst_office_id"`
	DestinationOfficeName string                                 `json:"dst_office_name"`
	TotalVolumeMl         int                                    `json:"total_volume_ml"`
	Tares                 []*RemainsLastMileReportsRouteInfoTare `json:"tares"`
}

type RemainsLastMileReportsRouteInfoTare struct {
	ID            int       `json:"tare_id"`
	PrepareDate   time.Time `json:"prepare_dt"`
	CountBarcodes int       `json:"count_shk"`
	SpID          int       `json:"sp_id"`
	SpName        string    `json:"sp_name"`
	VolumeMl      int       `json:"volume_ml_by_content"`
}
