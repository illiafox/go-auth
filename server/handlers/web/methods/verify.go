package methods

import (
	"context"
	"net/http"

	"auth-example/database/model"
	"go.uber.org/zap"
)

func (m Methods) Verify(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		m.ts.Message.Execute(w, "query field 'key' is empty")

		return
	}

	mail, secret, err := m.rep.Memcached.Mail.Get(key)
	if err != nil {
		m.log.Error("verify: memcached: get verify data", zap.Error(err), zap.String("key", key))

		m.ts.Message.Internal(w)

		return
	}

	ctx := context.Background()

	exists, err := m.rep.Postgres.User.Exists(ctx, mail)
	if err != nil {
		m.log.Error("verify: postgres: user exists", zap.Error(err), zap.String("mail", mail))
		m.ts.Message.Internal(w)

		return
	}

	if exists {
		m.ts.Message.Execute(w, "account with this mail already exists")

		return
	}

	err = m.rep.Postgres.User.New(ctx, model.Password, mail, secret)
	if err != nil {
		m.log.Error("verify: postgres: new user", zap.Error(err), zap.String("mail", mail))
		m.ts.Message.Internal(w)

		return
	}

	m.ts.Message.Execute(w, "Account is verified<br>Please, <a href='/login'>login</a>")
}
