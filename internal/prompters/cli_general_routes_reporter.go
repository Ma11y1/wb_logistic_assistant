package prompters

import "fmt"

type CLIReporterGeneralRoutesPrompter struct {
}

func (p *CLIReporterGeneralRoutesPrompter) PromptError(msg string) {

}
func (p *CLIReporterGeneralRoutesPrompter) PromptRender() { // todo сделать передачу времени выполнения
	fmt.Println("Отчет General Routes был обновлен")
}
