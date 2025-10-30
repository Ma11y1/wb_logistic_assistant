package prompters

import (
	"fmt"
	"time"
)

type CLIReporterShipmentClosePrompter struct {
}

func (p *CLIReporterShipmentClosePrompter) PromptError(message string) {
	fmt.Println("Ошибка составления отчета 'Shipment close':", message)

}
func (p *CLIReporterShipmentClosePrompter) PromptRender(duration time.Duration) {
	fmt.Println("Отчет 'Shipment close' был обновлен за: ", duration)
}
