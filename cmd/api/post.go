package main

import (
	"net/http"
	"social/internal/db"
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
