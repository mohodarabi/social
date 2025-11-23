package db

import (
	"net/http"
	"strconv"
)

type PagnaitedFeedQuery struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (feedQuery PagnaitedFeedQuery) Parse(request http.Request) (PagnaitedFeedQuery, error) {
	queryString := request.URL.Query()

	limit := queryString.Get("limit")
	if limit != "" {
		limit, err := strconv.Atoi(limit)
		feedQuery.Limit = limit
		if err != nil {
			return feedQuery, err
		}
	}

	offset := queryString.Get("offset")
	if offset != "" {
		offset, err := strconv.Atoi(offset)
		feedQuery.Offset = offset
		if err != nil {
			return feedQuery, err
		}
	}

	sort := queryString.Get("sort")
	if sort != "" {
		feedQuery.Sort = sort
	}

	return feedQuery, nil
}
