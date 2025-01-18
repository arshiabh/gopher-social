package main

import (
	"fmt"
	"log"

	"github.com/arshiabh/gopher-social/internal/db"
	"github.com/arshiabh/gopher-social/internal/store"
)

func main() {
	DB, err := db.New("host=localhost port=5432 user=postgres password=09300617050 dbname=gopher-database sslmode=disable", 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()
	store := store.NewPostgresStorage(DB)
	db.Seed(&store)
	fmt.Println("seeding is done")
}
