package main

import (
	"log"
	"net/http"
	"time"

	"github.com/arshiabh/gopher-social/internal/auth"
	"github.com/arshiabh/gopher-social/internal/mail"
	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/arshiabh/gopher-social/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type application struct {
	config config
	store  store.Storage
	cache  cache.Storage
	mail   mail.Client
	auth   auth.Authenticator
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/health", func(r chi.Router) {
			r.Use(app.BasicAuthMiddleware)
			r.Get("/", app.HandleGetHealth)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.JWTAuthMiddleware)
			r.Post("/", app.HandleCreatePosts)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContext)
				r.Get("/", app.HandleGetPost)
				r.Patch("/", app.checkPostOwnership("moderator", app.HandlePatchPost))
				r.Delete("/", app.checkPostOwnership("admin", app.HandleDeletePost))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.HandlePostActivate)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.JWTAuthMiddleware)
				r.Get("/", app.HandleGetUser)
				r.Put("/follow", app.HandleFollowUser)
				r.Put("/unfollow", app.HandleUnFollowUser)
			})
			r.Group(func(r chi.Router) {
				r.Use(app.JWTAuthMiddleware)
				r.Get("/feed", app.HandleGetFeed)
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.HandleRegisterUser)
			r.Post("/token", app.HandlePostToken)
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
