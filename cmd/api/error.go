package main

import "net/http"

func (app *application) ErrNotFound(w http.ResponseWriter) {
	writeErrJSON(w, http.StatusNotFound, "data not found")
}
