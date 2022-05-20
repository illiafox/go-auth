package oauth

import (
	"net/http"

	"auth-example/server/handlers/oauth/github"
	google "auth-example/server/handlers/oauth/google"
	"auth-example/server/repository"
)

func New(model repository.Model) http.Handler {
	root := http.NewServeMux()

	root.Handle("/google/", http.StripPrefix("/google", google.New(model)))
	root.Handle("/github/", http.StripPrefix("/github", github.New(model)))

	return root
}
