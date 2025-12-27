package reports

import (
	"fmt"
	"time"
	"wb_logistic_assistant/internal/errors"
)

type FinanceRoutesReportData struct {
	RouteID            int
	ShipmentID         string
	WaySheetID         string
	Parking            int
	DateOpen           time.Time
	DriverName         string
	VehicleNumberPlate string
	ShippedBarcodes    int
	ShippedTare        int
	TotalReturnTare    int
	CurrentReturnTare  int
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
}

type FinanceRoutesReport struct{}

func (r *FinanceRoutesReport) Render(data *FinanceRoutesReportData) (*ReportData, error) {
	if data == nil {
		return nil, errors.New("FinanceRoutesReport.Render()", "data is empty")
	}

	report := NewReportData()

	report.Header = &Item{
		Children: []*Item{
			{Text: time.Now().Format("02.01.2006 15:04 -07"), Quote: true},
		},
	}

	report.Body = &Item{
		Children: []*Item{
			{Text: "Маршрут:", Bold: true, Block: true}, {Text: itoa(data.RouteID)},
			{Text: "Парковка:", Bold: true, Block: true}, {Text: itoa(data.Parking)},
			{Text: "Отгрузка:", Bold: true, Block: true}, {Text: data.ShipmentID,
				Link: "https://logistics.wildberries.ru/external-logistics/shipments-shell/shipments/" + data.ShipmentID},
			{Text: "Путевой лист:", Bold: true, Block: true}, {Text: data.WaySheetID,
				Link: "https://logistics.wildberries.ru/external-logistics/waysheet-registry/LML/" + data.WaySheetID},
			{Text: "Дата погрузки:", Bold: true, Block: true}, {Text: data.DateOpen.Format("02.01.2006")},
			{Text: "Время открытия:", Bold: true, Block: true}, {Text: data.DateOpen.Format("15:04")},
			{Text: "Водитель:", Bold: true, Block: true}, {Text: data.DriverName},
			{Text: "Автомобиль:", Bold: true, Block: true}, {Text: data.VehicleNumberPlate},
			{Text: "Отгружено ШК:", Bold: true, Block: true}, {Text: itoa(data.ShippedBarcodes)},
			{Text: "Отгружено тар:", Bold: true, Block: true}, {Text: itoa(data.ShippedTare)},
			{Text: "Возвраты:", Bold: true, Block: true}, {Text: fmt.Sprintf("%d/%d", data.CurrentReturnTare, data.TotalReturnTare)},
			{Text: "Сумма задания:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Income)},
			{Text: "Сумма возвратов:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.IncomeReturn)},
			{Text: "Сумма штрафов:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Fine)},
			{Text: "Налог:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f (%.2f%%)", data.Tax, data.PercentTax)},
			{Text: "Брак:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f (%.2f%%)", data.Marriage, data.PercentMarriage)},
			{Text: "Ставка:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.SalaryRate)},
			{Text: "Ставка+:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.ExtendedSalaryRate)},
			{Text: "Итого:", Bold: true, Block: true}, {Text: fmt.Sprintf("%.2f р.", data.Margin)},
			{Block: true},
		},
	}

	return report, nil
}
