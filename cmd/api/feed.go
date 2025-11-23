package main

import (
	"net/http"
)

func (app *application) getUserFeedHandler(response http.ResponseWriter, request *http.Request) {

	ctx := request.Context()

	feed, err := app.db.Posts.GetUserFeed(ctx, int64(1))

	if err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := writeJson(response, http.StatusCreated, feed); err != nil {
		app.badRequestError(response, request, err)
		return
	}

}
