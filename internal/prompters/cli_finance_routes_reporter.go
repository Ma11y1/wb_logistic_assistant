package prompters

import (
	"fmt"
	"time"
)

const prefixCLIReporterFinanceRoutesPrompter = "[Отчет финансов маршрута]"

type CLIReporterFinanceRoutesPrompter struct {
}

func (p *CLIReporterFinanceRoutesPrompter) PromptStart() {
	fmt.Println(prefixCLIReporterFinanceRoutesPrompter, "Старт формирования...")
}

func (p *CLIReporterFinanceRoutesPrompter) PromptFinish(duration time.Duration) {
	fmt.Println(prefixCLIReporterFinanceRoutesPrompter, "Сформирован:", duration)
}

func (p *CLIReporterFinanceRoutesPrompter) PromptCountWaySheet(count int) {
	fmt.Printf("%s Количество открытых путевых листов: %d\n", prefixCLIReporterFinanceRoutesPrompter, count)
}

func (p *CLIReporterFinanceRoutesPrompter) PromptCloseWaySheet(routeID int, waySheetID, shipmentID string) {
	fmt.Printf("%s Закрыт путевой лист: %s  Маршрут: %d  Отгрузка: %s\n", prefixCLIReporterFinanceRoutesPrompter, waySheetID, routeID, shipmentID)
}

func (p *CLIReporterFinanceRoutesPrompter) PromptSendReport(routeID int, waySheetID, shipmentID string) {
	fmt.Printf("%s Отчет отправлен. Путевой лист: %s  Маршрут: %d  Отгрузка: %s\n", prefixCLIReporterFinanceRoutesPrompter, waySheetID, routeID, shipmentID)
}

func (p *CLIReporterFinanceRoutesPrompter) PromptError(message string) {
	fmt.Println(prefixCLIReporterFinanceRoutesPrompter, "Ошибка:", message)
}
