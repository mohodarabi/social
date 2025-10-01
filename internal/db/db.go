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
	}

	Users interface {
		Create(context.Context, *UserModel) error
		GetById(context.Context, int64) (*UserModel, error)
	}
}

func PostgresDb(connection *sql.DB) DbRepo {
	return DbRepo{
		Posts: &PostRepo{connection},
		Users: &UserRepo{connection},
	}
}
