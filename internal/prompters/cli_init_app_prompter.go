package prompters

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CLIInitAppPrompter struct {
}

//// Google sheets

func (p *CLIInitAppPrompter) PromptGoogleSheetsAuthStart() {
	fmt.Println("Авторизация Google Sheets...")
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsQuestionAuthNewCredentials() bool {
	var res string
	fmt.Print("Войти под новыми правами доступа? (Y/N): ")
	fmt.Scanln(&res)
	return res == "Y" || res == "y"
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsRequestAuthCodeAuto(url string, seconds int) {
	fmt.Printf("Ссылка авторизации пользователя: %s\nОжидание %d секунд\n", url, seconds)
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsReadCredentialsFailed() {
	fmt.Println("Не удалось прочитать файл прав доступа к приложению Google Sheets")
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsAuthAutoFailed() {
	fmt.Println("Не удалось провести автоматическую авторизацию, необходимо продолжить вручную")
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsAuthFailed() {
	fmt.Println("Не удалось провести авторизацию пользователя")
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsAuthStorageFailed() {
	fmt.Println("Не удалось провести авторизацию пользователя, используя данные из хранилища")
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsRequestAuthCode(url string) (string, error) {
	fmt.Printf("Ссылка авторизации пользователя: %s\n", url)
	fmt.Print("Код: ")
	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	code = strings.TrimSpace(code)

	if code == "" {
		fmt.Println("Введен невалидный код")
		return code, fmt.Errorf("code is empty")
	}

	return code, nil
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsInvalidAuthCode() {
	fmt.Println("Введен невалидный код авторизации")
}

func (p *CLIInitAppPrompter) PromptGoogleSheetsAuthSuccessful() {
	fmt.Println("Авторизация в Google Sheets прошла успешно")
}

//// WB logistic

func (p *CLIInitAppPrompter) PromptWBLogisticAuthStart() {
	fmt.Println("Авторизация WB Logistic...")
}

func (p *CLIInitAppPrompter) PromptWBLogisticQuestionAuthNewUser() bool {
	var res string
	fmt.Print("Войти под новым пользователем? (Y/N): ")
	fmt.Scanln(&res)
	return res == "Y" || res == "y"
}

func (p *CLIInitAppPrompter) PromptWBLogisticRequestAuthLogin() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Логин (79991112233): ")
	login, err := reader.ReadString('\n')
	if login == "" || err != nil {
		return ""
	}
	login = strings.TrimSpace(login)
	login = strings.ReplaceAll(login, " ", "")
	return strings.TrimPrefix(login, "+")
}

func (p *CLIInitAppPrompter) PromptWBLogisticInvalidAuthLogin() {
	fmt.Println("Введен невалидный логин. Пример: 79991112233")
}

func (p *CLIInitAppPrompter) PromptWBLogisticRequestAuthCode(method string, time int) int {
	fmt.Printf("Вам отправлено %s уведомление с кодом. Повторная отправка через %d сек!\n", method, time)

	fmt.Print("Ввести номер действия или код!\n1. Повторить отправку кода\n2. Ручной ввод токена\n3. Выход\nКод: ")
	reader := bufio.NewReader(os.Stdin)
	codeStr, err := reader.ReadString('\n')
	if err != nil {
		return 0
	}

	codeStr = strings.TrimSpace(codeStr)

	if codeStr == "" || codeStr == "1" || codeStr == "2" || codeStr == "3" {
		code, _ := strconv.Atoi(codeStr)
		return code
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return 0
	}

	return code
}

func (p *CLIInitAppPrompter) PromptWBLogisticRequestAccessTokenData() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Токен: ")
	token, err := reader.ReadString('\n')
	if token == "" || err != nil {
		return ""
	}
	token = strings.TrimSpace(token)
	token = strings.ReplaceAll(token, " ", "")
	return token
}

func (p *CLIInitAppPrompter) PromptWBLogisticInvalidAccessTokenData() {
	fmt.Println("Введены невалидные данные токена доступа!")
}

func (p *CLIInitAppPrompter) PromptWBLogisticInvalidAuthCode() {
	fmt.Println("Введен невалидный код авторизации")
}

func (p *CLIInitAppPrompter) PromptWBLogisticAuthFailed() {
	fmt.Println("Не удалось провести авторизацию пользователя")
}

func (p *CLIInitAppPrompter) PromptWBLogisticAuthStorageFailed() {
	fmt.Println("Не удалось провести авторизацию пользователя, используя данные из хранилища")
}

func (p *CLIInitAppPrompter) PromptWBLogisticAuthSuccessful(login, username string) {
	fmt.Printf("Авторизация в WB logistic прошла успешно.\n\tПользователь: %s (%s)\n", username, login)
}

//// Telegram bot

func (p *CLIInitAppPrompter) PromptTelegramBotAuthStart() {
	fmt.Println("Авторизация Telegram bot...")
}

func (p *CLIInitAppPrompter) PromptTelegramBotQuestionAuthNewBot() bool {
	var res string
	fmt.Print("Войти под новым пользователем? (Y/N): ")
	fmt.Scanln(&res)
	return res == "Y" || res == "y"
}

func (p *CLIInitAppPrompter) PromptTelegramBotRequestToken() (string, error) {
	fmt.Print("Токен: ")
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	token = strings.TrimSpace(token)

	if token == "" {
		fmt.Println("Введен невалидный токен.")
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

func (p *CLIInitAppPrompter) PromptTelegramBotInitStorageFailed() {
	fmt.Println("Не удалось пройти инициализацию Telegram Bot используя данные из хранилища")
}

func (p *CLIInitAppPrompter) PromptTelegramBotInitFailed() {
	fmt.Println("Не удалось пройти инициализацию Telegram Bot")
}

func (p *CLIInitAppPrompter) PromptTelegramBotAuthSuccessful(name string) {
	fmt.Printf("Авторизация Telegram Bot '%s' прошла успешно.\n", name)
}

func (p *CLIInitAppPrompter) PromptInitFinish() {
	fmt.Println("****************************************************\n")
}
