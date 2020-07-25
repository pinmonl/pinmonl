package web

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/exchange"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
)

type Server struct {
	Txer        database.Txer
	TokenSecret []byte
	TokenExpire time.Duration
	TokenIssuer string
	Queue       *queue.Manager
	Exchange    *exchange.Manager
	Pubsub      pubsub.Pubsuber

	ExchangeEnabled bool

	Images    *store.Images
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

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(s.authenticate())
	r.Mount("/api", s.APIRouter())
	r.Mount("/", s.WebRouter())
	return r
}

func (s *Server) APIRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", s.signupHandler)
	r.Post("/login", s.loginHandler)
	r.With(s.authorize()).
		Post("/refresh", s.refreshHandler)

	r.Route("/pinl", func(r chi.Router) {
		r.Use(s.authorize())
		r.With(s.pagination()).
			Get("/", s.pinlListHandler)
		r.Post("/", s.pinlCreateHandler)
		r.Route("/{pinl}", func(r chi.Router) {
			r.Use(s.bindPinl())
			r.Get("/", s.pinlHandler)
			r.Put("/", s.pinlUpdateHandler)
			r.Delete("/", s.pinlDeleteHandler)
			r.Post("/image", s.pinlUploadImageHandler)
		})
	})

	r.Get("/card/*", s.fetchCardHandler)

	r.Route("/tag", func(r chi.Router) {
		r.Use(s.authorize())
		r.With(s.pagination()).
			Get("/", s.tagListHandler)
		r.Post("/", s.tagCreateHandler)
		r.Route("/{tag}", func(r chi.Router) {
			r.Use(s.bindTag())
			r.Get("/", s.tagHandler)
			r.Put("/", s.tagUpdateHandler)
			r.Delete("/", s.tagDeleteHandler)
		})
	})

	r.Route("/pkg", func(r chi.Router) {
		r.With(s.pagination()).
			Get("/", s.pkgListHandler)
	})

	r.Route("/stat", func(r chi.Router) {
		r.With(s.pagination()).
			Get("/", s.statListHandler)
	})

	r.Route("/share", func(r chi.Router) {
		r.Use(s.authorize())
		r.With(s.pagination()).
			Get("/", s.shareListHandler)
		r.Route("/{slug}", func(r chi.Router) {
			r.Post("/", s.shareCreateHandler)
			r.Route("/", func(r chi.Router) {
				r.Use(s.bindShare())
				r.Get("/", s.shareHandler)
				r.Delete("/", s.shareDeleteHandler)
				r.Post("/publish", s.sharePublishHandler)

				r.With(s.pagination()).
					Get("/tag", s.sharetagListHandler)
			})
		})
	})

	r.Route("/exchange", func(r chi.Router) {
		r.Post("/signup", nil)
		r.Post("/login", nil)
		r.Get("/status", nil)
	})

	return r
}

func (s *Server) WebRouter() chi.Router {
	devServer, _ := url.Parse("http://node:8080")
	devHandle := httputil.NewSingleHostReverseProxy(devServer)

	r := chi.NewRouter()
	r.With(s.bindImage()).Get("/image/{image}", s.imageHandler)
	r.Handle("/*", devHandle)

	return r
}

func (s *Server) pagination() func(http.Handler) http.Handler {
	return request.Pagination("page", "page_size", 10)
}
