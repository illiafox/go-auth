package validator

import (
	"fmt"
	"net/http"
	"net/mail"

	"go-auth/server/handlers/web/validator/pass"
	"go-auth/utils/templates"
)

func Validate(w http.ResponseWriter, r *http.Request, ts *templates.Templates) (email, password string, passed bool) {

	err := r.ParseForm()
	if err != nil {
		ts.Message.Execute(w, fmt.Sprintf("parse form: %s", err))

	}

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

	_, err = mail.ParseAddress(email)
	if err != nil {
		ts.Message.Execute(w, fmt.Sprintf("wrong mail format: %s", err))

		return
	}

	err = pass.Validate(password)
	if err != nil {
		ts.Message.Execute(w, "password "+err.Error())

		return
	}

	return email, password, true
}
