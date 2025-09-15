package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type PostModel struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
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
