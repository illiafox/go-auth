package github

import (
	"net/url"
)

type Auth struct {
	authURL  string
	getToken string
}

func (a Auth) Auth(state string) string {
	return a.authURL + url.QueryEscape(state)
}

func (a Auth) Token(code string) string {
	return a.getToken + url.QueryEscape(code)
}
