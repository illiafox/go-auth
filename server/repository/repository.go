package repository

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"go-auth/database/model"
)

type Repository struct {
	// oauth states
	Memcached

	// store sessions
	Redis Redis

	// userdata
	Postgres Postgres

	Oauth Oauth

	Cookie interface {
		Set(r *http.Request, w http.ResponseWriter, token string) error
		Get(r *http.Request) (token string, err error)
		Delete(r *http.Request, w http.ResponseWriter) error
		GetDel(r *http.Request, w http.ResponseWriter) (token string, err error)
	}

	Mail interface {
		Send(to, subject string, body string) error
	}
}

type Memcached struct {
	State interface {
		Store(state string) (err error)
		Lookup(state string) (bool, error)
	}
	Mail interface {
		Store(key, mail string, secret []byte) error
		Get(key string) (mail, secret string, err error)
	}
}

type Redis struct {
	Session interface {
		New(token string, user int64) (err error)
		Get(token string) (int64, error)
		Delete(token string) error
		DeleteAll(user int64) error
	}
}

type Postgres struct {
	User interface {
		New(ctx context.Context, auth model.AuthType, mail, secret string) (err error)
		NewID(ctx context.Context, auth model.AuthType, mail, secret string) (id int64, err error)
		GetByMail(ctx context.Context, mail string) (userID int64, auth model.AuthType, err error)
		GetByID(ctx context.Context, id int64) (mail string, auth model.AuthType, err error)

		GetSecretByMail(ctx context.Context, mail string) (id int64, secret []byte, auth model.AuthType, err error)

		Exists(ctx context.Context, mail string) (exists bool, err error)
	}
}

type Oauth struct {
	Google Google
	Github Github
}

type Github interface {
	Auth(state string) (url string)
	Token(code string) (url string)
}

type Google interface {
	Github
	KeyFunc() jwt.Keyfunc
}
