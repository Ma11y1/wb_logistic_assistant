package prompters

import (
	"fmt"
	"time"
)

const prefixCLIReporterFinanceDailyPrompter = "[Отчет суточных финансов маршрута]"

type CLIReporterFinanceDailyPrompter struct {
}

func (p *CLIReporterFinanceDailyPrompter) PromptStart(date time.Time) {
	fmt.Printf("%s Старт формирования... Дата: %s\n", prefixCLIReporterFinanceDailyPrompter, date.Format("02.01.2006"))
}

func (p *CLIReporterFinanceDailyPrompter) PromptFinish(duration time.Duration) {
	fmt.Println(prefixCLIReporterFinanceDailyPrompter, "Сформирован:", duration)
}

func (p *CLIReporterFinanceDailyPrompter) PromptCountWaySheet(total, closed, opened int) {
	fmt.Printf("%s Количество путевых листов: %d  Закрыто: %d  Открыто: %d\n", prefixCLIReporterFinanceDailyPrompter, total, closed, opened)
}

func (p *CLIReporterFinanceDailyPrompter) PromptSendReport(routeID int) {
	fmt.Printf("%s Отчет отправлен. Маршрут: %d\n", prefixCLIReporterFinanceDailyPrompter, routeID)
}

func (p *CLIReporterFinanceDailyPrompter) PromptError(message string) {
	fmt.Println(prefixCLIReporterFinanceDailyPrompter, "Ошибка:", message)
}
