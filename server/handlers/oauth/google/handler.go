package google

import (
	"net/http"

	"go-auth/server/repository"
	"go-auth/utils/templates"
	"go.uber.org/zap"
)

type Methods struct {
	log *zap.Logger
	rep repository.Repository
	ts  *templates.Templates
}

func New(model repository.Model) http.Handler {
	router := http.NewServeMux()

	google := Methods{
		log: model.Log.Named("google"),
		rep: model.Rep,
		ts:  model.TS,
	}

	router.HandleFunc("/login", google.Login)
	router.HandleFunc("/callback", google.Callback)

	return router
}
