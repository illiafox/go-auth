package mail

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"go-auth/database/memcached/mail/base64"
)

type Mail struct {
	client *memcache.Client
}

const (
	prefix = "mail:"

	expire = 60 * 20 // 20 minutes
)

/*
STORING MODEL
//
1. KEY
mail:38613561363365622d626664662d343631372d393936612d363964346361333765313661
//
2. VALUE:
mail:type:secret
*/

func New(client *memcache.Client) Mail {
	return Mail{client: client}
}

func (m Mail) Store(key, mail string, secret []byte) error {

	data, err := base64.Encode(mail, secret)
	if err != nil {
		return fmt.Errorf("base64: encode: %w", err)
	}

	err = m.client.Set(&memcache.Item{
		Key:        prefix + key,
		Value:      data,
		Flags:      0,
		Expiration: expire,
	})

	if err != nil {
		return fmt.Errorf("set: %w", err)
	}

	return nil
}

func (m Mail) Get(key string) (mail, secret string, err error) {
	key = prefix + key

	item, err := m.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss { // Not Found
			err = nil
		}

		return "", "", err
	}

	err = m.client.Delete(key)
	if err != nil {
		if err == memcache.ErrCacheMiss { // Not Found
			err = nil
		}

		return "", "", err
	}

	mail, secret, err = base64.Decode(item.Value)
	if err != nil {
		return "", "", fmt.Errorf("base64: decode: %w", err)
	}

	return
}
