package google

import (
	"fmt"
	"net/url"

	"go-auth/oauth/parse"
)

func New(web *parse.Web) (*Auth, error) {
	// auth

	auth, err := url.Parse(web.AuthURL)
	if err != nil {
		return nil, fmt.Errorf("parse 'auth_uri': %w", err)
	}

	if l := len(web.RedirectUris); l != 1 {
		return nil, fmt.Errorf("wrong redirect url's length: expected 1, got %d", l)
	}

	query := auth.Query()
	query.Set("client_id", web.ClientID)
	query.Set("response_type", "code")
	query.Set("scope", "openid email")

	query.Set("redirect_uri", web.RedirectUris[0])

	web.AuthURL = auth.String() + "?" + query.Encode() + "&state="

	// token

	auth, err = url.Parse(web.TokenURL)
	if err != nil {
		return nil, fmt.Errorf("parse 'token_uri': %w", err)
	}

	query = auth.Query()
	query.Set("client_id", web.ClientID)
	query.Set("client_secret", web.ClientSecret)
	query.Set("redirect_uri", web.RedirectUris[0])
	query.Set("grant_type", "authorization_code")

	web.TokenURL = auth.String() + "?" + query.Encode() + "&code="

	// public keys
	f, err := parse.Key(web.CertsURL)
	if err != nil {
		return nil, fmt.Errorf("get jwt public keys: %w", err)
	}

	// return

	return &Auth{
		authURL:  web.AuthURL,
		getToken: web.TokenURL,
		Keys:     f,
	}, nil
}
