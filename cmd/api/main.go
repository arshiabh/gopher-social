package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/arshiabh/gopher-social/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: os.Getenv("addr"),
	}

	store := store.NewPostgresStorage(&sql.DB{})
	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
