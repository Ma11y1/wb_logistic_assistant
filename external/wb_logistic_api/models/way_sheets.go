package models

import "time"

type WaySheetsPage struct {
	TotalWaySheets int         `json:"total_waysheets"`
	TotalPrice     float64     `json:"total_price"`
	TotalFine      float64     `json:"total_fine"`
	Page           int         `json:"page"`
	Pages          int         `json:"pages"`
	TotalReturn    float64     `json:"total_return"`
	WaySheets      []*WaySheet `json:"waysheets"`
}

type WaySheet struct {
	WaySheetID            string    `json:"way_sheet_id"`
	WayTypeID             int       `json:"way_type_id"`
	OpenDt                time.Time `json:"open_dt"`
	CloseDt               time.Time `json:"close_dt"`
	SrcOfficeID           string    `json:"src_office_id"`
	RouteCarID            string    `json:"route_car_id"`
	RouteCarName          string    `json:"route_car_name"`
	DriverID              string    `json:"driver_id"`
	DriverName            string    `json:"driver_name"`
	VehicleNumberPlate    string    `json:"vehicle_number_plate"`
	SupplierID            string    `json:"supplier_id"`
	SupplierName          string    `json:"supplier_name"`
	PaymentType           string    `json:"payment_type"`
	CountBox              string    `json:"count_box"`
	CountContainer        string    `json:"count_container"`
	CountBarcodes         string    `json:"count_shk"`
	TotalPrice            float64   `json:"total_price"`
	SumFine               float64   `json:"sum_fine"`
	PlanMileage           string    `json:"plan_mileage"`
	SrcOfficeName         string    `json:"src_office_name"`
	CountPalletes         string    `json:"count_palletes"`
	SumReturn             float64   `json:"sum_return"`
	CountArrivalBox       string    `json:"count_arrival_box"`
	CountArrivalContainer string    `json:"count_arrival_container"`
	CountArrivalPalletes  string    `json:"count_arrival_palletes"`
	IsArrival             bool      `json:"is_arrival"`
}

type WaySheetInfo struct {
	ID                 string                       `json:"way_sheet_id"`
	WayTypeID          int                          `json:"way_type_id"`
	TotalBarcodesCount int                          `json:"total_shk_count"`
	TotalVolumeCount   int                          `json:"total_volume_count"`
	DateOpen           time.Time                    `json:"date_open"`
	DateClose          time.Time                    `json:"date_close"`
	PlanMileage        string                       `json:"plan_mileage"`
	SrcOffice          *WaySheetSourceOffice        `json:"src_office"`
	DstOffices         []*WaySheetDestinationOffice `json:"dst_offices"`
	Tares              []*WaySheetTare              `json:"tares"`
	Containers         []interface{}                `json:"containers"` // ???
	Pallets            []interface{}                `json:"pallets"`    // ???
	Route              *WaySheetRoute               `json:"route"`
	Drivers            []*WaySheetDriver            `json:"drivers"`
	Vehicles           *WaySheetVehicles            `json:"vehicles"`
	Shippings          []*WaySheetShipping          `json:"shippings"`
}

type WaySheetRoute struct {
	RouteCarID   string `json:"routecar_id"`
	RouteCarName string `json:"routecar_name"`
}

type WaySheetDriver struct {
	DriverName string `json:"driver_name"`
}

type WaySheetVehicles struct {
	ShippingCarNumber string `json:"shipping_car_number"`
}

type WaySheetSupplier struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type WaySheetShipping struct {
	ID          string    `json:"shipping_id"`
	TtnID       string    `json:"ttn_id"`
	DstOfficeID string    `json:"dst_office_id"`
	CloseDt     time.Time `json:"close_dt"`
}

// params

type GetWaySheetsParamsRequest struct {
	DateClose          time.Time `json:"date_close"`
	DateOpen           time.Time `json:"date_open"`
	SupplierID         int       `json:"supplier_id"`
	SrcOfficeID        int       `json:"src_office_id"`
	RouteCarID         int       `json:"routecar_id"`
	Limit              int       `json:"limit"`
	Offset             int       `json:"offset"`
	WayTypeID          int       `json:"way_type_id"`
	VehicleNumberPlate string    `json:"vehicle_number_plate"`
}
