package google

import (
	"net/url"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	authURL  string
	getToken string
	Keys     jwt.Keyfunc
}

func (a Auth) Auth(state string) string {
	return a.authURL + url.QueryEscape(state)
}

func (a Auth) Token(code string) string {
	return a.getToken + url.QueryEscape(code)
}

func (a Auth) KeyFunc() jwt.Keyfunc {
	return a.Keys
}
