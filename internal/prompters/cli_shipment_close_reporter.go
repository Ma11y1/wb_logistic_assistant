package prompters

import (
	"fmt"
	"time"
)

const prefixCLIReporterShipmentClosePrompter = "[Отчет закрытия отгрузок]"

type CLIReporterShipmentClosePrompter struct {
}

func (p *CLIReporterShipmentClosePrompter) PromptStart() {
	fmt.Println(prefixCLIReporterShipmentClosePrompter, "Старт формирования...")
}

func (p *CLIReporterShipmentClosePrompter) PromptFinish(duration time.Duration) {
	fmt.Println(prefixCLIReporterShipmentClosePrompter, "Сформирован:", duration)
}

func (p *CLIReporterShipmentClosePrompter) PromptShipmentOpened(routeID, shipmentID, opened int) {
	fmt.Printf("%s Открыта отгрузка: %d  Маршрут: %d  Всего: %d\n", prefixCLIReporterShipmentClosePrompter, shipmentID, routeID, opened)
}

func (p *CLIReporterShipmentClosePrompter) PromptShipmentClose(routeID, shipmentID int) {
	fmt.Printf("%s Закрыта отгрузка: %d  Маршрут: %d\n", prefixCLIReporterShipmentClosePrompter, shipmentID, routeID)
}

func (p *CLIReporterShipmentClosePrompter) PromptSendReport(routeID, shipmentID, waySheetID int) {
	fmt.Printf("%s Отчет отправлен. Отгрузка: %d  Маршрут: %d  Путевой лист: %d\n", prefixCLIReporterShipmentClosePrompter, shipmentID, routeID, waySheetID)
}

func (p *CLIReporterShipmentClosePrompter) PromptError(message string) {
	fmt.Println(prefixCLIReporterShipmentClosePrompter, "Ошибка:", message)
}
