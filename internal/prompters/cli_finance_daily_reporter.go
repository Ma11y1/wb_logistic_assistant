package prompters

import (
	"fmt"
	"time"
)

type CLIReporterFinanceDailyPrompter struct {
}

func (p *CLIReporterFinanceDailyPrompter) PromptStart(date time.Time) {
	fmt.Printf("Старт обновления [отчет суточных финансов маршрута] за %s...\n", date.Format("02.01.2006"))
}

func (p *CLIReporterFinanceDailyPrompter) PromptFinish(duration time.Duration) {
	fmt.Println("[Отчет суточных финансов маршрута] был обновлен. Время обновления: ", duration)
}

func (p *CLIReporterFinanceDailyPrompter) PromptCountWaySheet(total, closed, opened int) {
	fmt.Printf("[Отчет суточных финансов маршрута]: Количество путевых листов: всего: %d, закрытых: %d, открытых: %d\n", total, closed, opened)
}

func (p *CLIReporterFinanceDailyPrompter) PromptSendReport(routeID int) {
	fmt.Printf("[Отчет суточных финансов маршрута]: Отправлен отчет по маршруту: %d\n", routeID)
}

func (p *CLIReporterFinanceDailyPrompter) PromptError(message string) {
	fmt.Println("Ошибка составления [отчет суточных финансов маршрута]:", message)
}
