package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"social/internal/db"

	"github.com/google/uuid"
)

type UserAuthenticate struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=20"`
}

type UserWithToken struct {
	User  *db.UserModel
	Token string `json:"token"`
}

// CreateUser godoc
//
//	@Summary		create user and invitation
//	@Description	create user and invitation
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		UserAuthenticate	true	"User authentication payload"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/user [post]
func (app *application) authHandler(response http.ResponseWriter, request *http.Request) {
	var payload UserAuthenticate

	if err := readJson(response, request, &payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	ctx := request.Context()

	user := &db.UserModel{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(response, request, err)
		return
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.db.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.expTime)

	if err != nil {
		switch err {
		case db.ErrDuplicateEmail:
			app.badRequestError(response, request, err)
			return
		case db.ErrDuplicateUsername:
			app.badRequestError(response, request, err)
			return
		default:
			app.internalServerError(response, request, err)
			return
		}
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := writeJson(response, http.StatusOK, userWithToken); err != nil {
		app.internalServerError(response, request, err)
		return
	}

}
