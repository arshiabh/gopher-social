package main

import (
	"net/http"
	"strconv"

	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostsParams struct {
	Content string   `json:"content" validate:"required,max=100"`
	Title   string   `json:"title"   validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) HandleCreatePosts(w http.ResponseWriter, r *http.Request) {
	var PostsParams CreatePostsParams
	if err := readJSON(w, r, &PostsParams); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := validate.Struct(PostsParams); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &store.Post{
		Title:   PostsParams.Title,
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
	comments, err := app.store.Comments.GetByPostID(r.Context(), id)
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	post.Comments = comments
	writeJSON(w, http.StatusOK, post)
}

func (app *application) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	strID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	app.store.Posts.Delete(r.Context(), id)
	writeJSON(w, http.StatusAccepted, map[string]string{"message": "post deleted successfully"})
}

func (app *application) HandlePatchPost(w http.ResponseWriter, r *http.Request) {
	strID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	var PostsParams *CreatePostsParams
	if err := readJSON(w, r, &PostsParams); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := validate.Struct(PostsParams); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := app.store.Posts.GetByID(r.Context(), id)
	if err != nil {
		writeErrJSON(w, http.StatusNotFound, err.Error())
		return
	}
	if PostsParams.Title != "" {
		post.Title = PostsParams.Title
	}
	if PostsParams.Content != "" {
		post.Content = PostsParams.Content
	}
	if err := app.store.Posts.Patch(r.Context(), post); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, post)
}
