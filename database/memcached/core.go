package memcached

import (
	"fmt"

	"auth-example/database/memcached/mail"
	"auth-example/database/memcached/state"
	"auth-example/utils/config"
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcached struct {
	State state.State
	Mail  mail.Mail
}

func New(conf config.Memcached) (*Memcached, error) {
	mc := memcache.New(conf.Address)

	if err := mc.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Memcached{
		State: state.New(mc),
		Mail:  mail.New(mc),
	}, nil
}
