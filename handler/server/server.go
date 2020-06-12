package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
)

type Server struct {
	Txer        database.Txer
	TokenSecret []byte
	TokenExpire time.Duration
	TokenIssuer string
	Queue       *queue.Manager

	Monls     *store.Monls
	Monpkgs   *store.Monpkgs
	Pinls     *store.Pinls
	Pkgs      *store.Pkgs
	Sharepins *store.Sharepins
	Shares    *store.Shares
	Sharetags *store.Sharetags
	Stats     *store.Stats
	Taggables *store.Taggables
	Tags      *store.Tags
	Users     *store.Users
}

type ServerOpts struct {
	Txer        database.Txer
	TokenSecret []byte
	TokenExpire time.Duration
	TokenIssuer string
	Queue       *queue.Manager

	Monls     *store.Monls
	Monpkgs   *store.Monpkgs
	Pinls     *store.Pinls
	Pkgs      *store.Pkgs
	Sharepins *store.Sharepins
	Shares    *store.Shares
	Sharetags *store.Sharetags
	Stats     *store.Stats
	Taggables *store.Taggables
	Tags      *store.Tags
	Users     *store.Users
}

func NewServer(opts *ServerOpts) *Server {
	return &Server{
		Txer:        opts.Txer,
		TokenSecret: opts.TokenSecret,
		TokenExpire: opts.TokenExpire,
		TokenIssuer: opts.TokenIssuer,
		Queue:       opts.Queue,

		Monls:     opts.Monls,
		Monpkgs:   opts.Monpkgs,
		Pinls:     opts.Pinls,
		Pkgs:      opts.Pkgs,
		Sharepins: opts.Sharepins,
		Shares:    opts.Shares,
		Sharetags: opts.Sharetags,
		Stats:     opts.Stats,
		Taggables: opts.Taggables,
		Tags:      opts.Tags,
		Users:     opts.Users,
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Mount("/api", s.APIRouter())

	return r
}

func (s *Server) APIRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/info", s.infoHandler)
	r.Post("/login", s.loginHandler)
	r.Post("/signup", s.signupHandler)
	r.With(s.authenticate()).
		Post("/alive", s.aliveHandler)

	r.Route("/share", func(r chi.Router) {
		r.Use(s.authenticate())
		r.Get("/", nil)
		r.Route("/{slug}", func(r chi.Router) {
			r.Post("/prepare", s.prepareShareHandler)
			r.Route("/", func(r chi.Router) {
				r.Use(s.bindShareBySlug("slug"))
				r.Post("/publish", s.publishShareHandler)
				r.Post("/tag", s.createShareTagHandler)
				r.Post("/pinl", s.createSharePinlHandler)
				r.Delete("/", nil)
			})
		})
	})

	r.Route("/pkg", func(r chi.Router) {
		suffix := "{provider:[a-z-]+}/*"
		r.Get("/stats/"+suffix, nil)

		r.With(s.bindPkgURI()).
			Get("/"+suffix, s.pkgLatestHandler)
	})

	r.Route("/sharing", func(r chi.Router) {
		r.Route("/{user}/{share}", func(r chi.Router) {
			r.Get("/", nil)
			r.Get("/pinl", nil)
			r.Get("/tag", nil)
		})
	})

	return r
}

func (s *Server) pagination() func(http.Handler) http.Handler {
	return request.Pagination("page", "page_size")
}
