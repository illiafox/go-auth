package validator

import (
	"net/http"
	"net/mail"

	"auth-example/server/handlers/web/validator/pass"
	"auth-example/utils/templates"
)

func Validate(w http.ResponseWriter, r *http.Request, ts *templates.Templates) (email, password string, passed bool) {

	r.ParseForm()

	email = r.Form.Get("mail")
	if email == "" {
		ts.Message.Execute(w, "'mail' field is empty")

		return
	}

	password = r.Form.Get("password")
	if password == "" {
		ts.Message.Execute(w, "'password' field is empty")

		return
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		ts.Message.Execute(w, "wrong mail format: "+err.Error())

		return
	}

	err = pass.Validate(password)
	if err != nil {
		ts.Message.Execute(w, "password "+err.Error())

		return
	}

	return email, password, true
}
