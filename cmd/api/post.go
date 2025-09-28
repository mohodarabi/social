package main

import (
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
		UserID:  2,
	}

	ctx := request.Context()

	if err := app.db.Posts.Create(ctx, post); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := writeJson(response, http.StatusCreated, post); err != nil {
		app.badRequestError(response, request, err)
		return
	}

}

func (app *application) getPostHandler(response http.ResponseWriter, request *http.Request) {
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

	if err := writeJson(response, http.StatusCreated, post); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}
