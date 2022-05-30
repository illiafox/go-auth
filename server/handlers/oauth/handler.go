package oauth

import (
	"net/http"

	"go-auth/server/handlers/oauth/github"
	google "go-auth/server/handlers/oauth/google"
	"go-auth/server/repository"
)

func New(model repository.Model) http.Handler {
	root := http.NewServeMux()

	model.Log = model.Log.Named("oauth")

	root.Handle("/google/", http.StripPrefix("/google", google.New(model)))
	root.Handle("/github/", http.StripPrefix("/github", github.New(model)))

	return root
}
