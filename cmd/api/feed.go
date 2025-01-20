package main

import "net/http"

func (app *application) HandleGetFeed(w http.ResponseWriter, r *http.Request) {
	feed, err := app.store.Posts.GetUserFeed(r.Context(), int64(42))
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, feed)
}
