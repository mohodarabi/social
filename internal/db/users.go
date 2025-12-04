package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("duplicated user")
	ErrDuplicateUsername = errors.New("duplicated user")
)

type UserModel struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"passowrd"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type UserRepo struct {
	connection *sql.DB
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (user *UserRepo) Create(ctx context.Context, tx *sql.Tx, data *UserModel) error {
	query := `
		INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	err := tx.QueryRowContext(
		ctx,
		query,
		data.Username,
		data.Email,
		data.Password.hash,
	).Scan(
		&data.ID,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (user *UserRepo) GetById(ctx context.Context, userID int64) (*UserModel, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	var data UserModel
	err := user.connection.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&data.ID,
		&data.Username,
		&data.Email,
		&data.Password,
		&data.CreatedAt,
		&data.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &data, nil
}

func (user *UserRepo) CreateAndInvite(ctx context.Context, data *UserModel, token string, expTime time.Duration) error {
	return WithTx(user.connection, ctx, func(tx *sql.Tx) error {
		if err := user.Create(ctx, tx, data); err != nil {
			return err
		}

		if err := user.CreateUserInvitation(ctx, tx, data.ID, token, expTime); err != nil {
			return err
		}

		return nil
	})
}

func (user *UserRepo) CreateUserInvitation(ctx context.Context, tx *sql.Tx, userID int64, token string, expTime time.Duration) error {
	query := `
		INSERT INTO user_invitations (token, user_id, exptime) VALUES ($1, $2, $3)
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	_, err := tx.ExecContext(
		ctx,
		query,
		token,
		userID,
		time.Now().Add(expTime),
	)

	if err != nil {
		return err
	}

	return nil

}
