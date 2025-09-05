package db

import (
	"context"
	"database/sql"
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
