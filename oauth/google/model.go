package google

import (
	"net/url"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	authURL  string
	getToken string
	keyFunc  jwt.Keyfunc
}

func (a Auth) Keyfunc(token *jwt.Token) (interface{}, error) {
	return a.keyFunc(token)
}

func (a Auth) Auth(state string) string {
	return a.authURL + url.QueryEscape(state)
}

func (a Auth) Token(code string) string {
	return a.getToken + url.QueryEscape(code)
}
