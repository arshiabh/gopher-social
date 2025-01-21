package main

import (
	"net/http"

	"github.com/arshiabh/gopher-social/internal/store"
)

func (app *application) HandleGetFeed(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Order:  "desc",
	}
	if err := validate.Struct(fq); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := fq.Parse(r); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	feed, err := app.store.Posts.GetUserFeed(r.Context(), int64(34), fq)
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, feed)
}
