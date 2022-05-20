package redis

import (
	"fmt"

	"auth-example/database/redis/sessions"
	"auth-example/utils/config"
	"github.com/go-redis/redis"
)

type Redis struct {
	Session *sessions.Session
	closeDB func() error
}

func (r Redis) Close() error {
	return r.closeDB()
}

func New(conf config.Redis) (*Redis, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Pass,

		DB: conf.DB,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("redis: ping: %w", err)
	}

	return &Redis{
		Session: sessions.New(client),
		closeDB: client.Close,
	}, nil
}
