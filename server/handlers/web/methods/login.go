package methods

import (
	"context"
	"encoding/hex"
	"net/http"

	"auth-example/database/model"
	"auth-example/server/handlers/web/validator"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (m Methods) Login(w http.ResponseWriter, r *http.Request) {

	mail, pass, ok := validator.Validate(w, r, m.ts)
	if !ok {
		return
	}

	ctx := context.Background()

	id, secret, auth, err := m.rep.Postgres.User.GetSecretByMail(ctx, mail)
	if err != nil {
		m.log.Error("login: postgres: get secret by mail", zap.Error(err), zap.String("mail", mail))
		m.ts.Message.Internal(w)

		return
	}

	if secret == nil {
		m.ts.Message.Execute(w, `Account not found!<br>
			Please, <a href='/register'>Register</a>`)

		return

	}

	if auth != model.Password {
		m.ts.Message.Execute(w, "Account has another auth type, sign in via "+string(auth))

		return
	}

	err = bcrypt.CompareHashAndPassword(secret, []byte(pass))
	if err != nil {
		m.ts.Message.Execute(w, "Wrong password!")

		return
	}

	// // Delete old session

	session, err := m.rep.Cookie.Get(r)
	if err == nil && session != "" {
		err = m.rep.Redis.Session.Delete(session)
		if err != nil {
			m.log.Error("main: redis: delete session", zap.Error(err), zap.String("token", session))
			m.ts.Message.Internal(w)

			return
		}
	}

	// //

	key := hex.EncodeToString([]byte(uuid.NewString()))

	err = m.rep.Redis.Session.New(key, id)
	if err != nil {
		m.log.Error("login: redis: new session",
			zap.Error(err),
			zap.String("key", key),
			zap.Int64("id", id),
		)
		m.ts.Message.Internal(w)

		return
	}

	err = m.rep.Cookie.Set(r, w, key)
	if err != nil {
		m.ts.Message.Execute(w, "Cookie Error!<br><h4>"+err.Error()+"</h4>")

		return
	}

	m.ts.Message.Execute(w, `Redirecting to <a href='/'>
	Main Page </a> <meta http-equiv='refresh' content='2 url=/'>`)
}
