package methods

import (
	"bytes"
	"context"
	"encoding/hex"
	"net/http"

	"auth-example/server/handlers/web/validator"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (m Methods) Register(w http.ResponseWriter, r *http.Request) {

	mail, pass, ok := validator.Validate(w, r, m.ts)
	if !ok {
		return
	}

	ctx := context.Background()

	exists, err := m.rep.Postgres.User.Exists(ctx, mail)
	if err != nil {
		m.log.Error("register: postgres: user exists", zap.Error(err), zap.String("mail", mail))
		m.ts.Message.Internal(w)

		return
	}

	if exists {
		m.ts.Message.Execute(w, "account with this mail already exists")

		return
	}

	secret, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		m.log.Error("register: generate password bcrypt", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}

	key := hex.EncodeToString([]byte(uuid.NewString()))

	buf := new(bytes.Buffer)
	m.ts.Mail.Execute(buf, key)

	err = m.rep.Mail.Send(mail, "Verify your Mail!", buf.String())
	if err != nil {
		m.log.Warn("register: smtp: send mail", zap.Error(err))
		m.ts.Message.Execute(w, "Couldn't send mail,try again later!")

		return
	}

	err = m.rep.Memcached.Mail.Store(key, mail, secret)
	if err != nil {
		m.log.Error("register: memcached: store mail", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}

	m.ts.Message.Execute(w, "Check your <a href='https://gmail.com'>Mail</a>")
}
