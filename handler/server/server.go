package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
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
	r.Use(s.authenticate())
	r.Mount("/api", s.APIRouter())

	return r
}

func (s *Server) APIRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/info", s.infoHandler)
	r.Post("/login", s.loginHandler)
	r.Post("/signup", s.signupHandler)
	r.Post("/machine", s.machineSignupHandler)
	r.With(s.authorize()).
		Post("/alive", s.aliveHandler)

	r.Route("/share", func(r chi.Router) {
		r.Use(s.authorizeUserOnly())
		r.Get("/", nil)
		r.Route("/{slug}", func(r chi.Router) {
			r.Post("/", s.prepareShareHandler)
			r.With(s.bindShareBySlug()).
				Delete("/", s.deleteShareHandler)
			r.Route("/", func(r chi.Router) {
				r.Use(
					s.bindShareBySlug(),
					s.shareStatusMustBe(model.Preparing),
				)
				r.Post("/publish", s.publishShareHandler)
				r.Post("/tag", s.createShareTagHandler)
				r.Post("/pinl", s.createSharePinlHandler)
			})
		})
	})

	r.With(s.pagination()).
		Get("/pkgs/{proto:[a-z]+}://*", s.listPkgsHandler)
	r.Route("/pkg", func(r chi.Router) {
		suffix := "{provider:[a-z-]+}://*"

		r.With(
			s.checkPkgURI(),
			s.bindPkgByURI(),
			s.pagination(),
		).Get("/latest/"+suffix, s.listLatestStatsHandler)

		r.With(
			s.checkPkgURI(),
			s.bindPkgByURI(),
			s.pagination(),
		).Get("/stats/"+suffix, s.listStatsHandler)

		r.With(
			s.checkPkgURI(),
		).Get("/"+suffix, s.findPkgHandler)
	})

	r.Route("/sharing", func(r chi.Router) {
		r.Route("/{user}/{share}", func(r chi.Router) {
			r.Use(
				s.bindUserSharing(),
				s.shareStatusMustBe(model.Active),
			)
			r.Get("/", s.getSharingHandler)
			r.With(s.pagination()).
				Get("/pinl", s.listSharingPinlsHandler)
			r.With(s.pagination()).
				Get("/tag", s.listSharingTagsHandler)
		})
	})

	r.Route("/pinl", func(r chi.Router) {
		r.Use(s.authorize())
		r.With(s.pagination()).
			Get("/", s.listPinlsHandler)
		r.Post("/", s.createPinlHandler)
		r.Delete("/", s.clearPinlsHandler)
		r.Route("/{pinl}", func(r chi.Router) {
			r.Use(s.bindPinl())
			r.Delete("/", s.deletePinlHandler)
		})
	})

	return r
}

func (s *Server) pagination() func(http.Handler) http.Handler {
	return request.Pagination("page", "page_size", 10)
}
