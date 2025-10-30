package auth

import (
	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

type ActorType string

const (
	OAuthActorType   ActorType = "oauth"
	ServiceActorType ActorType = "service"
)

type Actor interface {
	Type() ActorType
	Service() *sheets.Service
	Spreadsheets() *sheets.SpreadsheetsService
	Token() *oauth2.Token
	IsValidToken() bool
	IsAuth() bool
	Close() error
}
