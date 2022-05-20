package parse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

type Web struct {
	ClientID                string   `json:"client_id"`
	ProjectID               string   `json:"project_id"`
	AuthURL                 string   `json:"auth_uri"`
	TokenURL                string   `json:"token_uri"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
	CertsURL                string   `json:"certs_uri"`
}

type JSON struct {
	Web Web
}

var client = http.Client{Timeout: time.Second * 5}

func Key(url string) (jwt.Keyfunc, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	js, err := keyfunc.NewJSON(data)
	if err != nil {
		return nil, fmt.Errorf("create keyfunc: %w", err)
	}

	return js.Keyfunc, nil
}

func Token(file string) (*Web, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var web JSON

	err = json.Unmarshal(data, &web)
	if err != nil {
		return nil, fmt.Errorf("unmarshall json: %w", err)
	}

	return &web.Web, nil
}
