package main

import (
	"log"

	"github.com/arshiabh/gopher-social/internal/db"
	"github.com/arshiabh/gopher-social/internal/mail"
	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/joho/godotenv"
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
	store := store.NewPostgresStorage(db)
	mailer := mail.NewSendGrip(cfg.mail.apiKey, cfg.mail.fromEmail)
	app := &application{
		config: *cfg,
		store:  store,
		mail:   mailer,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
