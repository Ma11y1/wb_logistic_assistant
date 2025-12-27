package wb_logistic

import (
	"wb_logistic_assistant/external/wb_logistic_api"
	"wb_logistic_assistant/external/wb_logistic_api/session"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
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

func (i *Initializer) Init() (*wb_logistic_api.Client, *session.Session, error) {
	logger.Log(logger.INFO, "Initializer.Init()", "Start init wb logistic")

	cfg := i.config.Logistic().WBClient()
	httpClient := transport.NewBaseHTTPClientWithParams(&transport.HTTPClientParameters{
		UserAgent:    cfg.UserAgent(),
		Platform:     cfg.Platform(),
		SecUserAgent: cfg.SecUserAgent(),
	})
	client := wb_logistic_api.NewClient(httpClient)

	i.prompter.PromptWBLogisticAuthStart()
	var s *session.Session
	var err error
	if i.prompter.PromptWBLogisticQuestionAuthNewUser() {
		s, err = i.AuthSession(client)
		if err != nil {
			i.prompter.PromptWBLogisticAuthFailed()
			return nil, nil, errors.Wrap(err, "Initializer.Init()", "failed to auth wb logistic session")
		}
	} else {
		s, err = i.authSessionStorage(client)
		if err != nil {
			i.prompter.PromptWBLogisticAuthStorageFailed()
			logger.Logf(logger.WARN, "Initializer.Init()", "failed to auth wb logistic session using storage: %v", err)

			s, err = i.AuthSession(client)
			if err != nil {
				i.prompter.PromptWBLogisticAuthFailed()
				return nil, nil, errors.Wrap(err, "Initializer.Init()", "failed to auth wb logistic session")
			}
		}
	}

	i.prompter.PromptWBLogisticAuthSuccessful(s.Login(), s.UserInfo().UserDetails.Name)

	logger.Log(logger.INFO, "Initializer.Init()", "Finish init wb logistic")
	return client, s, nil
}

// InitDirect init without storage data
func (i *Initializer) InitDirect() (*wb_logistic_api.Client, *session.Session, error) {
	logger.Log(logger.INFO, "Initializer.InitDirect()", "Start init wb logistic")

	cfg := i.config.Logistic().WBClient()
	httpClient := transport.NewBaseHTTPClientWithParams(&transport.HTTPClientParameters{
		UserAgent:    cfg.UserAgent(),
		Platform:     cfg.Platform(),
		SecUserAgent: cfg.SecUserAgent(),
	})
	client := wb_logistic_api.NewClient(httpClient)

	i.prompter.PromptWBLogisticAuthStart()
	s, err := i.AuthSession(client)
	if err != nil {
		i.prompter.PromptWBLogisticAuthFailed()
		return nil, nil, errors.Wrap(err, "Initializer.InitDirect()", "failed to auth wb logistic session")
	}
	i.prompter.PromptWBLogisticAuthSuccessful(s.Login(), s.UserInfo().UserDetails.Name)

	logger.Log(logger.INFO, "Initializer.InitDirect()", "Finish init wb logistic")
	return client, s, nil
}
