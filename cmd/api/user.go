package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"social/internal/db"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userCtxKey string

const userKey userCtxKey = "user"

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"Password" validate:"required,max=100"`
}

type FollowUserPayload struct {
	UserID int64 `json:"userId" validate:"required,max=100"`
}

func (app *application) createUserHandler(response http.ResponseWriter, request *http.Request) {

	var payload CreateUserPayload

	if err := readJson(response, request, &payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(response, request, err)
	}

	user := &db.UserModel{
		Email:    payload.Email,
		Password: payload.Password,
		Username: payload.Username,
	}

	ctx := request.Context()

	if err := app.db.Users.Create(ctx, user); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := writeJson(response, http.StatusCreated, user); err != nil {
		app.badRequestError(response, request, err)
		return
	}

}

func (app *application) getUserHandler(response http.ResponseWriter, request *http.Request) {
	user := getUserFromCtx(*request)
	if err := writeJson(response, http.StatusOK, user); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) followUserHandler(response http.ResponseWriter, request *http.Request) {
	followerId := getUserFromCtx(*request)

	var payload FollowUserPayload

	if err := readJson(response, request, &payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	err := app.db.Folowers.Follow(request.Context(), followerId.ID, payload.UserID)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			app.notFoundError(response, request, err)
		default:
			app.internalServerError(response, request, err)
		}
		return
	}

	responseDate := map[string]string{
		"message": "successfully added",
	}

	if err := writeJson(response, http.StatusOK, responseDate); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) unfollowUserHandler(response http.ResponseWriter, request *http.Request) {
	user := getUserFromCtx(*request)
	fmt.Println(user)

	if err := writeJson(response, http.StatusNoContent, nil); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) userContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		idParam := chi.URLParam(request, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestError(response, request, err)
			return
		}

		ctx := request.Context()

		user, err := app.db.Users.GetById(ctx, id)

		if err != nil {
			switch {
			case errors.Is(err, db.ErrNotFound):
				app.notFoundError(response, request, err)
			default:
				app.internalServerError(response, request, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userKey, user)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}

func getUserFromCtx(request http.Request) *db.UserModel {
	user, _ := request.Context().Value(userKey).(*db.UserModel)
	return user
}
