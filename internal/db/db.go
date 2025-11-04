package db

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

type DbRepo struct {
	Posts interface {
		Create(context.Context, *PostModel) error
		GetById(context.Context, int64) (*PostModel, error)
		DeleteById(context.Context, int64) error
		Update(context.Context, *PostModel) error
	}

	Users interface {
		Create(context.Context, *UserModel) error
		GetById(context.Context, int64) (*UserModel, error)
	}

	Comments interface {
		Create(context.Context, *CommentModel) error
		GetById(context.Context, int64) (*CommentModel, error)
		GetByPostId(context.Context, int64) ([]CommentModel, error)
	}
}

func PostgresDb(connection *sql.DB) DbRepo {
	return DbRepo{
		Posts:    &PostRepo{connection},
		Users:    &UserRepo{connection},
		Comments: &CommentRepo{connection},
	}
}
