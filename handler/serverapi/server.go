package serverapi

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/store"
)

type Server struct {
	users *store.Users
}

type ServerOpts struct {
	Users *store.Users
}

func NewServer(opts *ServerOpts) *Server {
	return &Server{
		users: opts.Users,
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Get("/version", Version)

	r.Post("/signup", Signup(s.users))

	return r
}
