package main

import (
	"log"
	"net/http"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.Test)
	return mux
}

func (app *application) run(mux *http.ServeMux) error {
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}
	log.Printf("server run at %s", app.config.addr)
	return srv.ListenAndServe()
}

func (app *application) Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
