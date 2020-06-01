package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	//
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	return r
}
