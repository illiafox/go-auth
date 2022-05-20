package pg

import (
	"context"
	"fmt"

	"auth-example/database/pg/user"
	"auth-example/utils/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	User    *user.User
	closeDB func()
}

func (p Postgres) Close() {
	p.closeDB()
}

func New(ctx context.Context, conf config.Postgres) (*Postgres, error) {

	pool, err := pgxpool.Connect(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%v:%v/%v?sslmode=disable",
			conf.User,
			conf.Pass,
			conf.IP,
			conf.Port,
			conf.DBName,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Postgres{
		User:    user.New(pool),
		closeDB: pool.Close,
	}, nil
}
