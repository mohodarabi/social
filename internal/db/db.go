package db

import (
	"context"
	"database/sql"
)

type DbRepo struct {
	Posts interface {
		Create(context.Context) error
	}

	Users interface {
		Create(context.Context) error
	}
}

func PostgresDb(connection *sql.DB) DbRepo {
	return DbRepo{
		Posts: &PostRepo{connection},
		Users: &UserRepo{connection},
	}
}
