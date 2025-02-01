package main

import (
	"log"

	"github.com/arshiabh/gopher-social/internal/auth"
	"github.com/arshiabh/gopher-social/internal/db"
	"github.com/arshiabh/gopher-social/internal/mail"
	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/arshiabh/gopher-social/internal/store/cache"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg := Config()
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	var rdb *redis.Client
	if cfg.redis.enable {
		rdb = cache.NewRedisClient(cfg.redis.addr, cfg.redis.password, cfg.redis.db)
	}
	cache := cache.NewRedisStorage(rdb)
	store := store.NewPostgresStorage(db)
	mailer := mail.NewSendGrip(cfg.mail.apiKey, cfg.mail.fromEmail)
	auth := auth.NewAuthentication(cfg.auth.secret)
	app := &application{
		config: *cfg,
		store:  store,
		cache:  cache,
		mail:   mailer,
		auth:   auth,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
