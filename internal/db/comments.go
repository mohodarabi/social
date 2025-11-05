package db

import (
	"context"
	"database/sql"
	"errors"
)

type CommentModel struct {
	ID        int64     `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	User      UserModel `json:"user"`
}

type CommentRepo struct {
	connection *sql.DB
}

func (comment *CommentRepo) Create(ctx context.Context, data *CommentModel) error {
	query := `
		INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

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

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

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

func (comment *CommentRepo) GetByPostId(ctx context.Context, postID int64) ([]CommentModel, error) {
	query := `
		SELECT comments.id, comments.post_id, comments.user_id, comments.content, comments.created_at, users.username, users.id FROM comments
		JOIN users on users.id = comments.user_id
		WHERE comments.post_id = $1
		ORDER BY comments.created_at DESC;
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	rows, err := comment.connection.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []CommentModel{}
	for rows.Next() {
		var comment CommentModel
		comment.User = UserModel{}
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.Username,
			&comment.User.ID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
