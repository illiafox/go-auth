package google

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"auth-example/database/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const TryAgain = "Oauth error, try <a href='/oauth/google/login'>again</a>"

var client = http.Client{
	Timeout: time.Second * 5,
}

func (m Methods) Callback(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	state := query.Get("state")
	if state == "" {
		m.ts.Message.Execute(w, "state is not valid")

		return
	}



	// find state in cache
	exist, err := m.rep.Memcached.State.Lookup(state)
	if err != nil { // only internal
		m.log.Error("oauth: google: memcached: lookup state",
			zap.Error(err),
			zap.String("state", state),
		)
		m.ts.Message.Internal(w)

		return
	}

	if !exist {
		m.ts.Message.Execute(w, "state is not valid")

		return
	}

	// //

	if problem := query.Get("error"); problem != "" {
		m.log.Warn("oauth: google: error inside query",
			zap.String("error", problem),
		)

		m.ts.Message.Internal(w)

		return
	}

	// //

	// get code
	code := query.Get("code")
	if code == "" {
		m.ts.Message.Execute(w, "code not found")

		return
	}

	// //

	resp, err := client.Post(m.rep.Oauth.Google.Token(code), "", nil)
	if err != nil {
		m.log.Error("oauth: google: get token request", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}
	defer resp.Body.Close()

	var data = struct {
		AccessToken string `json:"access_token"`
		JWT         string `json:"id_token"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		m.log.Error("oauth: google: decode json token request", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}

	if data.JWT == "" || data.AccessToken == "" {
		m.log.Warn("oauth: google: data fields are empty")
		m.ts.Message.Execute(w, TryAgain)

		return
	}

	token, err := jwt.Parse(data.JWT, m.rep.Oauth.Google.Keyfunc)
	if err != nil {
		m.log.Warn("oauth: google: jwt parse", zap.Error(err))
		m.ts.Message.Execute(w, TryAgain)

		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		m.log.Warn("oauth: google: cast to MapClaims")
		m.ts.Message.Execute(w, TryAgain)

		return
	}

	verified, ok := claims["email_verified"].(bool)
	if !ok {
		m.log.Warn("oauth: google: cast to 'email_verified'")
		m.ts.Message.Execute(w, TryAgain)

		return
	}
	mail, ok := claims["email"].(string)
	if !ok {
		m.log.Warn("oauth: google: cast to 'email'")
		m.ts.Message.Execute(w, TryAgain)

		return
	}

	if mail == "" || !verified {
		m.ts.Message.Execute(w, "verified or primary mail not found")

		return
	}

	ctx := context.Background()

	id, auth, err := m.rep.Postgres.User.GetByMail(ctx, mail)
	if err != nil {
		m.log.Error("oauth: google: postgresql: get user",
			zap.Error(err),
			zap.String("mail", mail),
		)
		m.ts.Message.Internal(w)

		return
	}

	if id < 0 {

		id, err = m.rep.Postgres.User.NewID(ctx, model.Google, mail, data.AccessToken)
		if err != nil {
			m.log.Error("oauth: google: postgresql: new user",
				zap.Error(err),
				zap.String("mail", mail),
			)
			m.ts.Message.Internal(w)

			return
		}

	} else if auth != model.Github {
		m.ts.Message.Execute(w, "Account has another auth type, sign in via "+string(auth))

		return
	}

	key := hex.EncodeToString([]byte(uuid.NewString()))

	err = m.rep.Redis.Session.New(key, id)
	if err != nil {
		m.log.Error("oauth: google: redis: new session",
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
