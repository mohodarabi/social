package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("duplicated user")
	ErrDuplicateUsername = errors.New("duplicated user")
)

type UserModel struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"passowrd"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	IsAcitve  bool     `json:"is_active"`
}

type UserRepo struct {
	connection *sql.DB
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (user *UserRepo) Create(ctx context.Context, tx *sql.Tx, data *UserModel) error {
	query := `
		INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	err := tx.QueryRowContext(
		ctx,
		query,
		data.Username,
		data.Email,
		data.Password.hash,
	).Scan(
		&data.ID,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (user *UserRepo) GetById(ctx context.Context, userID int64) (*UserModel, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	var data UserModel
	err := user.connection.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&data.ID,
		&data.Username,
		&data.Email,
		&data.Password,
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

func (user *UserRepo) CreateAndInvite(ctx context.Context, data *UserModel, token string, expTime time.Duration) error {
	return WithTx(user.connection, ctx, func(tx *sql.Tx) error {
		if err := user.Create(ctx, tx, data); err != nil {
			return err
		}

		if err := user.CreateUserInvitation(ctx, tx, data.ID, token, expTime); err != nil {
			return err
		}

		return nil
	})
}

func (user *UserRepo) CreateUserInvitation(ctx context.Context, tx *sql.Tx, userID int64, token string, expTime time.Duration) error {
	query := `
		INSERT INTO user_invitations (token, user_id, exptime) VALUES ($1, $2, $3)
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	_, err := tx.ExecContext(
		ctx,
		query,
		token,
		userID,
		time.Now().Add(expTime),
	)

	if err != nil {
		return err
	}

	return nil

}

func (user *UserRepo) Activate(ctx context.Context, token string) error {
	return WithTx(user.connection, ctx, func(tx *sql.Tx) error {

		userData, err := user.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		userData.IsAcitve = true

		if err := user.update(ctx, tx, userData); err != nil {
			return err
		}

		if err := user.dateUserInvitation(ctx, tx, userData.ID); err != nil {
			return err
		}

		return nil
	})
}

func (user *UserRepo) update(ctx context.Context, tx *sql.Tx, userData *UserModel) error {
	query := `
		UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	_, err := tx.ExecContext(
		ctx,
		query,
		userData.Username,
		userData.Email,
		userData.IsAcitve,
		userData.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (user *UserRepo) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*UserModel, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active 
		FROM users u
		JOIN user_invitations ui on u.id = ui.user_id
		WHERE ui.token = $1 AND ui.exptime > $2
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	userData := &UserModel{}

	err := tx.QueryRowContext(
		ctx,
		query,
		hashToken,
		time.Now(),
	).Scan(
		&userData.ID,
		&userData.Username,
		&userData.Email,
		&userData.CreatedAt,
		&userData.IsAcitve,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return userData, nil
}

func (user *UserRepo) dateUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM user_invitations WHERE user_id = $1
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeOutSecond)
	defer cancle()

	_, err := tx.ExecContext(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}
