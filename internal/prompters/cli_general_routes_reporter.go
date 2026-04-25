package prompters

import (
	"fmt"
	"time"
)

const prefixCLIReporterGeneralRoutesPrompter = "[Отчет маршрутов]"

type CLIReporterGeneralRoutesPrompter struct {
}

func (p *CLIReporterGeneralRoutesPrompter) PromptStart() {
	fmt.Println(prefixCLIReporterGeneralRoutesPrompter, "Старт формирования...")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptFinish(duration time.Duration) {
	fmt.Println(prefixCLIReporterGeneralRoutesPrompter, "Сформирован:", duration)
}

func (p *CLIReporterGeneralRoutesPrompter) PromptUpdateRoutes(count int) {
	fmt.Printf("%s Обновлены маршруты: %d\n", prefixCLIReporterGeneralRoutesPrompter, count)
}

func (p *CLIReporterGeneralRoutesPrompter) PromptUpdateRating() {
	fmt.Println(prefixCLIReporterGeneralRoutesPrompter, "Обновлен рейтинг")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptUpdateShipments() {
	fmt.Println(prefixCLIReporterGeneralRoutesPrompter, "Обновлены отгрузки")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptCloseShipment(id, remainsBarcodes int) {
	fmt.Printf("%s Закрыта отгрузка: %d  Остаток ШК: %d\n", prefixCLIReporterGeneralRoutesPrompter, id, remainsBarcodes)
}

func (p *CLIReporterGeneralRoutesPrompter) PromptUpdateWaySheets() {
	fmt.Println(prefixCLIReporterGeneralRoutesPrompter, "Обновлены путевые листы")
}

func (p *CLIReporterGeneralRoutesPrompter) PromptSendReport(target string) {
	fmt.Println(prefixCLIReporterGeneralRoutesPrompter, "Отправлен:", target)
}

func (p *CLIReporterGeneralRoutesPrompter) PromptError(message string) {
	fmt.Println(prefixCLIReporterGeneralRoutesPrompter, "Ошибка:", message)

}
