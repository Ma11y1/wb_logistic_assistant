package models

import (
	"encoding/json"
	"fmt"
)

type OAuthCredentials struct {
	Installed struct {
		ProjectID               string   `json:"project_id"`
		ClientID                string   `json:"client_id"`
		ClientSecret            string   `json:"client_secret"`
		AuthUri                 string   `json:"auth_uri"`
		TokenUri                string   `json:"token_uri"`
		AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
		RedirectUris            []string `json:"redirect_uris"`
	} `json:"installed"`
}

func (c *OAuthCredentials) Validate() error {
	if c.Installed.ProjectID == "" {
		return fmt.Errorf("missing installed project id")
	}
	if c.Installed.ClientID == "" {
		return fmt.Errorf("missing installed client id")
	}
	if c.Installed.ClientSecret == "" {
		return fmt.Errorf("missing installed client secret")
	}
	if c.Installed.AuthUri == "" {
		return fmt.Errorf("missing installed session uri")
	}
	if c.Installed.TokenUri == "" {
		return fmt.Errorf("missing installed token uri")
	}
	if c.Installed.AuthProviderX509CertUrl == "" {
		return fmt.Errorf("missing installed session provider x509 cert url")
	}
	return nil
}

func (c *OAuthCredentials) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

type ServiceCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

func (c *ServiceCredentials) Validate() error {
	if c.Type == "" {
		return fmt.Errorf("missing credential type")
	}
	if c.ProjectID == "" {
		return fmt.Errorf("missing credential project id")
	}
	if c.PrivateKeyID == "" {
		return fmt.Errorf("missing credential private key id")
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("missing credential private key")
	}
	if c.ClientID == "" {
		return fmt.Errorf("missing credential client id")
	}
	if c.ClientEmail == "" {
		return fmt.Errorf("missing credential client email")
	}
	if c.AuthUri == "" {
		return fmt.Errorf("missing credential session uri")
	}
	if c.TokenUri == "" {
		return fmt.Errorf("missing credential token uri")
	}
	if c.AuthProviderX509CertUrl == "" {
		return fmt.Errorf("missing credential session provider x509 cert url")
	}
	if c.ClientX509CertUrl == "" {
		return fmt.Errorf("missing credential client x509 cert url")
	}
	if c.UniverseDomain == "" {
		return fmt.Errorf("missing credential universe domain")
	}
	return nil
}

func (c *ServiceCredentials) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}
