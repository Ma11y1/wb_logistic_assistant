package prompters

import (
	"fmt"
	"time"
)

type CLIReporterGeneralRoutesPrompter struct {
}

func (p *CLIReporterGeneralRoutesPrompter) PromptStart() {
	fmt.Println("Старт обновления [отчет маршрутов]...")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptFinish(duration time.Duration) {
	fmt.Println("[Отчет маршрутов] был обновлен. Время обновления: ", duration)
}

func (p *CLIReporterGeneralRoutesPrompter) PromptUpdateRoutes(count int) {
	fmt.Printf("[Отчет маршрутов]: обновлены маршруты %d\n", count)
}

func (p *CLIReporterGeneralRoutesPrompter) PromptUpdateRating() {
	fmt.Println("[Отчет маршрутов]: обновлен рейтинг")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptUpdateShipments() {
	fmt.Println("[Отчет маршрутов]: обновлены отгрузки")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptCloseShipment(id, remainsBarcodes int) {
	fmt.Printf("[Отчет маршрутов]: закрыта отгрузка %d, остаток ШК: %d\n", id, remainsBarcodes)
}

func (p *CLIReporterGeneralRoutesPrompter) PromptUpdateWaySheets() {
	fmt.Println("[Отчет маршрутов]: обновлены путевые листы")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptSendReport() {
	fmt.Println("Отправлен [отчет маршрутов]")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptError(message string) {
	fmt.Println("Ошибка составления [отчет маршрутов]:", message)

}
