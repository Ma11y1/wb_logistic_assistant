package google_sheets

import (
	"wb_logistic_assistant/external/google_sheets_api"
	"wb_logistic_assistant/external/google_sheets_api/auth"
	"wb_logistic_assistant/internal/config"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/prompters"
	"wb_logistic_assistant/internal/storage"
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

func (i *Initializer) Init() (*google_sheets_api.Client, auth.Actor, error) {
	logger.Log(logger.INFO, "Initializer.Init()", "Start init Google Sheets")

	client := google_sheets_api.NewClient()

	i.prompter.PromptGoogleSheetsAuthStart()
	var actor auth.Actor
	var err error
	if i.config.GoogleSheets().Client().IsOAuth() {
		actor, err = i.authOAuthSession()
	} else {
		actor, err = i.authServiceSession()
	}
	if err != nil {
		i.prompter.PromptGoogleSheetsAuthFailed()
		return nil, nil, errors.Wrap(err, "Initializer.Init()", "failed to auth google sheets client session")
	}

	err = i.UpdateStorage()
	if err != nil {
		i.prompter.PromptGoogleSheetsAuthFailed()
		return nil, nil, errors.Wrap(err, "Initializer.Init()", "failed to update storage")
	}

	i.prompter.PromptGoogleSheetsAuthSuccessful()

	logger.Log(logger.INFO, "Initializer.Init()", "Finish init Google Sheets")
	return client, actor, nil
}
