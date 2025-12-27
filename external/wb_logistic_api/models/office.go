package models

type AssociationOfficeInfoByName struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsDc      bool    `json:"isDc"`
	IsSc      bool    `json:"isSc"`
	IsWh      bool    `json:"isWh"`
}

type OfficeCoordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type WaySheetSourceOffice struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Coordinates *OfficeCoordinates `json:"coordinates"`
}

type WaySheetDestinationOffice struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Coordinates  *OfficeCoordinates `json:"coordinates"`
	Distance     string             `json:"distance"`
	SumOffice    float64            `json:"sum_office"`
	AvgTime      int                `json:"avg_time"`
	NormTime     int                `json:"norm_time"`
	Sequence     string             `json:"sequence"`
	SequenceFact string             `json:"sequence_fact"`
}
