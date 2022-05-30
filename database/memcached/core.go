package memcached

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"go-auth/database/memcached/mail"
	"go-auth/database/memcached/state"
	"go-auth/utils/config"
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
