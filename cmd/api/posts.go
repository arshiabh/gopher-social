package main

import (
	"context"
	"errors"
	"fmt"
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
	user := getUserFromCtx(r)

	post := &store.Post{
		Title:   PostsParams.Title,
		Content: PostsParams.Content,
		Tags:    PostsParams.Tags,
		UserID:  user.ID,
	}
	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusAccepted, &post)
}

func (app *application) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	post, err := getPostFromCtx(r)
	if err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
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
	post.Comments = comments
	jsonResponse(w, http.StatusOK, post)
}

func (app *application) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	post, _ := getPostFromCtx(r)
	user := getUserFromCtx(r)
	if user.ID != post.UserID {
		writeErrJSON(w, http.StatusForbidden, "user is not auhtorized")
		return
	}
	if err := app.store.Posts.Delete(r.Context(), post.ID); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusAccepted, map[string]string{"message": "post deleted successfully"})
}

func (app *application) HandlePatchPost(w http.ResponseWriter, r *http.Request) {
	var PostsParams *CreatePostsParams
	if err := readJSON(w, r, &PostsParams); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	user := getUserFromCtx(r)
	post, _ := getPostFromCtx(r)
	if user.ID != post.UserID {
		writeErrJSON(w, http.StatusForbidden, "user is not auhtorized")
		return
	}
	if err := validate.Struct(PostsParams); err != nil {
		writeErrJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	post.Title = PostsParams.Title
	post.Content = PostsParams.Content

	if err := app.store.Posts.Patch(r.Context(), post); err != nil {
		writeErrJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonResponse(w, http.StatusAccepted, post)
}

type PostCtx string

func (app *application) postContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strID := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(strID, 10, 64)
		if err != nil {
			writeErrJSON(w, http.StatusBadRequest, "invalid type for id")
			return
		}
		post, err := app.store.Posts.GetByID(r.Context(), postID)
		if err != nil {
			writeErrJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		var postctx PostCtx = "post"
		ctx := r.Context()
		ctx = context.WithValue(ctx, postctx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) (*store.Post, error) {
	var postctx PostCtx = "post"
	post, ok := r.Context().Value(postctx).(*store.Post)
	if !ok {
		return nil, fmt.Errorf("failed to get post")
	}
	return post, nil
}
