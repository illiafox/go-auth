package mail

import (
	"fmt"

	"go-auth/utils/config"
	"gopkg.in/gomail.v2"
)

func NewMail(conf config.SMTP) (Mail, error) {
	m := Mail{
		Dialer: gomail.NewDialer(conf.Hostname, conf.Port, conf.Mail, conf.Password),
		From:   conf.Mail,
	}

	// Connect to the remote SMTP server.
	rc, err := m.Dialer.Dial()
	if err != nil {
		return Mail{}, fmt.Errorf("dial: %w", err)
	}

	err = rc.Close()
	if err != nil {
		return Mail{}, fmt.Errorf("close test dial: %w", err)
	}

	return m, nil
}
