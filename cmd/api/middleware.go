package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
		name := app.config.auth.baseconfig.name
		password := app.config.auth.baseconfig.password
		fmt.Println(params[0])
		if params[0] != name || params[1] != password {
			writeErrJSON(w, http.StatusBadRequest, "invalid name or password")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			writeErrJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		parts := strings.Split(token, " ")
		if len(parts) != 2 && parts[0] != "Bearer" {
			writeErrJSON(w, http.StatusUnauthorized, "invalid authorization")
			return
		}
		tokenstr := parts[1]
		jwtToken, err := app.auth.ValidateToken(tokenstr)
		if err != nil {
			writeErrJSON(w, http.StatusUnauthorized, "invalid authorization")
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		sub := claims["sub"]
		type userCtx interface{}
		var userctx userCtx = "userID"
		ctx := context.WithValue(r.Context(), userctx, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
