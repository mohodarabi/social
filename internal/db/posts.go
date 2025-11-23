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
	Version   int            `json:"version"`
	Comments  []CommentModel `json:"comments"`
	User      UserModel      `json:"user"`
}

type PostWithMetadata struct {
	PostModel
	CommentsCount int `json:"comments_count"`
}

type PostRepo struct {
	connection *sql.DB
}

func (post *PostRepo) GetUserFeed(ctx context.Context, userID int64, feedQuery PagnaitedFeedQuery) ([]PostWithMetadata, error) {
	query := `
		SELECT
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON follower_id = p.user_id OR p.user_id = $1
		WHERE f.user_id = $1 or p.user_id  = $1
		GROUP BY p.id, u.username
		ORDER BY ` + feedQuery.Sort + `
		LIMIT $2 OFFSET $3
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	rows, err := post.connection.QueryContext(
		ctx,
		query,
		userID,
		feedQuery.Limit,
		feedQuery.Offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata

	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)

		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil

}

func (post *PostRepo) Create(ctx context.Context, data *PostModel) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

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
		SELECT id, user_id, title, content, created_at, updated_at, tags, version
		FROM posts 
		WHERE id = $1
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

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
		&data.Version,
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
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4; 
		RETURNING version
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	err := post.connection.QueryRowContext(
		ctx,
		query,
		data.Title,
		data.Content,
		data.ID,
		data.Version,
	).Scan(
		&data.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (post *PostRepo) DeleteById(ctx context.Context, postID int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1; 
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

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
