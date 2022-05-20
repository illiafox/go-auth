package sessions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type Session struct {
	client *redis.Client
}

func New(client *redis.Client) *Session {
	return &Session{client: client}
}

const (
	expire = time.Hour * 24 * 30 // 30 days
)

func (r Session) New(token string, user int64) error {

	return r.client.Set(strconv.FormatInt(user, 10)+":"+token, nil, expire).Err()
}

// Get returns user id (-1 if not exists) and error
func (r Session) Get(token string) (int64, error) {
	it := r.client.Scan(0, "*:"+token, 1)

	if err := it.Err(); err != nil {
		return -1, err
	}

	val, _ := it.Val()

	switch l := len(val); l {
	case 0:

		return -1, nil

	case 1:
		v := val[0]
		i := strings.IndexByte(v, ':')

		if i < 1 {
			return int64(i), fmt.Errorf("wrong id format: %s", v)
		}

		return strconv.ParseInt(v[:i], 10, 64)

	default:

		return -1, fmt.Errorf("wrong keys length: expected 1, got %d", l)
	}
}

func (r Session) Delete(token string) error {
	it := r.client.Scan(0, "*:"+token, 1)

	if err := it.Err(); err != nil {
		return err
	}

	val, _ := it.Val()

	switch l := len(val); l {

	case 1:
		return r.client.Del(val[0]).Err()

	case 0:

		return nil

	default:

		return fmt.Errorf("wrong keys length: expected 1, got %d", l)
	}
}

func (r Session) DeleteAll(user int64) error {

	it := r.client.Scan(0, strconv.FormatInt(user, 10)+":*", 1)

	if err := it.Err(); err != nil {
		return err
	}

	keys, _ := it.Val()

	return r.client.Del(keys...).Err()

}
