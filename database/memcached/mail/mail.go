package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type Mail struct {
	client *memcache.Client
}

const (
	prefix = "mail:"

	expire = 60 * 20 // 20 minutes
)

var separator = []byte(":")

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
	buf := new(bytes.Buffer)

	encoder := base64.NewEncoder(base64.StdEncoding, buf)

	_, err := encoder.Write([]byte(mail))
	if err != nil {
		return fmt.Errorf("encode mail: %w", err)
	}

	buf.Write(separator)

	_, err = encoder.Write(secret)
	if err != nil {
		return fmt.Errorf("encode secret: %w", err)
	}

	err = encoder.Close()
	if err != nil {
		return fmt.Errorf("close encoder: %w", err)
	}

	return m.client.Set(&memcache.Item{
		Key:        prefix + key,
		Value:      buf.Bytes(),
		Flags:      0,
		Expiration: expire,
	})
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

	i := bytes.Index(item.Value, separator)
	if i < 0 {
		return "", "", fmt.Errorf("wrong separator index %d", i)
	}
	//
	buf := item.Value[:i]
	_, err = base64.StdEncoding.Decode(buf, buf)
	if err != nil {
		return "", "", fmt.Errorf("base64: decode first part: %w", err)
	}
	mail = string(buf[:base64.StdEncoding.DecodedLen(len(buf))])
	//
	buf = item.Value[i+1:]
	_, err = base64.StdEncoding.Decode(buf, buf)
	if err != nil {
		return "", "", fmt.Errorf("base64: decode second part: %w", err)
	}
	secret = string(buf[:base64.StdEncoding.DecodedLen(len(buf))])
	//
	return mail, secret, nil
}
