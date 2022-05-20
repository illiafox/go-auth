package user

import (
	"context"

	"auth-example/database/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	client *pgxpool.Pool
}

func New(client *pgxpool.Pool) *User {
	return &User{client: client}
}

func (u User) NewID(ctx context.Context, auth model.AuthType, mail, secret string) (int64, error) {
	var id int64 = -1
	err := u.client.QueryRow(ctx, "INSERT INTO users (mail,secret_type,secret) VALUES ($1,$2,$3) RETURNING user_id",
		mail, auth, secret).Scan(&id)

	return id, err
}

func (u User) New(ctx context.Context, auth model.AuthType, mail, secret string) error {
	_, err := u.client.Exec(ctx, "INSERT INTO users (mail,secret_type,secret) VALUES ($1,$2,$3)", mail, auth, secret)

	return err
}

func (u User) GetByMail(ctx context.Context, mail string) (userID int64, auth model.AuthType, err error) {

	err = u.client.QueryRow(ctx, "SELECT user_id,secret_type FROM users WHERE mail = $1", mail).
		Scan(&userID, &auth)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}

		return -1, "", err
	}

	return
}

func (u User) Exists(ctx context.Context, mail string) (exists bool, err error) {
	err = u.client.QueryRow(ctx, "SELECT exists(SELECT 1 FROM users WHERE mail = $1)", mail).
		Scan(&exists)

	return
}

func (u User) GetSecretByMail(ctx context.Context, mail string)(id int64,
	secret []byte, auth model.AuthType, err error) {

	err = u.client.QueryRow(ctx, "SELECT user_id,secret,secret_type FROM users WHERE mail = $1", mail).
		Scan(&id, &secret, &auth)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}

		return -1, nil, "", err
	}

	return
}

func (u User) GetByID(ctx context.Context, id int64) (mail string, auth model.AuthType, err error) {

	err = u.client.QueryRow(ctx, "SELECT mail,secret_type FROM users WHERE user_id = $1", id).
		Scan(&mail, &auth)

	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}

		return "", "", err
	}

	return
}
