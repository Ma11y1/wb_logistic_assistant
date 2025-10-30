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
	fmt.Println("Начало авторизации Google Sheets")
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

//// Ol logistic

func (p *CLIInitAppPrompter) PromptOLLogisticAuthStart() {
	fmt.Println("Начало авторизации Ol logistic")
}

func (p *CLIInitAppPrompter) PromptOLLogisticQuestionAuthNewUser() bool {
	var res string
	fmt.Print("Войти под новым пользователем? (Y/N): ")
	fmt.Scanln(&res)
	return res == "Y" || res == "y"
}

func (p *CLIInitAppPrompter) PromptOLLogisticRequestAuthData() (string, string, error) {
	fmt.Print("Логин: ")
	reader := bufio.NewReader(os.Stdin)
	login, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	login = strings.TrimSpace(login)

	if login == "" {
		fmt.Println("Введен невалидный логин")
		return login, "", fmt.Errorf("login is empty")
	}

	login = strings.ReplaceAll(login, " ", "")
	login = strings.TrimPrefix(login, "+")

	fmt.Print("Пароль: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return login, "", err
	}

	password = strings.TrimSpace(password)

	if password == "" {
		fmt.Println("Введен невалидный пароль.")
		return login, password, fmt.Errorf("password is empty")
	}

	return login, password, nil
}

func (p *CLIInitAppPrompter) PromptOLLogisticInvalidAuthData() {
	fmt.Println("Введены невалидные данные авторизации")
}

func (p *CLIInitAppPrompter) PromptOLLogisticAuthFailed() {
	fmt.Println("Не удалось провести авторизацию пользователя")
}

func (p *CLIInitAppPrompter) PromptOLLogisticAuthStorageFailed() {
	fmt.Println("Не удалось провести авторизацию пользователя, используя данные из хранилища")
}

func (p *CLIInitAppPrompter) PromptOLLogisticAuthSuccessful(login, username string) {
	fmt.Printf("Авторизация в Ol logistic прошла успешно.\n\tПользователь: %s (%s)\n", username, login)
}

//// WB logistic

func (p *CLIInitAppPrompter) PromptWBLogisticAuthStart() {
	fmt.Println("Начало авторизации WB Logistic")
}

func (p *CLIInitAppPrompter) PromptWBLogisticQuestionAuthNewUser() bool {
	var res string
	fmt.Print("Войти под новым пользователем? (Y/N): ")
	fmt.Scanln(&res)
	return res == "Y" || res == "y"
}

func (p *CLIInitAppPrompter) PromptWBLogisticRequestAuthLogin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Логин: ")
	login, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	login = strings.TrimSpace(login)

	if login == "" {
		fmt.Println("Введен невалидный логин")
		return login, fmt.Errorf("login is empty")
	}

	login = strings.ReplaceAll(login, " ", "")
	login = strings.TrimPrefix(login, "+")

	return login, nil
}

func (p *CLIInitAppPrompter) PromptWBLogisticRequestAuthCode(method string, time int) (int, error) {
	fmt.Printf("Вам отправлено %s уведомление с кодом. Ожидание %d\n", method, time)

	fmt.Print("Код: ")
	reader := bufio.NewReader(os.Stdin)
	codeStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	codeStr = strings.TrimSpace(codeStr)

	if codeStr == "" {
		fmt.Println("Введен невалидный код.")
		return 0, fmt.Errorf("code is empty")
	}

	if len(codeStr) != 6 {
		fmt.Println("Введен невалидный код. Длина кода должна быть 6 символов.")
		return 0, fmt.Errorf("code is empty")
	}
	fmt.Println(codeStr)
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		fmt.Println("Введен невалидный код. Код состоит только из цифр.")
		return 0, fmt.Errorf("code is invalid: %v", err)
	}

	return code, nil
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
	fmt.Println("Начало авторизации Telegram bot")
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
