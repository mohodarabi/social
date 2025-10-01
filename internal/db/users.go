package db

import (
	"context"
	"database/sql"
	"errors"
)

type UserModel struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"passowrd"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserRepo struct {
	connection *sql.DB
}

func (user *UserRepo) Create(ctx context.Context, data *UserModel) error {

	query := `
		INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id
	`

	err := user.connection.QueryRowContext(
		ctx,
		query,
		data.Username,
		data.Email,
		data.Password,
	).Scan(
		&data.ID,
	)

	if err != nil {
		return err
	}

	return nil

}

func (user *UserRepo) GetById(ctx context.Context, userID int64) (*UserModel, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

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
