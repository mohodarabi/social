package db

import (
	"context"
	"database/sql"
	"errors"
)

type FollowerModel struct {
	ID         int64  `json:"id"`
	UserID     string `json:"user_id"`
	FollowerID string `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type FollowerRepo struct {
	connection *sql.DB
}

func (follower *FollowerRepo) GetById(ctx context.Context, userID int64) (*FollowerModel, error) {
	query := `
		SELECT id, user_id, follower_id, created_at, updated_at
		FROM followers 
		WHERE id = $1
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	var data FollowerModel
	err := follower.connection.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&data.ID,
		&data.UserID,
		&data.FollowerID,
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
