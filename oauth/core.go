package oauth

import (
	"fmt"

	"auth-example/oauth/github"
	"auth-example/oauth/google"
	"auth-example/oauth/parse"
	"auth-example/utils/config"
)

type Oauth struct {
	Google google.Auth
	Github github.Auth
}

func New(oauth config.Oauth) (*Oauth, error) {

	web, err := parse.Token(oauth.GooglePath)
	if err != nil {
		return nil, fmt.Errorf("parse google data (%s): %w", oauth.GooglePath, err)
	}

	g, err := google.New(web)
	if err != nil {
		return nil, fmt.Errorf("google: %w", err)
	}

	web, err = parse.Token(oauth.GitPath)
	if err != nil {
		return nil, fmt.Errorf("parse github data (%s): %w", oauth.GitPath, err)
	}

	hub, err := github.New(web)
	if err != nil {
		return nil, fmt.Errorf("github: %w", err)
	}

	return &Oauth{Google: *g, Github: *hub}, nil
}
