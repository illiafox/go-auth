package methods

import (
	"net/http"

	"go.uber.org/zap"
)

func (m Methods) Logout(w http.ResponseWriter, r *http.Request) {

	token, err := m.rep.Cookie.GetDel(r, w)
	if err != nil {
		m.ts.Message.Execute(w, "Cookie Error!<br><h4>"+err.Error()+"</h4>")

		return
	}

	if token != "" {
		err = m.rep.Redis.Session.Delete(token)
		if err != nil {
			m.log.Error("logout: redis: delete session", zap.Error(err), zap.String("token", token))
		}
	}

	m.ts.Message.Execute(w, `Logged out<br>
	Redirecting to <a href='/login'> login </a> page
	<meta http-equiv='refresh' content='2 url=/login'>`)
}
