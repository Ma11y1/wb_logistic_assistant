package telegram_bot

import (
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/prompters"
	"wb_logistic_assistant/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Initializer struct {
	config   *config.Config
	storage  storage.Storage
	prompter prompters.InitializeAppPrompter
}

func NewInitializer(config *config.Config, storage storage.Storage, prompter prompters.InitializeAppPrompter) *Initializer {
	return &Initializer{
		config:   config,
		storage:  storage,
		prompter: prompter,
	}
}

func (i *Initializer) Init() (*tgbotapi.BotAPI, error) {
	logger.Log(logger.INFO, "Initializer.Init()", "Start init Telegram Bot")
	i.prompter.PromptTelegramBotAuthStart()

	var token string
	var err error
	if i.prompter.PromptTelegramBotQuestionAuthNewBot() {
		token, err = i.prompter.PromptTelegramBotRequestToken()
		if err != nil {
			i.prompter.PromptTelegramBotInitFailed()
			return nil, errors.Wrap(err, "Initializer.Init()", "failed to receive telegram bot token")
		}
	} else {
		token, err = i.GetToken()
		if err != nil {
			i.prompter.PromptTelegramBotInitStorageFailed()
			logger.Logf(logger.WARN, "Initializer.Init()", "failed to receive telegram bot token using storage: %v", err)

			token, err = i.prompter.PromptTelegramBotRequestToken()
			if err != nil {
				i.prompter.PromptTelegramBotInitFailed()
				return nil, errors.Wrap(err, "Initializer.Init()", "failed to receive telegram bot token")
			}
		}
	}
	if token == "" {
		return nil, errors.New("Initializer.Init()", "received empty telegram bot token")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		i.prompter.PromptTelegramBotInitFailed()
		return nil, errors.Wrap(err, "Initializer.Init()", "failed to init telegram bot")
	}

	err = i.SetToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "Initializer.Init()", "failed to save telegram bot to storage")
	}

	err = i.UpdateStorage()
	if err != nil {
		i.prompter.PromptTelegramBotInitFailed()
		return nil, errors.Wrap(err, "Initializer.Init()", "failed to update storage")
	}

	i.prompter.PromptTelegramBotAuthSuccessful(bot.Self.UserName)
	logger.Log(logger.INFO, "Initializer.Init()", "Finish init Telegram Bot")
	return bot, nil
}
