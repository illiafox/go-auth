package mail

import (
	"gopkg.in/gomail.v2"
)

type Mail struct {
	Dialer *gomail.Dialer
	From   string
}

func (m Mail) Send(to, subject string, body string) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", m.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return m.Dialer.DialAndSend(msg)
}
