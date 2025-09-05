package db

import (
	"context"
	"database/sql"
)

type DbRepo struct {
	Posts interface {
		Create(context.Context, *PostModel) error
	}

	Users interface {
		Create(context.Context, *UserModel) error
	}
}

func PostgresDb(connection *sql.DB) DbRepo {
	return DbRepo{
		Posts: &PostRepo{connection},
		Users: &UserRepo{connection},
	}
}
