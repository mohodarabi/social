package db

import (
	"context"
	"database/sql"
)

type PostRepo struct {
	connection *sql.DB
}

func (post *PostRepo) Create(ctx context.Context) error {
	return nil
}
