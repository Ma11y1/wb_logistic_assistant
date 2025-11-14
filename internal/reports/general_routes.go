package reports

import (
	"time"
	"wb_logistic_assistant/internal/errors"
)

type GeneralRoutesReportMetaData struct {
	TimeLastResetCountBarcodes time.Time
}

type GeneralRoutesReportData struct {
	ID                  int
	Parking             int
	CountTares          int
	Barcodes            int
	LastBarcodes        int
	CountChangeBarcodes int
	RemainsBarcodes     int
	Rating              float64
	ShipmentID          int
	ShipmentDateCreate  time.Time
	ShipmentDateClose   time.Time
	ShipmentMH          string
	WaySheetID          int
}

type GeneralRoutesReport struct {
	data                      *ReportData
	indexColumnRouteID        int
	indexColumnParking        int
	indexColumnTares          int
	indexColumnBarcodes       int
	indexColumnChangeBarcodes int
	indexColumnRating         int
	indexShipmentID           int
	indexShipmentDateCreate   int
	indexShipmentDateClose    int
	indexShipmentMH           int
	indexRemainsBarcodes      int
	indexWaySheetID           int
}

func NewGeneralRoutesReport() *GeneralRoutesReport {
	r := &GeneralRoutesReport{
		data:                      NewReportDataBySize(),
		indexColumnRouteID:        0,
		indexColumnParking:        1,
		indexColumnTares:          2,
		indexColumnBarcodes:       3,
		indexColumnChangeBarcodes: 4,
		indexColumnRating:         5,
		indexShipmentID:           6,
		indexShipmentDateCreate:   7,
		indexShipmentDateClose:    8,
		indexShipmentMH:           9,
		indexRemainsBarcodes:      10,
		indexWaySheetID:           11,
	}

	//r.data.Header[1] = []Item{
	//	{Text: "Маршрут", Bold: true},
	//	{Text: "Парковка", Bold: true},
	//	{Text: "Тара", Bold: true},
	//	{Text: "ШК", Bold: true},
	//	{Text: "ШК, изменение", Bold: true},
	//	{Text: "Рейтинг", Bold: true},
	//	{Text: "Отгрузка", Bold: true},
	//	{Text: "Открытие", Bold: true},
	//	{Text: "Закрытие", Bold: true},
	//	{Text: "МХ", Bold: true},
	//	{Text: "Остаток ШК", Bold: true},
	//	{Text: "Путевой лист", Bold: true},
	//}
	return r
}

func (r *GeneralRoutesReport) Render(meta *GeneralRoutesReportMetaData, routesData []*GeneralRoutesReportData) (*ReportData, error) {
	if meta == nil {
		return nil, errors.New("GeneralRoutesReport.Render()", "meta data is empty")
	}
	if routesData == nil || len(routesData) == 0 {
		return nil, errors.New("GeneralRoutesReport.Render()", "routes data is empty")
	}

	//if len(r.data.Header) > 0 && len(r.data.Header[0]) > 0 {
	//	r.data.Header[0][0] = Item{Text: time.Now().Format("15:04 - 02.01.2006")}
	//	r.data.Header[0][r.indexColumnChangeBarcodes] = Item{Text: meta.TimeLastResetCountBarcodes.Format("15:04")}
	//}
	//
	//width := r.data.Width()
	//countRoutes := len(routesData)
	//
	//body := r.data.Body
	//if cap(body) < countRoutes {
	//	body = make([][]Item, countRoutes)
	//} else {
	//	body = body[:countRoutes]
	//}
	//
	//for i := 0; i < countRoutes; i++ {
	//	routeData := routesData[i]
	//	if routeData == nil {
	//		continue
	//	}
	//
	//	row := body[i]
	//	if row == nil || cap(row) < width {
	//		row = make([]Item, width)
	//		body[i] = row
	//	} else {
	//		row = row[:width]
	//	}
	//
	//	row[r.indexColumnRouteID] = Item{Text: itoa(routeData.ID)}
	//	row[r.indexColumnParking] = Item{Text: itoa(routeData.Parking)}
	//	row[r.indexColumnTares] = Item{Text: itoa(routeData.CountTares)}
	//	row[r.indexColumnBarcodes] = Item{Text: itoa(routeData.Barcodes)}
	//	row[r.indexColumnChangeBarcodes] = Item{Text: itoa(routeData.CountChangeBarcodes)}
	//
	//	if routeData.Rating > 0 {
	//		row[r.indexColumnRating] = Item{Text: ftoa(routeData.Rating)}
	//	} else {
	//		row[r.indexColumnRating] = Item{}
	//	}
	//
	//	if routeData.ShipmentID > 0 {
	//		row[r.indexShipmentID] = Item{
	//			Text: itoa(routeData.ShipmentID),
	//			Link: "https://logistics.wildberries.ru/external-logistics/shipments-shell/shipments/" + itoa(routeData.ShipmentID),
	//		}
	//	} else {
	//		row[r.indexShipmentID] = Item{}
	//	}
	//
	//	if !routeData.ShipmentDateCreate.IsZero() {
	//		row[r.indexShipmentDateCreate] = Item{Text: routeData.ShipmentDateCreate.Format("02-01 15:04")}
	//	} else {
	//		row[r.indexShipmentDateCreate] = Item{}
	//	}
	//
	//	if !routeData.ShipmentDateClose.IsZero() {
	//		row[r.indexShipmentDateClose] = Item{Text: routeData.ShipmentDateClose.Format("02-01 15:04")}
	//	} else {
	//		row[r.indexShipmentDateClose] = Item{}
	//	}
	//
	//	row[r.indexShipmentMH] = Item{Text: routeData.ShipmentMH}
	//
	//	if routeData.RemainsBarcodes > 0 {
	//		row[r.indexRemainsBarcodes] = Item{Text: itoa(routeData.RemainsBarcodes)}
	//	} else {
	//		row[r.indexRemainsBarcodes] = Item{}
	//	}
	//
	//	if routeData.WaySheetID > 0 {
	//		row[r.indexWaySheetID] = Item{
	//			Text: itoa(routeData.WaySheetID),
	//			Link: "https://logistics.wildberries.ru/external-logistics/waysheet-registry/" + itoa(routeData.WaySheetID),
	//		}
	//	} else {
	//		row[r.indexWaySheetID] = Item{}
	//	}
	//}
	//r.data.Body = body
	return r.data, nil
}
