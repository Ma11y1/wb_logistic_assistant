package reports

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

type GeneralRoutesReportMetaData struct {
	Update                   time.Time
	TimeUpdateChangeBarcodes time.Time
	TimeUpdateRating         time.Time
	TimeUpdateShipments      time.Time
	TimeRemainsBarcodes      time.Time
	TimeUpdateWaySheets      time.Time
}

type GeneralRoutesReportData struct {
	RouteID                      int
	Parking                      int
	Tares                        int
	VolumeLiters                 float32
	VolumeNormativeLiters        float32
	VolumeNormativeLitersPercent float32
	Barcodes                     int
	ChangeBarcodes               int
	Rating                       float32
	RemainsBarcodes              int

	ShipmentID         int
	ShipmentCreateDate time.Time
	ShipmentCloseDate  time.Time

	WaySheetID                   int
	WaySheetDateLastOperation    time.Time
	WaySheetTotalAddresses       int
	WaySheetCurrentAddresses     int
	WaySheetTotalReturnedTares   int
	WaySheetCurrentReturnedTares int

	PrevWaySheetID                   int
	PrevWaySheetDateLastOperation    time.Time
	PrevWaySheetTotalAddresses       int
	PrevWaySheetCurrentAddresses     int
	PrevWaySheetTotalReturnedTares   int
	PrevWaySheetCurrentReturnedTares int

	WaySheetsInterval time.Duration
}

// GeneralRoutesSheetReport Generating report for sheet
type GeneralRoutesSheetReport struct {
	isSort      bool
	isAscending bool
	sortColumn  int

	countHeaderRows int
	rowHeaderTimes  int
	rowHeaderNames  int

	countColumn                          int
	columnRouteID                        int
	columnParking                        int
	columnTares                          int
	columnBarcodes                       int
	columnChangeBarcodes                 int
	columnVolumeLiters                   int
	columnVolumeNormativeLiters          int
	columnRating                         int
	columnShipmentID                     int
	columnShipmentCreateDate             int
	columnShipmentCloseDate              int
	columnRemainsBarcodes                int
	columnWaySheetID                     int
	columnWaySheetDateCloseAddress       int
	columnWaySheetAddresses              int
	columnWaySheetTotalReturnedTares     int
	columnPrevWaySheetID                 int
	columnPrevWaySheetDateCloseAddress   int
	columnPrevWaySheetAddresses          int
	columnPrevWaySheetTotalReturnedTares int
	columnWaySheetInterval               int
}

func NewGeneralRoutesSheetReport(sort bool, sortColumn int, ascending bool) *GeneralRoutesSheetReport {
	return &GeneralRoutesSheetReport{
		isSort:      sort,
		isAscending: ascending,
		sortColumn:  sortColumn,

		countHeaderRows: 2,
		rowHeaderTimes:  0,
		rowHeaderNames:  1,

		countColumn:                          21,
		columnRouteID:                        0,
		columnParking:                        1,
		columnTares:                          2,
		columnVolumeLiters:                   3,
		columnVolumeNormativeLiters:          4,
		columnBarcodes:                       5,
		columnChangeBarcodes:                 6,
		columnRating:                         7,
		columnWaySheetID:                     8,
		columnShipmentID:                     9,
		columnWaySheetAddresses:              10,
		columnWaySheetDateCloseAddress:       11,
		columnShipmentCreateDate:             12,
		columnShipmentCloseDate:              13,
		columnWaySheetTotalReturnedTares:     14,
		columnRemainsBarcodes:                15,
		columnPrevWaySheetID:                 16,
		columnPrevWaySheetAddresses:          17,
		columnPrevWaySheetDateCloseAddress:   18,
		columnPrevWaySheetTotalReturnedTares: 19,
		columnWaySheetInterval:               20,
	}
}

func (r *GeneralRoutesSheetReport) Render(meta *GeneralRoutesReportMetaData, routes []*GeneralRoutesReportData) (*ReportData, error) {
	report := &ReportData{
		Header: &Item{Children: make([]*Item, r.countHeaderRows)},
		Body:   &Item{Children: make([]*Item, len(routes))},
	}

	report.Header.Children[r.rowHeaderTimes] = &Item{Children: make([]*Item, r.countColumn), Block: true}
	report.Header.Children[r.rowHeaderNames] = &Item{Children: make([]*Item, r.countColumn), Block: true}

	headerTimes := report.Header.Children[r.rowHeaderTimes].Children
	headerTimes[r.columnRouteID] = &Item{Text: meta.Update.Format("15:04 - 02.01.2006")}
	headerTimes[r.columnChangeBarcodes] = &Item{Text: meta.TimeUpdateChangeBarcodes.Format("15:04")}
	headerTimes[r.columnRating] = &Item{Text: meta.TimeUpdateRating.Format("15:04")}
	headerTimes[r.columnShipmentID] = &Item{Text: meta.TimeUpdateShipments.Format("15:04")}
	headerTimes[r.columnRemainsBarcodes] = &Item{Text: meta.TimeRemainsBarcodes.Format("15:04")}
	headerTimes[r.columnWaySheetID] = &Item{Text: meta.TimeUpdateWaySheets.Format("15:04")}
	headerTimes[r.columnPrevWaySheetID] = &Item{Text: meta.TimeUpdateWaySheets.Format("15:04")}

	headerNames := report.Header.Children[r.rowHeaderNames].Children
	headerNames[r.columnRouteID] = &Item{Text: "Маршрут"}
	headerNames[r.columnParking] = &Item{Text: "Парковка"}
	headerNames[r.columnTares] = &Item{Text: "Тара"}
	headerNames[r.columnVolumeLiters] = &Item{Text: "Объем, л"}
	headerNames[r.columnVolumeNormativeLiters] = &Item{Text: "Норматив, % (л)"}
	headerNames[r.columnBarcodes] = &Item{Text: "ШК"}
	headerNames[r.columnChangeBarcodes] = &Item{Text: "ШК, изменение"}
	headerNames[r.columnRating] = &Item{Text: "Рейтинг"}
	headerNames[r.columnWaySheetID] = &Item{Text: "Путевой лист 1"}
	headerNames[r.columnWaySheetAddresses] = &Item{Text: "Адреса"}
	headerNames[r.columnWaySheetDateCloseAddress] = &Item{Text: "Время сдачи"}
	headerNames[r.columnWaySheetTotalReturnedTares] = &Item{Text: "Возвраты"}
	headerNames[r.columnShipmentID] = &Item{Text: "Отгрузка"}
	headerNames[r.columnShipmentCreateDate] = &Item{Text: "Открыта"}
	headerNames[r.columnShipmentCloseDate] = &Item{Text: "Закрыта"}
	headerNames[r.columnRemainsBarcodes] = &Item{Text: "Остаток"}
	headerNames[r.columnPrevWaySheetID] = &Item{Text: "Путевой лист 2"}
	headerNames[r.columnPrevWaySheetAddresses] = &Item{Text: "Адреса"}
	headerNames[r.columnPrevWaySheetDateCloseAddress] = &Item{Text: "Время сдачи"}
	headerNames[r.columnPrevWaySheetTotalReturnedTares] = &Item{Text: "Возвраты"}
	headerNames[r.columnWaySheetInterval] = &Item{Text: "Промежуток"}

	if r.isSort {
		sort.Slice(routes, func(i, j int) bool {
			switch r.sortColumn {
			case r.columnRouteID:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].RouteID < routes[j].RouteID
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].RouteID > routes[j].RouteID
					})
				}

			case r.columnParking:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].Parking < routes[j].Parking
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].Parking > routes[j].Parking
					})
				}

			case r.columnTares:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].Tares < routes[j].Tares
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].Tares > routes[j].Tares
					})
				}

			case r.columnBarcodes:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].Barcodes < routes[j].Barcodes
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].Barcodes > routes[j].Barcodes
					})
				}

			case r.columnChangeBarcodes:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].ChangeBarcodes < routes[j].ChangeBarcodes
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].ChangeBarcodes > routes[j].ChangeBarcodes
					})
				}

			case r.columnVolumeLiters:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].VolumeLiters < routes[j].VolumeLiters
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].VolumeLiters > routes[j].VolumeLiters
					})
				}

			case r.columnVolumeNormativeLiters:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].VolumeNormativeLiters < routes[j].VolumeNormativeLiters
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].VolumeNormativeLiters > routes[j].VolumeNormativeLiters
					})
				}

			case r.columnRating:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].Rating < routes[j].Rating
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].Rating > routes[j].Rating
					})
				}

			case r.columnShipmentID:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].ShipmentID < routes[j].ShipmentID
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].ShipmentID > routes[j].ShipmentID
					})
				}

			case r.columnShipmentCreateDate:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].ShipmentCreateDate.Before(routes[j].ShipmentCreateDate)
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].ShipmentCreateDate.After(routes[j].ShipmentCreateDate)
					})
				}

			case r.columnShipmentCloseDate:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].ShipmentCloseDate.Before(routes[j].ShipmentCloseDate)
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].ShipmentCloseDate.After(routes[j].ShipmentCloseDate)
					})
				}

			case r.columnRemainsBarcodes:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].RemainsBarcodes < routes[j].RemainsBarcodes
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].RemainsBarcodes > routes[j].RemainsBarcodes
					})
				}

			case r.columnWaySheetID:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetID < routes[j].WaySheetID
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetID > routes[j].WaySheetID
					})
				}

			case r.columnWaySheetDateCloseAddress:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetDateLastOperation.Before(routes[j].WaySheetDateLastOperation)
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetDateLastOperation.After(routes[j].WaySheetDateLastOperation)
					})
				}

			case r.columnWaySheetAddresses:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetTotalAddresses < routes[j].WaySheetTotalAddresses
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetTotalAddresses > routes[j].WaySheetTotalAddresses
					})
				}

			case r.columnWaySheetTotalReturnedTares:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetTotalReturnedTares < routes[j].WaySheetTotalReturnedTares
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetTotalReturnedTares > routes[j].WaySheetTotalReturnedTares
					})
				}

			case r.columnPrevWaySheetID:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].PrevWaySheetID < routes[j].PrevWaySheetID
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].PrevWaySheetID > routes[j].PrevWaySheetID
					})
				}

			case r.columnPrevWaySheetDateCloseAddress:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].PrevWaySheetDateLastOperation.Before(routes[j].PrevWaySheetDateLastOperation)
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].PrevWaySheetDateLastOperation.After(routes[j].PrevWaySheetDateLastOperation)
					})
				}

			case r.columnPrevWaySheetAddresses:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].PrevWaySheetTotalAddresses < routes[j].PrevWaySheetTotalAddresses
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].PrevWaySheetTotalAddresses > routes[j].PrevWaySheetTotalAddresses
					})
				}

			case r.columnPrevWaySheetTotalReturnedTares:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].PrevWaySheetTotalReturnedTares < routes[j].PrevWaySheetTotalReturnedTares
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].PrevWaySheetTotalReturnedTares > routes[j].PrevWaySheetTotalReturnedTares
					})
				}

			case r.columnWaySheetInterval:
				if r.isAscending {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetsInterval < routes[j].WaySheetsInterval
					})
				} else {
					sort.Slice(routes, func(i, j int) bool {
						return routes[i].WaySheetsInterval > routes[j].WaySheetsInterval
					})
				}

			}
			return false
		})
	}

	body := report.Body.Children
	for i, route := range routes {
		if route == nil {
			continue
		}
		if i >= len(body) {
			body = append(body, &Item{})
			report.Body.Children = body
		}

		body[i] = &Item{Children: make([]*Item, r.countColumn), Block: true}
		row := body[i].Children

		row[r.columnRouteID] = &Item{Text: itoa(route.RouteID)}
		if route.Parking == 0 {
			row[r.columnParking] = &Item{Text: " "}
		} else {
			row[r.columnParking] = &Item{Text: itoa(route.Parking)}
		}
		row[r.columnTares] = &Item{Text: itoa(route.Tares)}
		row[r.columnBarcodes] = &Item{Text: itoa(route.Barcodes)}
		if route.ChangeBarcodes == 0 {
			row[r.columnChangeBarcodes] = &Item{Text: " "}
		} else {
			row[r.columnChangeBarcodes] = &Item{Text: itoa(route.ChangeBarcodes)}
		}
		row[r.columnVolumeLiters] = &Item{Text: fmt.Sprintf("%.1f", route.VolumeLiters)}
		row[r.columnVolumeNormativeLiters] = &Item{Text: fmt.Sprintf("%.1f%%, (%.1f)", route.VolumeNormativeLitersPercent, route.VolumeNormativeLiters)}

		if route.Rating > 0 {
			row[r.columnRating] = &Item{Text: fmt.Sprintf("%.1f", route.Rating)}
		} else {
			row[r.columnRating] = &Item{Text: " "}
		}

		if route.ShipmentID > 0 {
			row[r.columnShipmentID] = &Item{Text: itoa(route.ShipmentID), Link: "https://logistics.wildberries.ru/external-logistics/shipments-shell/shipments/" + strconv.Itoa(route.ShipmentID)}
		} else {
			row[r.columnShipmentID] = &Item{Text: " "}
		}

		if route.ShipmentCreateDate.IsZero() {
			row[r.columnShipmentCreateDate] = &Item{Text: " "}
		} else {
			row[r.columnShipmentCreateDate] = &Item{Text: route.ShipmentCreateDate.Format("15:04 (02.01.06)")}
		}
		if route.ShipmentCloseDate.IsZero() {
			if !route.ShipmentCreateDate.IsZero() {
				row[r.columnShipmentCloseDate] = &Item{Text: "Открыто"}
			} else {
				row[r.columnShipmentCloseDate] = &Item{Text: " "}
			}
		} else {
			row[r.columnShipmentCloseDate] = &Item{Text: route.ShipmentCloseDate.Format("15:04 (02.01.06)")}
		}
		if route.RemainsBarcodes == 0 {
			row[r.columnRemainsBarcodes] = &Item{Text: " "}
		} else {
			row[r.columnRemainsBarcodes] = &Item{Text: itoa(route.RemainsBarcodes)}
		}

		if route.WaySheetID > 0 {
			row[r.columnWaySheetID] = &Item{Text: itoa(route.WaySheetID), Link: "https://ol.wildberries.ru/#/layout/external-waysheet/" + strconv.Itoa(route.WaySheetID)}
		} else {
			row[r.columnWaySheetID] = &Item{Text: " "}
		}
		if route.WaySheetDateLastOperation.IsZero() {
			row[r.columnWaySheetDateCloseAddress] = &Item{Text: " "}
		} else {
			row[r.columnWaySheetDateCloseAddress] = &Item{Text: route.WaySheetDateLastOperation.Format("15:04 (02.01.06)")}
		}
		if route.WaySheetTotalAddresses == 0 {
			row[r.columnWaySheetAddresses] = &Item{Text: " "}
		} else {
			row[r.columnWaySheetAddresses] = &Item{Text: fmt.Sprintf("%d/%d", route.WaySheetCurrentAddresses, route.WaySheetTotalAddresses)}
		}
		if route.WaySheetTotalReturnedTares == 0 {
			row[r.columnWaySheetTotalReturnedTares] = &Item{Text: " "}
		} else {
			row[r.columnWaySheetTotalReturnedTares] = &Item{Text: fmt.Sprintf("%d/%d", route.WaySheetCurrentReturnedTares, route.WaySheetTotalReturnedTares)}
		}

		if route.PrevWaySheetID > 0 {
			row[r.columnPrevWaySheetID] = &Item{Text: itoa(route.PrevWaySheetID), Link: "https://ol.wildberries.ru/#/layout/external-waysheet/" + strconv.Itoa(route.PrevWaySheetID)}
		} else {
			row[r.columnPrevWaySheetID] = &Item{Text: " "}
		}
		if route.PrevWaySheetDateLastOperation.IsZero() {
			row[r.columnPrevWaySheetDateCloseAddress] = &Item{Text: " "}
		} else {
			row[r.columnPrevWaySheetDateCloseAddress] = &Item{Text: route.PrevWaySheetDateLastOperation.Format("15:04 (02.01.06)")}
		}
		if route.PrevWaySheetTotalAddresses == 0 {
			row[r.columnPrevWaySheetAddresses] = &Item{Text: " "}
		} else {
			row[r.columnPrevWaySheetAddresses] = &Item{Text: fmt.Sprintf("%d/%d", route.PrevWaySheetTotalAddresses, route.PrevWaySheetCurrentAddresses)}
		}
		if route.PrevWaySheetTotalReturnedTares == 0 {
			row[r.columnPrevWaySheetTotalReturnedTares] = &Item{Text: " "}
		} else {
			row[r.columnPrevWaySheetTotalReturnedTares] = &Item{Text: fmt.Sprintf("%d/%d", route.PrevWaySheetCurrentReturnedTares, route.PrevWaySheetTotalReturnedTares)}
		}

		if route.WaySheetsInterval == 0 {
			row[r.columnWaySheetInterval] = &Item{Text: " "}
		} else {
			row[r.columnWaySheetInterval] = &Item{Text: fmt.Sprintf("%02d:%02d", int(route.WaySheetsInterval.Hours()), int(route.WaySheetsInterval.Minutes())%60)}
		}
	}

	return report, nil
}
