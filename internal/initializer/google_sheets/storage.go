package google_sheets

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"os"
	models2 "wb_logistic_assistant/external/google_sheets_api/models"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/logger"
	"wb_logistic_assistant/internal/models"
)

func (i *Initializer) GetOAuthCredentialsFile() (*models2.OAuthCredentials, error) {
	logger.Log(logger.INFO, "Initializer.GetOAuthCredentialsFile()", "start getting oauth credentials from file")

	path := i.config.GoogleSheets().Client().OAuthCredentials()
	if path == "" {
		return nil, errors.New("Initializer.GetOAuthCredentialsFile()", "empty path for oauth credentials")
	}

	credentialsData, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Initializer.GetOAuthCredentialsFile()", "failed to read oauth credentials from file %s", path)
	}

	credentials := &models2.OAuthCredentials{}
	err = json.Unmarshal(credentialsData, credentials)
	if err != nil {
		return nil, errors.Wrapf(err, "Initializer.GetOAuthCredentialsFile()", "failed to decode oauth credentials from file %s", path)
	}

	return credentials, nil
}

func (i *Initializer) GetServiceCredentialsFile() (*models2.ServiceCredentials, error) {
	logger.Log(logger.INFO, "Initializer.GetServiceCredentialsFile()", "start getting service credentials from file")

	path := i.config.GoogleSheets().Client().ServiceCredentials()
	if path == "" {
		return nil, errors.New("Initializer.GetServiceCredentialsFile()", "empty path for service credentials")
	}

	credentialsData, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Initializer.GetServiceCredentialsFile()", "failed to read service credentials from file %s", path)
	}

	credentials := &models2.ServiceCredentials{}
	err = json.Unmarshal(credentialsData, credentials)
	if err != nil {
		return nil, errors.Wrapf(err, "Initializer.GetServiceCredentialsFile()", "failed to decode service credentials from file %s", path)
	}

	return credentials, nil
}

func (i *Initializer) GetOAuthTokenStorage() (*models2.OAuthToken, error) {
	tokenStorage := i.storage.ConfigStore().GetGoogleSheetsOAuthTokenModel()
	if tokenStorage == nil {
		return nil, errors.New("Initializer.GetOAuthTokenStorage()", "oauth token are missing from storage")
	}

	token := &models2.OAuthToken{
		TokenType:    tokenStorage.TokenType(),
		AccessToken:  tokenStorage.AccessToken(),
		RefreshToken: tokenStorage.RefreshToken(),
		Expiry:       tokenStorage.Expiry(),
		ExpiresIn:    tokenStorage.ExpiresIn(),
	}

	if len(token.AccessToken) == 0 || len(token.RefreshToken) == 0 {
		return nil, errors.New("Initializer.GetOAuthTokenStorage()", "invalid access token or refresh token from storage")
	}
	return token, nil
}

func (i *Initializer) SetOAuthTokenStorage(token *oauth2.Token) error {
	if token == nil {
		return errors.New("Initializer.SetOAuthTokenStorage()", "invalid oauth token provided")
	}

	tokenModel := &models.GoogleSheetsOAuthTokenModel{}
	tokenModel.SetAll(
		token.TokenType,
		token.AccessToken,
		token.RefreshToken,
		token.Expiry,
		token.ExpiresIn,
	)

	i.storage.ConfigStore().SetGoogleSheetsOAuthToken(tokenModel)

	logger.Log(logger.INFO, "Initializer.SetOAuthTokenStorage()", "google sheets oauth access token save to storage")
	return nil
}

func (i *Initializer) GetOAuthCredentialsStorage() (*models2.OAuthCredentials, error) {
	credentialsModel := i.storage.ConfigStore().GetGoogleSheetsOAuthCredentials()
	if credentialsModel == nil {
		return nil, errors.New("Initializer.GetOAuthCredentialsStorage()", "oauth credentials are missing from storage")
	}

	credentials := &models2.OAuthCredentials{}
	credentials.Installed.ProjectID = credentialsModel.ProjectID()
	credentials.Installed.ClientID = credentialsModel.ClientID()
	credentials.Installed.ClientSecret = credentialsModel.ClientSecret()
	credentials.Installed.AuthUri = credentialsModel.AuthUri()
	credentials.Installed.TokenUri = credentialsModel.TokenUri()
	credentials.Installed.AuthProviderX509CertUrl = credentialsModel.AuthProviderX509CertUrl()
	credentials.Installed.RedirectUris = credentialsModel.RedirectUris()

	if err := credentials.Validate(); err != nil {
		return nil, errors.Wrap(err, "Initializer.GetOAuthCredentialsStorage()", "invalid oauth credentials from storage")
	}

	return credentials, nil
}

func (i *Initializer) SetOAuthCredentialsStorage(credentials *models2.OAuthCredentials) error {
	if credentials == nil {
		return errors.New("Initializer.SetOAuthCredentialsStorage()", "invalid oauth credentials provided")
	}

	credentialsModel := &models.GoogleSheetsOAuthCredentialsModel{}
	credentialsModel.SetAll(
		credentials.Installed.ClientID,
		credentials.Installed.ProjectID,
		credentials.Installed.AuthUri,
		credentials.Installed.TokenUri,
		credentials.Installed.AuthProviderX509CertUrl,
		credentials.Installed.ClientSecret,
		credentials.Installed.RedirectUris,
	)
	i.storage.ConfigStore().SetGoogleSheetsOAuthCredentials(credentialsModel)

	logger.Logf(logger.INFO, "Initializer.SetOAuthCredentialsStorage()", "google sheets oauth credentials of project '%s' save to storage", credentials.Installed.ProjectID)
	return nil
}

func (i *Initializer) GetServiceCredentialsStorage() (*models2.ServiceCredentials, error) {
	credentialsModel := i.storage.ConfigStore().GetGoogleSheetsServiceCredentials()
	if credentialsModel == nil {
		return nil, errors.New("Initializer.GetServiceCredentialsStorage()", "service credentials are missing from storage")
	}

	credentials := &models2.ServiceCredentials{}
	credentials.Type = credentialsModel.Type()
	credentials.ProjectID = credentialsModel.ProjectID()
	credentials.PrivateKeyID = credentialsModel.PrivateKeyID()
	credentials.PrivateKey = credentialsModel.PrivateKey()
	credentials.ClientEmail = credentialsModel.ClientEmail()
	credentials.ClientID = credentialsModel.ClientID()
	credentials.AuthUri = credentialsModel.AuthUri()
	credentials.TokenUri = credentialsModel.TokenUri()
	credentials.AuthProviderX509CertUrl = credentialsModel.AuthProviderX509CertUrl()
	credentials.ClientX509CertUrl = credentialsModel.ClientX509CertUrl()
	credentials.UniverseDomain = credentialsModel.UniverseDomain()

	if err := credentials.Validate(); err != nil {
		return nil, errors.Wrap(err, "Initializer.GetServiceCredentialsStorage()", "invalid service credentials from storage")
	}

	return credentials, nil
}

func (i *Initializer) SetServiceCredentialsStorage(credentials *models2.ServiceCredentials) error {
	if credentials == nil {
		return errors.New("Initializer.SetServiceCredentialsStorage()", "invalid service credentials provided")
	}

	credentialsModel := &models.GoogleSheetsServiceCredentialsModel{}
	credentialsModel.SetAll(
		credentials.Type,
		credentials.ProjectID,
		credentials.PrivateKeyID,
		credentials.PrivateKey,
		credentials.ClientEmail,
		credentials.ClientID,
		credentials.AuthUri,
		credentials.TokenUri,
		credentials.AuthProviderX509CertUrl,
		credentials.ClientX509CertUrl,
		credentials.UniverseDomain,
	)

	i.storage.ConfigStore().SetGoogleSheetsServiceCredentials(credentialsModel)

	logger.Logf(logger.INFO, "Initializer.SetServiceCredentialsStorage()", "google sheets service credentials of project '%s' save to storage", credentials.ProjectID)
	return nil
}
