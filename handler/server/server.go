package server

import (
	"net/http"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Txer        database.Txer
	TokenSecret []byte
	TokenExpire time.Duration
	TokenIssuer string
	Queue       *queue.Manager
	Version     *semver.Version

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
	r.Use(s.logMachine())
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

	// r.Route("/share", func(r chi.Router) {
	// 	r.Use(s.authorizeUserOnly())
	// 	r.Get("/", nil)
	// 	r.Route("/{slug}", func(r chi.Router) {
	// 		r.Post("/", s.sharePrepareHandler)
	// 		r.With(s.bindShareBySlug()).
	// 			Delete("/", s.shareDeleteHandler)
	// 		r.Route("/", func(r chi.Router) {
	// 			r.Use(
	// 				s.bindShareBySlug(),
	// 				s.shareStatusMustBe(model.Preparing),
	// 			)
	// 			r.Post("/publish", s.sharePublishHandler)
	// 			r.Post("/tag/must", s.sharetagCreateHandler(model.SharetagMust))
	// 			r.Post("/tag/any", s.sharetagCreateHandler(model.SharetagAny))
	// 			r.Post("/pinl", s.sharepinCreateHandler)
	// 		})
	// 	})
	// })

	r.Route("/pkg", func(r chi.Router) {
		r.With(
			s.pagination(),
		).Get("/", s.pkgListHandler)
	})

	r.Route("/stat", func(r chi.Router) {
		r.With(
			s.pagination(),
		).Get("/", s.statListHandler)
	})

	// r.Route("/sharing", func(r chi.Router) {
	// 	r.Route("/{user}/{share}", func(r chi.Router) {
	// 		r.Use(
	// 			s.bindUser(),
	// 			s.bindUserSharing(),
	// 			s.shareStatusMustBe(model.Active),
	// 		)
	// 		r.Get("/", s.sharingHandler)
	// 		r.With(s.pagination()).
	// 			Get("/pinl", s.sharingPinlListHandler)
	// 		r.With(s.pagination()).
	// 			Get("/tag", s.sharingTagListHandler)
	// 	})
	// })

	// r.Route("/pinl", func(r chi.Router) {
	// 	r.Use(s.authorize())
	// 	r.With(s.pagination()).
	// 		Get("/", s.pinlListHandler)
	// 	r.Post("/", s.pinlCreateHandler)
	// 	r.Delete("/", s.pinlClearHandler)
	// 	r.Route("/{pinl}", func(r chi.Router) {
	// 		r.Use(s.bindPinl())
	// 		r.Delete("/", s.pinlDeleteHandler)
	// 	})
	// })

	return r
}

func (s *Server) pagination() func(http.Handler) http.Handler {
	return request.Pagination("page", "page_size", 10)
}

func (s *Server) logMachine() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logrus.Debugf("server: connected from %s", r.RemoteAddr)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
