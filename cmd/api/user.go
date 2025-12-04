package main

import (
	"context"
	"errors"
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

// GetUser godoc
//
//	@Summary		get user
//	@Description	get user by id, id should send in params
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"userID"
//	@Success		200	{object}	db.UserModel
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(response http.ResponseWriter, request *http.Request) {
	user := getUserFromCtx(*request)
	if err := writeJson(response, http.StatusOK, user); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) followUserHandler(response http.ResponseWriter, request *http.Request) {
	follower := getUserFromCtx(*request)

	var payload FollowUserPayload

	if err := readJson(response, request, &payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	err := app.db.Folowers.Follow(request.Context(), follower.ID, payload.UserID)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			app.notFoundError(response, request, err)
		case errors.Is(err, db.ErrUserAlreadyFollowed):
			app.badRequestError(response, request, err)
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
	follower := getUserFromCtx(*request)

	var payload FollowUserPayload

	if err := readJson(response, request, &payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	err := app.db.Folowers.UnFollow(request.Context(), follower.ID, payload.UserID)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			app.notFoundError(response, request, err)
		case errors.Is(err, db.ErrUserAlreadyFollowed):
			app.badRequestError(response, request, err)
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
