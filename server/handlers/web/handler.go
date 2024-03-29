package web

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"go-auth/server/handlers/web/methods"
	"go-auth/server/repository"
)

func New(model repository.Model) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	model.Log = model.Log.Named("web")

	static := http.FileServer(http.Dir(path.Dir(model.Config.Static)))

	router.PathPrefix("/images/").Handler(static)
	router.PathPrefix("/js/").Handler(static)

	router.Handle("/favicon.ico", static)

	m := methods.New(model)

	router.HandleFunc("/", m.Main)

	router.HandleFunc("/register/", m.Register).Methods(http.MethodPost)
	router.Handle("/register/", static).Methods(http.MethodGet)
	router.PathPrefix("/register/").Handler(static).Methods(http.MethodGet)

	router.HandleFunc("/login/", m.Login).Methods(http.MethodPost)
	router.Handle("/login/", static).Methods(http.MethodGet)
	router.PathPrefix("/login/").Handler(static).Methods(http.MethodGet)

	router.HandleFunc("/verify", m.Verify)
	router.HandleFunc("/logout", m.Logout)

	return router
}
