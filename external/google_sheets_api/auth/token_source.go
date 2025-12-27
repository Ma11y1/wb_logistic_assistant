package auth

import "golang.org/x/oauth2"

type RefreshTokenHandler func(*oauth2.Token)
type ErrorHandler func(error)

type TokenSource struct {
	oauth2.TokenSource
	RefreshTokenHandler RefreshTokenHandler
	ErrorHandler        func(error)
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token, err := t.TokenSource.Token()
	if err != nil {
		if t.ErrorHandler != nil {
			t.ErrorHandler(err)
		}
		return nil, err
	}
	if t.RefreshTokenHandler != nil {
		t.RefreshTokenHandler(token)
	}
	return token, nil
}
