package reports

import (
	"fmt"
	"time"
	"wb_logistic_assistant/internal/errors"
)

type FinanceDailyGeneralReportData struct {
	DateStart          time.Time
	DateEnd            time.Time
	Flights            int
	OpenedFlights      int
	ShippedBarcodes    int
	Tare               int
	ShippedTare        int
	ReturnedTare       int
	Income             float64
	IncomeReturn       float64
	Fine               float64
	SalaryRate         float64
	ExtendedSalaryRate float64
	Marriage           float64
	PercentMarriage    float64
	Tax                float64
	PercentTax         float64
	Margin             float64
	Expenses           float64
	TotalMargin        float64
	OpenedWaySheets    []string
}

type FinanceDailyGeneralReport struct{}

func (r *FinanceDailyGeneralReport) Render(data *FinanceDailyGeneralReportData) (*ReportData, error) {
	if data == nil {
		return nil, errors.New("FinanceDailyGeneralReport.Render()", "data is empty")
	}

	report := NewReportData()

	report.Header = &Item{
		Children: []*Item{
			{Text: time.Now().Format("02.01.2006 15:04 -07"), Quote: true},
		},
	}

	report.Body = &Item{
		Children: []*Item{
			{Text: "ОБЩИЕ РЕЗУЛЬТАТЫ", Bold: true, Block: true},
			{Block: true},
			{Text: "Начало:", Bold: true, Block: true}, {Text: data.DateStart.Format("02.01.2006 15:04")},
			{Text: "Конец:", Bold: true, Block: true}, {Text: data.DateEnd.Format("02.01.2006 15:04")},
			{Text: "Рейсы:", Bold: true, Block: true}, {Text: itoa(data.Flights)},
			{Text: "Незавершенные рейсы:", Bold: true, Block: true}, {Text: itoa(data.OpenedFlights)},
			{Text: "ШК:", Bold: true, Block: true}, {Text: itoa(data.ShippedBarcodes)},
			{Text: "Тара:", Bold: true, Block: true}, {Text: itoa(data.Tare)},
			{Text: "Доставлено тар:", Bold: true, Block: true}, {Text: itoa(data.ShippedTare)},
			{Text: "Возврат тар:", Bold: true, Block: true}, {Text: itoa(data.ReturnedTare)},
			{Text: "Задание:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Income)},
			{Text: "Возврат:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.IncomeReturn)},
			{Text: "Штраф:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Fine)},
			{Text: "Брак:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р. (%.2f%%)", data.Marriage, data.PercentMarriage)},
			{Text: "Налог:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р. (%.2f%%)", data.Tax, data.PercentTax)},
			{Text: "Ставка:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.SalaryRate)},
			{Text: "Ставка+:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.ExtendedSalaryRate)},
			{Text: "Маржа:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Margin)},
			{Text: "Расходы:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Expenses)},
			{Text: "Итого:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.TotalMargin)},
			{Block: true},
		},
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
