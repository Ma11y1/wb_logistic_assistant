package reports

import (
	"fmt"
	"slices"
	"strconv"
	"time"
	"wb_logistic_assistant/internal/errors"
)

type ShipmentCloseReportData struct {
	RouteID                  int
	ShipmentID               int
	WaySheetID               int
	Date                     time.Time
	DateCreate               time.Time
	DateClose                time.Time
	Parking                  int
	DriverName               string
	VehicleNumberPlate       string
	BarcodesTotalRemains     int
	BarcodesTotalTransfer    int
	BarcodesStandard         float64
	BarcodesDeviationPercent float64
	TareTotalRemains         int
	TareTotalTransfer        int
	SpName                   string
	RemainsTaresInfo         []*ShipmentCloseRemainsTareInfo
}

type ShipmentCloseRemainsTareInfo struct {
	ID              int
	DstOfficeID     int
	DstOfficeName   string
	LastOperationDt time.Time
	CountBarcodes   int
}

type ShipmentCloseReport struct{}

func (r *ShipmentCloseReport) Render(data *ShipmentCloseReportData) (*ReportData, error) {
	if data == nil {
		return nil, errors.New("ShipmentCloseReport.Render()", "data is empty")
	}

	report := NewReportData()

	report.Header = &Item{
		Children: []*Item{
			{Text: time.Now().Format("02.01.2006 15:04 -07"), Quote: true},
		},
	}

	report.Body = &Item{
		Children: []*Item{
			{Text: "Маршрут:", Bold: true, Block: true}, {Text: strconv.Itoa(data.RouteID)},
			{Text: "Парковка:", Bold: true, Block: true}, {Text: strconv.Itoa(data.Parking)},
			{Text: "Отгрузка:", Bold: true, Block: true}, {Text: strconv.Itoa(data.ShipmentID),
				Link: "https://logistics.wildberries.ru/external-logistics/shipments-shell/shipments/" + strconv.Itoa(data.ShipmentID)},
			{Text: "Путевой лист:", Bold: true, Block: true}, {Text: strconv.Itoa(data.WaySheetID),
				Link: "https://ol.wildberries.ru/#/layout/external-waysheet/" + strconv.Itoa(data.WaySheetID)},
			{Text: "Дата:", Bold: true, Block: true}, {Text: data.Date.Format("02.01.2006")},
			{Text: "Открытие:", Bold: true, Block: true}, {Text: data.DateCreate.Format("15:04")},
			{Text: "Закрытие:", Bold: true, Block: true}, {Text: data.DateClose.Format("15:04")},
			{Text: "Водитель:", Bold: true, Block: true}, {Text: data.DriverName},
			{Text: "Автомобиль:", Bold: true, Block: true}, {Text: data.VehicleNumberPlate},
			{Text: "ШК остаток:", Bold: true, Block: true}, {Text: strconv.Itoa(data.BarcodesTotalRemains)},
			{Text: "ШК отгружено:", Bold: true, Block: true}, {Text: strconv.Itoa(data.BarcodesTotalTransfer)},
			{Text: "ШК норматив:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.0f", data.BarcodesStandard)},
			{Text: "ШК отклонение:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.1f%%", data.BarcodesDeviationPercent)},
			{Text: "Тара остаток:", Bold: true, Block: true}, {Text: strconv.Itoa(data.TareTotalRemains)},
			{Text: "Тара отгружено:", Bold: true, Block: true}, {Text: strconv.Itoa(data.TareTotalTransfer)},
			{Text: "МХ:", Bold: true, Block: true}, {Text: data.SpName},
			{Block: true},
		},
	}

	slices.SortFunc(data.RemainsTaresInfo, func(a, b *ShipmentCloseRemainsTareInfo) int {
		if a.LastOperationDt.Equal(b.LastOperationDt) {
			return 0
		}
		if a.LastOperationDt.After(b.LastOperationDt) {
			return 1
		}
		return -1
	})

	if len(data.RemainsTaresInfo) > 0 {
		tares := &Item{
			HiddenQuote: true,
			Block:       true,
			Children:    make([]*Item, 0, len(data.RemainsTaresInfo)*5),
		}

		for i, info := range data.RemainsTaresInfo {
			if info == nil {
				continue
			}

			tares.Children = append(tares.Children,
				&Item{Text: strconv.Itoa(i+1) + ") " +
					strconv.Itoa(info.ID) + " " +
					strconv.Itoa(info.DstOfficeID) + " " +
					info.DstOfficeName,
					Block: true},
				&Item{Text: "ШК:", Bold: true, Block: true}, &Item{Text: strconv.Itoa(info.CountBarcodes)},
				&Item{Text: "Дата операции:", Bold: true}, &Item{Text: info.LastOperationDt.Format("02.01.2006 15:04")},
			)
		}

		report.Body.AddChild(tares)
	}

	return report, nil
}
