package db

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound            = errors.New("record not found")
	ErrUserAlreadyFollowed = errors.New("user already followed")
	QueryTimeOutSecond     = time.Second * 5
)

type DbRepo struct {
	Posts interface {
		Create(context.Context, *PostModel) error
		GetById(context.Context, int64) (*PostModel, error)
		DeleteById(context.Context, int64) error
		Update(context.Context, *PostModel) error
		GetUserFeed(context.Context, int64, PagnaitedFeedQuery) ([]PostWithMetadata, error)
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

	Folowers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		UnFollow(ctx context.Context, followerID, userID int64) error
	}
}

func PostgresDb(connection *sql.DB) DbRepo {
	return DbRepo{
		Posts:    &PostRepo{connection},
		Users:    &UserRepo{connection},
		Comments: &CommentRepo{connection},
		Folowers: &FollowerRepo{connection},
	}
}
