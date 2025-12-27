package models

import "time"

type ShipmentsMeta struct {
	TotalCount int `json:"total_count"`
}

type Shipment struct {
	ShipmentID         int       `json:"shipment_id"`
	CloseDt            time.Time `json:"close_dt"`
	CreateDt           time.Time `json:"create_dt"`
	RouteCarId         int       `json:"route_car_id"`
	SupplierId         int       `json:"supplier_id"`
	SupplierName       string    `json:"supplier_name"`
	SrcOfficeId        int       `json:"src_office_id"`
	OfficeName         string    `json:"office_name"`
	DriverName         string    `json:"driver_name"`
	DstOfficeName      []string  `json:"dst_office_name"`
	VehicleNumberPlate string    `json:"vehicle_number_plate"`
	IsLastMile         bool      `json:"is_last_mile"`
}

type ShipmentInfo struct {
	ID                         int                                  `json:"id"`
	CreateDt                   time.Time                            `json:"create_dt"`
	CloseDt                    time.Time                            `json:"close_dt"`
	RouteID                    int                                  `json:"route_id"`
	RouteTypeID                int                                  `json:"route_type_id"`
	WaySheetID                 int                                  `json:"way_sheet_id"`
	CountryID                  int                                  `json:"country_id"`
	RpoID                      int                                  `json:"rpo_id"`
	DriverFreelancerID         int                                  `json:"driver_freelancer_id"`
	DriverID                   int                                  `json:"driver_id"`
	DriverName                 string                               `json:"driver_name"`
	IsLastMile                 bool                                 `json:"is_last_mile"`
	IsReturnShipment           bool                                 `json:"is_return_shipment"`
	ResponsibleUserAccountID   int                                  `json:"responsible_user_account_id"`
	ResponsibleUserAccountName string                               `json:"responsible_user_account_name"`
	SrcOfficeID                int                                  `json:"src_office_id"`
	SrcOfficeName              string                               `json:"src_office_name"`
	SupplierID                 int                                  `json:"supplier_id"`
	SupplierName               string                               `json:"supplier_name"`
	TowedVehicleID             int                                  `json:"towed_vehicle_id"`
	TowedVehicleName           string                               `json:"towed_vehicle_name"`
	TowedVehicleNumberPlate    string                               `json:"towed_vehicle_number_plate"`
	VehicleID                  int                                  `json:"vehicle_id"`
	VehicleName                string                               `json:"vehicle_name"`
	VehicleNumberPlate         string                               `json:"vehicle_number_plate"`
	DestinationOfficesInfo     []*ShipmentInfoDestinationOfficeInfo `json:"ttns"`
	PointType                  int                                  `json:"point_type"`
	Seal                       interface{}                          `json:"seal"`   // ???
	Blrspt                     interface{}                          `json:"blrspt"` // ???
}

type ShipmentInfoDestinationOfficeInfo struct {
	ID             int    `json:"id"`
	CountryID      int    `json:"country_id"`
	DstOfficeID    int    `json:"dst_office_id"`
	DstOfficeName  string `json:"dst_office_name"`
	PointType      int    `json:"point_type"`
	ResponsibleFIO string `json:"responsible_fio"`
}

type ShipmentTransfers struct {
	Containers    []*ShipmentContainer   `json:"containers"`
	TransferBoxes []*ShipmentTransferBox `json:"transfer_boxes"`
}

type ShipmentTransferBox struct {
	BoxPrice            int       `json:"box_price"`
	BoxType             int       `json:"box_type"`
	CreateDt            time.Time `json:"create_dt"`
	DstOfficeID         int       `json:"dst_office_id"`
	DstOfficeName       string    `json:"dst_office_name"`
	DtArrival           time.Time `json:"dt_arrival"`
	LastOperationDt     time.Time `json:"last_operation_dt"`
	PaymentSum          int       `json:"payment_sum"`
	PlaceName           string    `json:"place_name"`
	ResponsibleUserName string    `json:"responsible_user_name"`
	CountBarcodes       int       `json:"shks_count"`
	SrcOfficeName       string    `json:"src_office_name"`
	StatusID            int       `json:"status_id"`
	TtnID               int       `json:"ttn_id"`
	UnloadUserName      string    `json:"unload_user_name"`
	VolumeMl            int       `json:"volume_ml"`
	WeightMg            int       `json:"weight_mg"`
}

type ShipmentContainer struct {
	// ???
}

// Params request

type GetShipmentParamsRequest struct {
	DataStart                time.Time   `json:"dt_start"`
	DataEnd                  time.Time   `json:"dt_end"`
	SrcOfficeID              int         `json:"src_office_id"`
	PageIndex                int         `json:"page_index"`
	Limit                    int         `json:"limit"`
	SupplierID               int         `json:"supplier_id"`
	Direction                int         `json:"direction"`
	Sorter                   string      `json:"sorter"` // Types: "id", "supplier", "driver", "route_car_id", "updated_at"
	FilterShipmentID         int         `json:"shipment_id"`
	FilterShipmentType       string      `json:"shipment_type"` // Types: "last-mile" or "truck"
	FilterVehicleNumberPlate int         `json:"vehicle_number_plate"`
	FilterDstOffice          interface{} `json:"dst_office"` // id(int) or name(string)
	FilterShowOnlyOpen       bool        `json:"show_only_open"`
}
