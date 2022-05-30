package cookie

import (
	"net/http"

	"github.com/gorilla/sessions"
	"go-auth/utils/config"
)

func New(conf config.Cookie) Cookie {

	store := sessions.NewCookieStore([]byte(conf.Key))
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteStrictMode
	store.Options.MaxAge = 60 * 60 * 24 * 28 // 28 days

	return Cookie{
		store: store,
		part:  conf.Part,
		name:  conf.Name,
	}
}

type Cookie struct {
	store      *sessions.CookieStore
	name, part string
}

func (c Cookie) Set(r *http.Request, w http.ResponseWriter, token string) error {
	session, err := c.store.Get(r, c.name)
	if err != nil {
		return err
	}

	session.Values[c.part] = token

	return sessions.Save(r, w)
}

func (c Cookie) Get(r *http.Request) (string, error) {
	session, err := c.store.Get(r, c.name)
	if err != nil {
		return "", err
	}

	token, ok := session.Values[c.part].(string)
	if !ok || token == "" {
		return "", nil
	}

	return token, nil
}

func (c Cookie) Delete(r *http.Request, w http.ResponseWriter) error {
	session, err := c.store.Get(r, c.name)
	if err != nil {
		return err
	}

	delete(session.Values, c.part)

	return sessions.Save(r, w)
}

func (c Cookie) GetDel(r *http.Request, w http.ResponseWriter) (string, error) {
	session, err := c.store.Get(r, c.name)
	if err != nil {
		return "", err
	}

	token, _ := session.Values[c.part].(string)

	delete(session.Values, c.part)

	return token, sessions.Save(r, w)
}
