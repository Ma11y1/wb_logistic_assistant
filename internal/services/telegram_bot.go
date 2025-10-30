package services

import (
	"fmt"
	"strings"
	"wb_logistic_assistant/internal/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const TelegramMessageSizeLimit = 4096

type TelegramBotService interface {
	SendMessage(chatID int64, message string, parseMode string) error
	SendPhotoFile(chatID int64, photoPath, caption string) error
	SendDocumentFile(chatID int64, filePath, caption string) error
	GetUpdates(offset, limit, timeout int) ([]tgbotapi.Update, error)
	GetBotInfo() (*tgbotapi.User, error)
	HandleCommands(updates []tgbotapi.Update, handlers map[string]func(update tgbotapi.Update)) error
}

type TelegramBotAPIService struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramBotAPIService(bot *tgbotapi.BotAPI) *TelegramBotAPIService {
	return &TelegramBotAPIService{bot: bot}
}

func (s *TelegramBotAPIService) SendMessage(chatID int64, message string, parseMode string) error {
	if parseMode == "" {
		parseMode = tgbotapi.ModeHTML
	}

	runes := []rune(message)
	totalLen := len(runes)

	for start := 0; start < totalLen; start += TelegramMessageSizeLimit {
		end := start + TelegramMessageSizeLimit
		if end > totalLen {
			end = totalLen
		}

		part := string(runes[start:end])
		msg := tgbotapi.NewMessage(chatID, part)
		msg.ParseMode = parseMode

		_, err := s.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "TelegramBotAPIService.SendMessage()", "failed send message part")
		}
	}

	return nil
}

func (s *TelegramBotAPIService) SendPhotoFile(chatID int64, photoPath, caption string) error {
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(photoPath))
	photo.Caption = caption
	_, err := s.bot.Send(photo)
	if err != nil {
		return errors.Wrap(err, "TelegramBotAPIService.SendPhotoFile()", "failed send photo file")
	}
	return nil
}

func (s *TelegramBotAPIService) SendDocumentFile(chatID int64, filePath, caption string) error {
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filePath))
	doc.Caption = caption
	_, err := s.bot.Send(doc)
	if err != nil {
		return errors.Wrap(err, "TelegramBotAPIService.SendDocumentFile()", "failed send document file")
	}
	return nil
}

func (s *TelegramBotAPIService) GetUpdates(offset, limit, timeout int) ([]tgbotapi.Update, error) {
	cfg := tgbotapi.NewUpdate(offset)
	cfg.Limit = limit
	cfg.Timeout = timeout

	updates, err := s.bot.GetUpdates(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "TelegramBotAPIService.GetUpdates()", "failed get updates")
	}
	return updates, nil
}

func (s *TelegramBotAPIService) GetBotInfo() (*tgbotapi.User, error) {
	user, err := s.bot.GetMe()
	if err != nil {
		return nil, errors.Wrap(err, "TelegramBotAPIService.GetBotInfo()", "failed get bot info")
	}
	return &user, nil
}

// HandleCommands calls command handlers ("/start", "/help", etc.)
func (s *TelegramBotAPIService) HandleCommands(updates []tgbotapi.Update, handlers map[string]func(update tgbotapi.Update)) error {
	if updates == nil {
		return errors.New("TelegramBotAPIService.HandleCommands()", "updates is nil")
	}
	if handlers == nil {
		return errors.New("TelegramBotAPIService.HandleCommands()", "handlers is nil")
	}

	for _, update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		cmd := strings.ToLower(update.Message.Command())
		if handler, ok := handlers[cmd]; ok {
			handler(update)
		} else {
			_ = s.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Invalid command: /%s", cmd), "")
		}
	}

	return nil
}
