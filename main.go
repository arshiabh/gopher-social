package main

import "net/http"

type Server struct {
	Addr string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from index"))
}

func main() {
	s := &Server{":8080"}
	http.ListenAndServe(s.Addr, s)
}
