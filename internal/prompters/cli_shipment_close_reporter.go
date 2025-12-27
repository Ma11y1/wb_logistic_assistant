package prompters

import (
	"fmt"
	"time"
)

type CLIReporterShipmentClosePrompter struct {
}

func (p *CLIReporterShipmentClosePrompter) PromptStart() {
	fmt.Println("Старт обновления [отчет закрытия отгрузок]...")
}

func (p *CLIReporterShipmentClosePrompter) PromptFinish(duration time.Duration) {
	fmt.Println("[Отчет закрытия отгрузок] был обновлен. Время обновления: ", duration)
}

func (p *CLIReporterShipmentClosePrompter) PromptShipmentOpened(routeID, shipmentID, opened int) {
	fmt.Printf("[Отчет закрытия отгрузок]: Открыта отгрузка: %d на маршруте: %d. Всего: %d\n", shipmentID, routeID, opened)
}

func (p *CLIReporterShipmentClosePrompter) PromptShipmentClose(routeID, shipmentID int) {
	fmt.Printf("[Отчет закрытия отгрузок]: Закрыта отгрузка: %d на маршруте: %d\n", shipmentID, routeID)
}

func (p *CLIReporterShipmentClosePrompter) PromptSendReport(routeID, shipmentID, waySheetID int) {
	fmt.Printf("[Отчет закрытия отгрузок]: Отправлен отчет по отгрузке: %d на маршруте: %d. Установлен путевой лист: %d\n", shipmentID, routeID, waySheetID)
}

func (p *CLIReporterShipmentClosePrompter) PromptError(message string) {
	fmt.Println("Ошибка составления [отчет закрытия отгрузок]:", message)
}
