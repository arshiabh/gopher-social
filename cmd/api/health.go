package main

import "net/http"

func (app *application) HandleGetHealth(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"message": "hello"})
}
