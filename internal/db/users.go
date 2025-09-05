package db

import (
	"context"
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func (user *UserRepo) Create(ctx context.Context) error {
	return nil
}
