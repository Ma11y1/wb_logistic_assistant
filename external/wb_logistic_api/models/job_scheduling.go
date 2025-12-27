package models

type JobsScheduling struct {
	ActiveTask          int                    `json:"active_task"`
	PlanTask            int                    `json:"plan_task"`
	CountAutoAccess     int                    `json:"count_auto_access"`
	CountRouteAccess    int                    `json:"count_route_access"`
	AutoPassApplication []interface{}          `json:"auto_pass_application"` // empty? not clear from regular queries
	VehicleFree         []interface{}          `json:"vehicle_free"`          // empty? not clear from regular queries
	Route               []*JobsSchedulingRoute `json:"route"`
}

type JobsSchedulingRoute struct {
	RouteID          int           `json:"route_id"`
	RouteName        string        `json:"route_name"`
	Price            int           `json:"price"`
	FactVolumeTare   string        `json:"fact_volume_tare"`
	PlanVolumeTare   string        `json:"plan_volume_tare"`
	Rating           *RouteRating  `json:"rating"`
	Task             []interface{} `json:"task"` // empty? not clear from regular queries
	SrcOfficeId      int           `json:"src_office_id"`
	SrcOfficeName    string        `json:"src_office_name"`
	CrsRequestStatus string        `json:"crs_request_status"`
	TransferDt       int           `json:"transfer_dt"` // empty? not clear from regular queries
	SupplierId       string        `json:"supplier_id"`
	SupplierName     string        `json:"supplier_name"`
	NewSupplierId    string        `json:"new_supplier_id"`
	NewSupplierName  string        `json:"new_supplier_name"`
	ProcessingUntil  interface{}   `json:"processing_until"` // empty? not clear from regular queries
}

type RouteRating struct {
	OverallRating     float64 `json:"overall_rating"`
	BufferSpeedRating float64 `json:"buffer_speed_rating"`
	RoadSpeedRating   float64 `json:"road_speed_rating"`
	BrakRating        float64 `json:"brak_rating"`
	PretensionsRating float64 `json:"pretensions_rating"`
	ActiveDaysRating  float64 `json:"active_days_rating"`
	NoReturnRating    float64 `json:"no_return_rating"`
	AuthorizedRating  float64 `json:"authorized_rating"`
}
