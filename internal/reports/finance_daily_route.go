package reports

import (
	"fmt"
	"time"
	"wb_logistic_assistant/internal/errors"
)

type FinanceDailyRouteReportData struct {
	Date               time.Time
	RouteID            int
	Parking            int
	Flights            int
	ShippedBarcodes    int
	Tare               int
	ShippedTare        int
	ReturnedTare       int
	Income             float64
	IncomeReturn       float64
	Fine               float64
	TotalSalaryRate    float64
	SalaryRate         float64
	ExtendedSalaryRate float64
	Marriage           float64
	PercentMarriage    float64
	Tax                float64
	PercentTax         float64
	Margin             float64
	WaySheetIDs        []string
	OpenedWaySheets    []string
}

type FinanceDailyRouteReport struct{}

func (r *FinanceDailyRouteReport) Render(data *FinanceDailyRouteReportData) (*ReportData, error) {
	if data == nil {
		return nil, errors.New("FinanceDailyRouteReport.Render()", "data is empty")
	}

	report := NewReportData()

	report.Header = &Item{
		Children: []*Item{
			{Text: time.Now().Format("02.01.2006 15:04 -07"), Quote: true},
		},
	}

	report.Body = &Item{
		Children: []*Item{
			{Text: "Дата:", Bold: true, Block: true}, {Text: data.Date.Format("02.01.2006")},
			{Text: "Маршрут:", Bold: true, Block: true}, {Text: itoa(data.RouteID)},
			{Text: "Парковка:", Bold: true, Block: true}, {Text: itoa(data.Parking)},
			{Text: "Рейсы:", Bold: true, Block: true}, {Text: itoa(data.Flights)},
			{Text: "ШК:", Bold: true, Block: true}, {Text: itoa(data.ShippedBarcodes)},
			{Text: "Тара:", Bold: true, Block: true}, {Text: itoa(data.Tare)},
			{Text: "Доставлено тар:", Bold: true, Block: true}, {Text: itoa(data.ShippedTare)},
			{Text: "Возврат тар:", Bold: true, Block: true}, {Text: itoa(data.ReturnedTare)},
			{Text: "Задание:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Income)},
			{Text: "Возврат:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.IncomeReturn)},
			{Text: "Штраф:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Fine)},
			{Text: "Брак:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р. (%.2f%%)", data.Marriage, data.PercentMarriage)},
			{Text: "Налог:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р. (%.2f%%)", data.Tax, data.PercentTax)},
			{Text: "Ставка:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р. (%.2f р.)", data.TotalSalaryRate, data.SalaryRate)},
			{Text: "Ставка+:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.ExtendedSalaryRate)},
			{Text: "Маржа:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Margin)},
			{Block: true},
		},
	}

	if len(data.WaySheetIDs) > 0 {
		report.Body.AddChild(&Item{Text: "Путевые листы:", Bold: true, Block: true})

		ids := &Item{
			HiddenQuote: true,
			Block:       true,
			Children:    make([]*Item, 0, len(data.WaySheetIDs)),
		}

		for _, id := range data.WaySheetIDs {
			ids.Children = append(ids.Children, &Item{Text: id, Block: true, Link: "https://ol.wildberries.ru/#/layout/external-waysheet/" + id})
		}

		report.Body.AddChild(ids)
	}

	if len(data.OpenedWaySheets) > 0 {
		report.Body.AddChild(&Item{Text: "Открытые путевые листы:", Bold: true, Block: true})

		ids := &Item{
			HiddenQuote: true,
			Block:       true,
			Children:    make([]*Item, 0, len(data.OpenedWaySheets)),
		}

		for _, id := range data.OpenedWaySheets {
			ids.Children = append(ids.Children, &Item{Text: id, Block: true, Link: "https://ol.wildberries.ru/#/layout/external-waysheet/" + id})
		}

		report.Body.AddChild(ids)
	}

	return report, nil
}
