package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/db"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=100"`
}

type postCtxKey string
const postKey postCtxKey = "post"

func (app *application) createPostHandler(response http.ResponseWriter, request *http.Request) {

	var payload CreatePostPayload

	if err := readJson(response, request, &payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(response, request, err)
	}

	post := &db.PostModel{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := request.Context()

	if err := app.db.Posts.Create(ctx, post); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := writeJsonResponse(response, http.StatusCreated, post); err != nil {
		app.badRequestError(response, request, err)
		return
	}

}

func (app *application) getPostHandler(response http.ResponseWriter, request *http.Request) {
	post := getPostFromCtx(*request)

	comments, err := app.db.Comments.GetByPostId(request.Context(), post.ID)
	if err != nil {
		app.internalServerError(response, request, err)
		return
	}

	post.Comments = comments

	if err := writeJsonResponse(response, http.StatusCreated, post); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) deletePostHandler(response http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "postID")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.badRequestError(response, request, err)
		return
	}

	ctx := request.Context()

	err = app.db.Posts.DeleteById(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			app.notFoundError(response, request, err)
		default:
			app.internalServerError(response, request, err)
		}
		return
	}

	if err := writeJsonResponse(response, http.StatusCreated, "post successfully deleted"); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) updatePostHandler(response http.ResponseWriter, request *http.Request) {
	post := getPostFromCtx(*request)

	var payload UpdatePostPayload

	if err := readJson(response, request, &payload); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(response, request, err)
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	err := app.db.Posts.Update(request.Context(), post)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			app.notFoundError(response, request, err)
		default:
			app.internalServerError(response, request, err)
		}
		return
	}

	if err := writeJsonResponse(response, http.StatusCreated, post); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) postContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		idParam := chi.URLParam(request, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.badRequestError(response, request, err)
			return
		}

		ctx := request.Context()

		post, err := app.db.Posts.GetById(ctx, id)

		if err != nil {
			switch {
			case errors.Is(err, db.ErrNotFound):
				app.notFoundError(response, request, err)
			default:
				app.internalServerError(response, request, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postKey, post)
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}

func getPostFromCtx(request http.Request) *db.PostModel {
	post, _ := request.Context().Value(postKey).(*db.PostModel)
	return post
}
