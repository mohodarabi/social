package db

import (
	"context"
	"database/sql"
	"errors"
)

type CommentModel struct {
	ID        int64  `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CommentRepo struct {
	connection *sql.DB
}

func (comment *CommentRepo) Create(ctx context.Context, data *CommentModel) error {

	query := `
		INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id
	`

	err := comment.connection.QueryRowContext(
		ctx,
		query,
		data.PostID,
		data.UserID,
		data.Content,
	).Scan(
		&data.ID,
	)

	if err != nil {
		return err
	}

	return nil

}

func (comment *CommentRepo) GetById(ctx context.Context, commentID int64) (*CommentModel, error) {
	query := `
		SELECT id, post_id, user_id, content, created_at, updated_at
		FROM comments 
		WHERE id = $1
	`

	var data CommentModel
	err := comment.connection.QueryRowContext(
		ctx,
		query,
		commentID,
	).Scan(
		&data.ID,
		&data.PostID,
		&data.UserID,
		&data.Content,
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
