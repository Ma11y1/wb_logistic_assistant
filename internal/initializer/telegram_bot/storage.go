package telegram_bot

import "wb_logistic_assistant/internal/errors"

func (i *Initializer) GetToken() (string, error) {
	token := i.storage.ConfigStore().GetTelegramBotToken()
	if token == "" {
		return "", errors.New("Initializer.GetToken()", "no token found")
	}
	return token, nil
}

func (i *Initializer) SetToken(token string) error {
	if token == "" {
		return errors.New("Initializer.SetToken()", "token is empty")
	}
	i.storage.ConfigStore().SetTelegramBotToken(token)
	return nil
}

func (i *Initializer) UpdateStorage() error {
	err := i.storage.Save(i.config.Storage().Path())
	if err != nil {
		return errors.Wrap(err, "Initializer.UpdateStorage()", "failed to update storage")
	}
	return nil
}
