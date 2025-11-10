package db

import (
	"context"
	"database/sql"
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

func (follower *FollowerRepo) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2); 
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	result, err := follower.connection.ExecContext(
		ctx,
		query,
		userID,
		followerID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (follower *FollowerRepo) UnFollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2); 
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	result, err := follower.connection.ExecContext(
		ctx,
		query,
		userID,
		followerID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
