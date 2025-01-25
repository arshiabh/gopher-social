package main

import (
	"log"
	"net/http"
	"time"

	"github.com/arshiabh/gopher-social/internal/mail"
	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
	mail   mail.Client
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/health", func(r chi.Router) {
			r.Use(app.BasicAuthMiddleware)
			r.Get("/", app.HandleGetHealth)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.HandleCreatePosts)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.HandleGetPost)
				r.Delete("/", app.HandleDeletePost)
				r.Patch("/", app.HandlePatchPost)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.HandlePostActivate)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.UserContextMiddleware)
				r.Get("/", app.HandleGetUser)
				r.Put("/follow", app.HandleFollowUser)
				r.Put("/unfollow", app.HandleUnFollowUser)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.HandleGetFeed)
			})
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.HandleRegisterUser)
		})
	})
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server run at %s", app.config.addr)
	return srv.ListenAndServe()
}
