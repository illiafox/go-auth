package google

import (
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (m Methods) Login(w http.ResponseWriter, r *http.Request) {
	// Delete old session
	session, err := m.rep.Cookie.Get(r)
	if err == nil && session != "" {
		err = m.rep.Redis.Session.Delete(session)
		if err != nil {
			m.log.Error("redis: delete session", zap.Error(err), zap.String("token", session))
			m.ts.Message.Internal(w)

			return
		}
	}

	// //

	state := hex.EncodeToString([]byte(uuid.NewString()))

	err = m.rep.Memcached.State.Store(state)
	if err != nil {
		m.log.Error("memcached: store state", zap.Error(err), zap.String("state", state))
		m.ts.Message.Internal(w)

		return
	}

	http.Redirect(w, r, m.rep.Oauth.Google.Auth(state), http.StatusTemporaryRedirect)
}
