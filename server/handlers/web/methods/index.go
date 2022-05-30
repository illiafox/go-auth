package methods

import (
	"context"
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

func New(model repository.Model) Methods {
	return Methods{
		log: model.Log,
		rep: model.Rep,
		ts:  model.TS,
	}
}

func (m Methods) Main(w http.ResponseWriter, r *http.Request) {
	session, err := m.rep.Cookie.Get(r)
	if err != nil {
		m.ts.Message.Execute(w, "Cookie Error!<br><h4>"+err.Error()+"</h4>")

		return
	}

	if session == "" {
		m.ts.Message.Execute(w, `Session not found!<br>
		Please, <a href='/login'>Login</a> or <a href='/register'>Register</a>`)

		return
	}

	id, err := m.rep.Redis.Session.Get(session)
	if err != nil {
		m.log.Error("main: redis: get session", zap.Error(err), zap.String("token", session))
		m.ts.Message.Internal(w)

		return
	}

	if id < 0 {
		m.ts.Message.Execute(w, `Session is not valid!<br>
		Please, <a href='/login'>Login</a> again`)

		return
	}

	ctx := context.Background()

	mail, auth, err := m.rep.Postgres.User.GetByID(ctx, id)
	if err != nil {
		m.log.Error("main: memcached: get user by id", zap.Error(err), zap.Int64("id", id))
		m.ts.Message.Internal(w)

		return
	}

	if mail == "" {
		m.ts.Message.Execute(w, `Account not found!<br>
		Please, <a href='/login'>Login</a> again`)

		return
	}

	err = m.ts.Main.Any(w, struct {
		Mail, Auth string
	}{
		Mail: mail,
		Auth: string(auth),
	})

	if err != nil {
		m.log.Error("main: execute template", zap.Error(err), zap.String("mail", mail))
		m.ts.Message.Internal(w)

		return
	}
}
