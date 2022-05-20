package server

import (
	"net/http"
	"time"

	"auth-example/server/handlers/oauth"
	"auth-example/server/handlers/web"
	"auth-example/server/repository"
	"auth-example/utils/config"
	"github.com/gorilla/mux"
)

func New(conf config.Host, model repository.Model) *http.Server {
	router := mux.NewRouter()

	router.PathPrefix("/oauth/").Handler(http.StripPrefix(
		"/oauth", oauth.New(model),
	))

	router.PathPrefix("/").Handler(web.New(model))

	return &http.Server{
		Addr: "0.0.0.0:" + conf.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
}
