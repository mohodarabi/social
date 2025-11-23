package main

import (
	"net/http"
	"social/internal/db"
)

func (app *application) getUserFeedHandler(response http.ResponseWriter, request *http.Request) {

	feedQuery := db.PagnaitedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	feedQuery, err := feedQuery.Parse(*request)
	if err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := Validate.Struct(feedQuery); err != nil {
		app.badRequestError(response, request, err)
		return
	}

	ctx := request.Context()

	feed, err := app.db.Posts.GetUserFeed(ctx, int64(1), feedQuery)

	if err != nil {
		app.badRequestError(response, request, err)
		return
	}

	if err := writeJson(response, http.StatusCreated, feed); err != nil {
		app.badRequestError(response, request, err)
		return
	}

}
