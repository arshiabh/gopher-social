package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/arshiabh/gopher-social/internal/store"
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
		userID := claims["sub"].(float64)
		user, err := app.getUser(r.Context(), int64(userID))
		if err != nil {
			writeErrJSON(w, http.StatusUnauthorized, err.Error())
			return

		}
		//set user when jwt done to context
		var userStr userCtx = "user"
		ctx := context.WithValue(r.Context(), userStr, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)
		roleName, err := app.store.Users.GetUserRole(r.Context(), user.ID)
		if err != nil {
			writeErrJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}
		if roleName != role {
			writeErrJSON(w, http.StatusForbidden, "forbidden action")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	user, err := app.cache.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err = app.store.Users.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if err := app.cache.User.Set(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}
