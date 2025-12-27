package prompters

import (
	"fmt"
	"time"
)

type CLIReporterFinanceRoutesPrompter struct {
}

func (p *CLIReporterFinanceRoutesPrompter) PromptStart() {
	fmt.Println("Старт обновления [отчет финансов маршрута]...")
}

func (p *CLIReporterFinanceRoutesPrompter) PromptFinish(duration time.Duration) {
	fmt.Println("[Отчет финансов маршрута] был обновлен. Время обновления: ", duration)
}

func (p *CLIReporterFinanceRoutesPrompter) PromptCountWaySheet(count int) {
	fmt.Printf("[Отчет финансов маршрута]: Количество открытых путевых листов: %d\n", count)
}

func (p *CLIReporterFinanceRoutesPrompter) PromptCloseWaySheet(routeID int, waySheetID, shipmentID string) {
	fmt.Printf("[Отчет финансов маршрута]: Закрыт путевой лист: %s на маршруте %d, отгрузка %s\n", waySheetID, routeID, shipmentID)
}

func (p *CLIReporterFinanceRoutesPrompter) PromptSendReport(routeID int, waySheetID, shipmentID string) {
	fmt.Printf("[Отчет финансов маршрута]: Отправлен отчет по путевому листу: %s на маршруте %d, отгрузка %s\n", waySheetID, routeID, shipmentID)
}

func (p *CLIReporterFinanceRoutesPrompter) PromptError(message string) {
	fmt.Println("Ошибка составления [отчет финансов маршрута]:", message)
}
