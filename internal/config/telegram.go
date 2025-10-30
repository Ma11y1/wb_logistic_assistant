package config

import (
	"encoding/json"
)

type TelegramBot struct {
	chatID int64
}

type telegramBot struct {
	ChatID int64 `json:"chat_id"`
}

func newTelegramBot() *TelegramBot {
	return &TelegramBot{}
}

func (t *TelegramBot) ChatID() int64 {
	return t.chatID
}

func (t *TelegramBot) UnmarshalJSON(b []byte) error {
	temp := &telegramBot{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	t.chatID = temp.ChatID
	return nil
}

func (t *TelegramBot) MarshalJSON() ([]byte, error) {
	return json.Marshal(&telegramBot{ChatID: t.chatID})
}
