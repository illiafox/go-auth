package github

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go-auth/database/model"
	"go.uber.org/zap"
)

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

	// //

	// find state in cache
	exist, err := m.rep.Memcached.State.Lookup(state)
	if err != nil { // only internal
		m.log.Error("memcached: lookup state",
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
		if problem == "access_denied" {
			m.ts.Message.Execute(w, "Access denied, <a href='/oauth/github/login'>login</a> again")

			return
		}

		desc := query.Get("error_description")
		uri := query.Get("error_uri")

		m.log.Warn("error inside query",
			zap.String("error", problem),
			zap.String("description", desc),
			zap.String("url", uri),
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

	//

	resp, err := client.Get(m.rep.Oauth.Github.Token(code))
	if err != nil {
		m.log.Error("get token request", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}
	defer resp.Body.Close()

	data := new(strings.Builder)
	_, err = io.Copy(data, resp.Body)
	if err != nil {
		m.log.Error("read body", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}

	values, err := url.ParseQuery(data.String())
	if err != nil {
		m.log.Error("parse query",
			zap.Error(err),
			zap.String("query", data.String()),
		)
		m.ts.Message.Internal(w)

		return
	}

	token := values.Get("access_token")
	if token == "" {
		m.log.Error("token",
			zap.String("error", values.Get("error_description")),
			zap.String("query", data.String()),
		)
		m.ts.Message.Internal(w)

		return
	}

	//	request, err := http.NewRequest("GET", m.rep.Oauth.Github.TokenURL(), nil)
	request, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)

	if err != nil {
		m.log.Error("new user data request", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(request)
	if err != nil {
		m.log.Error("do user data request", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}
	defer resp.Body.Close()

	var mails []struct {
		Email      string  `json:"email"`
		Primary    bool    `json:"primary"`
		Verified   bool    `json:"verified"`
		Visibility *string `json:"visibility"`
	}

	err = json.NewDecoder(resp.Body).Decode(&mails)
	if err != nil {
		m.log.Error("user data: decode", zap.Error(err))
		m.ts.Message.Internal(w)

		return
	}

	var mail string
	for _, m := range mails {
		if m.Primary && m.Verified {
			mail = m.Email
		}
	}

	if mail == "" {
		m.ts.Message.Execute(w, "verified or primary mail not found")

		return
	}

	ctx := context.Background()

	id, auth, err := m.rep.Postgres.User.GetByMail(ctx, mail)
	if err != nil {
		m.log.Error("postgresql: get user",
			zap.Error(err),
			zap.String("mail", mail),
		)
		m.ts.Message.Internal(w)

		return
	}

	if id < 0 {

		id, err = m.rep.Postgres.User.NewID(ctx, model.Github, mail, token)
		if err != nil {
			m.log.Error("postgresql: new user",
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
		m.log.Error("redis: new session",
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
