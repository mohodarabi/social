package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type PostModel struct {
	ID        int64          `json:"id"`
	Content   string         `json:"content"`
	Title     string         `json:"title"`
	UserID    int64          `json:"user_id"`
	Tags      []string       `json:"tags"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Comments  []CommentModel `json:"comments"`
}

type PostRepo struct {
	connection *sql.DB
}

func (post *PostRepo) Create(ctx context.Context, data *PostModel) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := post.connection.QueryRowContext(
		ctx,
		query,
		data.Content,
		data.Title,
		data.UserID,
		pq.Array(data.Tags),
	).Scan(
		&data.ID,
		&data.CreatedAt,
		&data.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}

func (post *PostRepo) GetById(ctx context.Context, postID int64) (*PostModel, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags
		FROM posts 
		WHERE id = $1
	`

	var data PostModel
	err := post.connection.QueryRowContext(
		ctx,
		query,
		postID,
	).Scan(
		&data.ID,
		&data.UserID,
		&data.Title,
		&data.Content,
		&data.CreatedAt,
		&data.UpdatedAt,
		pq.Array(&data.Tags),
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

func (post *PostRepo) Update(ctx context.Context, data *PostModel) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2
		WHERE id = $3; 
	`
	result, err := post.connection.ExecContext(
		ctx,
		query,
		data.Title,
		data.Content,
		data.ID,
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

func (post *PostRepo) DeleteById(ctx context.Context, postID int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1; 
	`

	result, err := post.connection.ExecContext(ctx, query, postID)

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
