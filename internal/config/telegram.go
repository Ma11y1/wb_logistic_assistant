package config

import (
	"encoding/json"
)

type TelegramBot struct {
	shipmentClose *TelegramBotParams
	financeRoutes *TelegramBotParams
	financeDaily  *TelegramBotParams
}

type telegramBot struct {
	ShipmentClose *TelegramBotParams `json:"shipment_close"`
	FinanceRoutes *TelegramBotParams `json:"finance_routes"`
	FinanceDaily  *TelegramBotParams `json:"finance_daily"`
}

func newTelegramBot() *TelegramBot {
	return &TelegramBot{
		shipmentClose: newTelegramBotParams(),
		financeRoutes: newTelegramBotParams(),
		financeDaily:  newTelegramBotParams(),
	}
}

func (t *TelegramBot) ShipmentClose() *TelegramBotParams {
	return t.shipmentClose
}
func (t *TelegramBot) FinanceRoutes() *TelegramBotParams {
	return t.financeRoutes
}
func (t *TelegramBot) FinanceDaily() *TelegramBotParams {
	return t.financeDaily
}

func (t *TelegramBot) UnmarshalJSON(b []byte) error {
	temp := &telegramBot{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	t.shipmentClose = temp.ShipmentClose
	t.financeRoutes = temp.FinanceRoutes
	t.financeDaily = temp.FinanceDaily
	return nil
}

func (t *TelegramBot) MarshalJSON() ([]byte, error) {
	return json.Marshal(&telegramBot{
		ShipmentClose: t.shipmentClose,
		FinanceRoutes: t.financeRoutes,
		FinanceDaily:  t.financeDaily,
	})
}

type TelegramBotParams struct {
	chatID int64
}

type telegramBotParams struct {
	ChatID int64 `json:"chat_id"`
}

func newTelegramBotParams() *TelegramBotParams {
	return &TelegramBotParams{}
}

func (t *TelegramBotParams) ChatID() int64 {
	return t.chatID
}

func (t *TelegramBotParams) UnmarshalJSON(b []byte) error {
	temp := &telegramBotParams{}
	err := json.Unmarshal(b, temp)
	if err != nil {
		return err
	}
	t.chatID = temp.ChatID
	return nil
}

func (t *TelegramBotParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(&telegramBotParams{ChatID: t.chatID})
}
