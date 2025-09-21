package main

import (
	"errors"
	"net/http"
	"social/internal/db"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(response http.ResponseWriter, request *http.Request) {

	var payload CreatePostPayload

	if err := readJson(response, request, &payload); err != nil {
		writeJson(response, http.StatusBadRequest, err.Error())
		return
	}

	post := &db.PostModel{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  2,
	}

	ctx := request.Context()

	if err := app.db.Posts.Create(ctx, post); err != nil {
		writeJson(response, http.StatusBadRequest, err.Error())
		return
	}

	if err := writeJson(response, http.StatusCreated, post); err != nil {
		writeJson(response, http.StatusBadRequest, err.Error())
		return
	}

}

func (app *application) getPostHandler(response http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJson(response, http.StatusBadRequest, err.Error())
		return
	}

	ctx := request.Context()

	post, err := app.db.Posts.GetById(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrNotFound):
			writeJson(response, http.StatusNotFound, err.Error())
		default:
			writeJson(response, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJson(response, http.StatusCreated, post); err != nil {
		writeJson(response, http.StatusInternalServerError, err.Error())
		return
	}
}
