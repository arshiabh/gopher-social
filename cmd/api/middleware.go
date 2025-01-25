package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			writeErrJSON(w, http.StatusBadRequest, "authorization failed")
			return
		}
		ls := strings.Split(header, " ")
		if ls[0] != "Basic" {
			writeErrJSON(w, http.StatusBadRequest, "invalid authrization")
			return
		}
		src, _ := base64.StdEncoding.DecodeString(ls[1])
		params := strings.Split(string(src), ":")
		name := app.config.auth.name
		password := app.config.auth.password
		fmt.Println(params[0])
		if params[0] != name || params[1] != password {
			writeErrJSON(w, http.StatusBadRequest, "invalid name or password")
			return
		}
		next.ServeHTTP(w, r)
	})
}
