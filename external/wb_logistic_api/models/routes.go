package models

import "time"

type AssociationRouteInfoByName struct {
	ID                   int                                   `json:"id"`
	SrcOfficeName        string                                `json:"srcOfficeName"`
	Distance             int                                   `json:"distance"`
	TotalCost            float64                               `json:"totalCost"`
	UnloadBoxCost        float64                               `json:"unloadBoxCost"`
	Price                float64                               `json:"price"`
	Date                 time.Time                             `json:"date"`
	Code1C               int                                   `json:"code1C"`
	IsTechnical          bool                                  `json:"isTechnical"`
	IsDeleted            bool                                  `json:"isDeleted"`
	IsExternal           bool                                  `json:"isExternal"`
	IsTransit            bool                                  `json:"isTransit"`
	IsCascade            bool                                  `json:"isCascade"`
	GeoLine              interface{}                           `json:"geoLine"`
	Color                interface{}                           `json:"color"`
	OfficeName           string                                `json:"officeName"`
	Type                 int                                   `json:"type"`
	Details              interface{}                           `json:"details"`
	Suppliers            []*AssociationRouteInfoByNameSupplier `json:"suppliers"`
	PaymentType          int                                   `json:"paymentType"`
	PaymentData          string                                `json:"paymentData"`
	RpoID                int                                   `json:"rpoId"`
	CostItemID           int                                   `json:"costItemId"`
	CostItemName         string                                `json:"costItemName"`
	HasMultipleCars      bool                                  `json:"hasMultipleCars"`
	TimeTravel           time.Time                             `json:"timeTravel"`
	DepartureTime        time.Time                             `json:"departureTime"`
	ArrivalTime          time.Time                             `json:"arrivalTime"`
	IsCircular           bool                                  `json:"isCircular"`
	Workload             int                                   `json:"workload"`
	Periodicity          int                                   `json:"periodicity"`
	PointsCount          int                                   `json:"pointsCount"`
	NeedsApproval        bool                                  `json:"needsApproval"`
	Comment              string                                `json:"comment"`
	ApprovedEmployeeID   int                                   `json:"approvedEmployeeId"`
	ApprovedEmployeeName string                                `json:"approvedEmployeeName"`
	IsBalance            bool                                  `json:"isBalance"`
	FuelRate             float64                               `json:"fuelRate"`
	DcOffices            int                                   `json:"dcOffices"`
	ChainedRoute         int                                   `json:"chainedRoute"`
	IsBackRoute          bool                                  `json:"isBackRoute"`
	CalcID               int                                   `json:"calcId"`
	IsPriceFixed         bool                                  `json:"isPriceFixed"`
	TimeDepartures       time.Time                             `json:"timeDepartures"`
	DeparturesCount      int                                   `json:"departuresCount"`
	IsInternational      bool                                  `json:"isInternational"`
	ContainerCount       int                                   `json:"containerCount"`
	VehicleBodyVolume    float64                               `json:"vehicleBodyVolume"`
	IsReturn             bool                                  `json:"isReturn"`
	WorkloadWBDrive      interface{}                           `json:"workloadWBDrive"`
	WorkloadWBGo         interface{}                           `json:"workloadWBGo"`
	WorkloadTitle        string                                `json:"workloadTitle"`
	TrunkRoadPrice       int                                   `json:"trunkRoadPrice"`
	IsCluster            bool                                  `json:"isCluster"`
	Name                 string                                `json:"name"`
	SrcOfficeID          int                                   `json:"srcOfficeId"`
	Sequence             int                                   `json:"sequence"`
	IsWarehouse          bool                                  `json:"isWarehouse"`
}

type AssociationRouteInfoByNameSupplier struct {
	ID           int    `json:"suplierId"` // typo on wb service's side
	Name         string `json:"supplierName"`
	RouteCarID   int    `json:"routeCarId"`
	RouteCarName string `json:"routeCarName"`
}
