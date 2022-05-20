package database

import (
	"context"
	"fmt"

	"auth-example/database/memcached"
	"auth-example/database/pg"
	"auth-example/database/redis"
	"auth-example/utils/config"
	"go.uber.org/zap"
)

type Database struct {
	Postgres *pg.Postgres

	Redis *redis.Redis

	Memcached *memcached.Memcached

	closeFunc func() error
}

func (d Database) Close(logger *zap.Logger) {
	logger.Info("Closing database")
	if err := d.closeFunc(); err != nil {
		logger.Error("Error", zap.Error(err))
	}
}

func New(ctx context.Context, conf *config.Config) (*Database, error) {
	post, err := pg.New(ctx, conf.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	rdb, err := redis.New(conf.Redis)
	if err != nil {
		post.Close()

		return nil, fmt.Errorf("redis: %w", err)
	}

	mem, err := memcached.New(conf.Memcached)
	if err != nil {
		post.Close()

		err = fmt.Errorf("memcached: %w", err)

		if r := rdb.Close(); r != nil {
			err = fmt.Errorf("%w || closing redis: %s", err, r)
		}
		
		return nil, err

	}

	return &Database{
		Postgres:  post,
		Redis:     rdb,
		Memcached: mem,
		//
		closeFunc: func() error {
			post.Close()

			return rdb.Close()
		},
	}, nil
}
