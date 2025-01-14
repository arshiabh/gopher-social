package main

import (
	"net/http"
	"strconv"

	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostsParams struct {
	Content string   `json:"content"`
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
}

func (app *application) HandleCreatePosts(w http.ResponseWriter, r *http.Request) {
	var PostsParams CreatePostsParams
	if err := readJSON(w, r, &PostsParams); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	post := &store.Post{
		Tite:    PostsParams.Title,
		Content: PostsParams.Content,
		Tags:    PostsParams.Tags,
		UserID:  1,
	}
	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, &post)
}

func (app *application) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, post)
}
