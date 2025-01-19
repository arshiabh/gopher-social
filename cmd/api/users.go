package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userCtx string

func (app *application) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		log.Fatal(err)
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) HandleFollowUser(w http.ResponseWriter, r *http.Request) {
	follower := getUserFromCtx(r)
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.store.Followers.Follow(r.Context(), payload.UserID, follower); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusAccepted, map[string]string{"message": "successfully followed"})
}

func (app *application) HandleUnFollowUser(w http.ResponseWriter, r *http.Request) {
	follower := getUserFromCtx(r)
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.store.Followers.UnFollow(r.Context(), payload.UserID, follower); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusAccepted, map[string]string{"message": "successfully unfollowed"})
}

func (app *application) UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strID := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(strID, 10, 64)
		if err != nil {
			writeErrJSON(w, http.StatusBadRequest, "invalid type for id")
			return
		}
		user, err := app.store.Users.GetByUserID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.ErrNotFound(w)
				return
			default:
				writeErrJSON(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		ctx := r.Context()
		var userStr userCtx = "user"
		ctx = context.WithValue(ctx, userStr, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	var userctx userCtx = "user"
	user, _ := r.Context().Value(userctx).(*store.User)
	return user
}
