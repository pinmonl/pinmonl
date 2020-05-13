package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/api/image"
	"github.com/pinmonl/pinmonl/handler/api/pinl"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/session"
	"github.com/pinmonl/pinmonl/handler/api/share"
	"github.com/pinmonl/pinmonl/handler/api/sharing"
	"github.com/pinmonl/pinmonl/handler/api/tag"
	"github.com/pinmonl/pinmonl/handler/api/user"
	"github.com/pinmonl/pinmonl/handler/middleware"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/store"
)

var (
	// DefaultPageSize defines the default size of pagination.
	DefaultPageSize int64 = 50
)

// Handler returns the handlers of api.
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.EnableTransaction(s.store))
	if s.singleUser {
		r.Use(s.authAsFirstUser(s.users))
	} else {
		r.Use(user.Authenticate(s.cookie, s.users))
	}

	r.Route("/user", func(r chi.Router) {
		r.Post("/", user.HandleCreate(s.users))
	})

	r.Route("/me", func(r chi.Router) {
		r.Use(user.Authorize())
		r.Get("/", user.HandleGetMe())
		r.Put("/", user.HandleUpdateMe(s.users))
	})

	r.Route("/session", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(user.Guest())
			r.Post("/", session.HandleCreate(s.cookie, s.users))
		})
		r.Group(func(r chi.Router) {
			r.Use(user.Authorize())
			r.Delete("/", session.HandleDelete(s.cookie))
		})
	})

	r.Route("/pinl", func(r chi.Router) {
		r.Use(user.Authorize())
		r.With(pagination).
			Get("/", pinl.HandleList(s.pinls, s.taggables, s.monpkgs, s.stats))
		r.Get("/page-info", pinl.HandlePageInfo(s.pinls))
		r.Post("/", pinl.HandleCreate(s.pinls, s.tags, s.taggables, s.dp, s.images, s.pkgs, s.stats, s.pubsub))
		r.Route("/{pinl}", func(r chi.Router) {
			r.Use(pinl.BindByID("pinl", s.pinls))
			r.Use(pinl.RequireOwner())
			r.Get("/", pinl.HandleFind(s.taggables, s.monpkgs, s.pkgs, s.stats))
			r.Put("/", pinl.HandleUpdate(s.pinls, s.tags, s.taggables, s.dp, s.images, s.pkgs, s.monpkgs, s.stats, s.pubsub))
			r.Delete("/", pinl.HandleDelete(s.pinls, s.taggables, s.pubsub))
		})
	})

	r.Route("/tag", func(r chi.Router) {
		r.Use(user.Authorize())
		r.With(pagination).
			Get("/", tag.HandleList(s.tags))
		r.Get("/page-info", tag.HandlePageInfo(s.tags))
		r.Post("/", tag.HandleCreate(s.tags))
		r.Route("/{tag}", func(r chi.Router) {
			r.With(
				tag.BindByID("tag", s.tags),
				tag.RequireOwner(),
			).Get("/", tag.HandleFind())
			r.Route("/", func(r chi.Router) {
				r.Use(tag.BindByID("tag", s.tags))
				r.Use(tag.RequireOwner())
				r.Put("/", tag.HandleUpdate(s.tags))
				r.Delete("/", tag.HandleDelete(s.tags))
			})
		})
	})

	r.Route("/share", func(r chi.Router) {
		r.Use(user.Authorize())
		r.With(pagination).
			Get("/", share.HandleList(s.shares, s.sharetags))
		r.Get("/page-info", share.HandlePageInfo(s.shares))
		r.Post("/", share.HandleCreate(s.shares, s.sharetags, s.tags))
		r.Route("/{share}", func(r chi.Router) {
			r.Use(share.BindByID("share", s.shares))
			r.Use(share.RequireOwner())
			r.Get("/", share.HandleFind(s.shares, s.sharetags))
			r.Put("/", share.HandleUpdate(s.shares, s.sharetags, s.tags))
			r.Delete("/", share.HandleDelete(s.shares, s.sharetags))
		})
	})

	r.Route("/image", func(r chi.Router) {
		r.Route("/{image}", func(r chi.Router) {
			r.Use(image.BindByID("image", s.images))
			r.Get("/", image.HandleFind(s.images))
		})
	})

	r.Route("/sharing", func(r chi.Router) {
		r.Route("/{user}/{share}", func(r chi.Router) {
			r.Use(sharing.BindUser("user", s.users),
				sharing.BindUserShare("share", s.shares))
			r.Get("/", sharing.HandleFind(s.shares, s.sharetags))
			r.Get("/tag", sharing.HandleListTags(s.sharetags))
			r.Get("/pinl", sharing.HandleListPinls(s.sharetags, s.pinls, s.taggables))
			r.Get("/pinl/{pinl}", sharing.HandleFindPinl(s.sharetags, s.pinls, s.taggables, s.pkgs, s.stats))
		})
	})

	return r
}

func (s *Server) authAsFirstUser(users store.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			users, err := users.List(ctx, &store.UserOpts{ListOpts: store.ListOpts{Limit: 1}})
			if err != nil {
				logx.Debugln(err)
				return
			}
			if len(users) == 0 {
				logx.Debugln("no user is found in single user mode")
				return
			}
			user := users[0]
			r = r.WithContext(
				request.WithUser(ctx, user),
			)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func pagination(next http.Handler) http.Handler {
	return middleware.Pagination(DefaultPageSize)(next)
}
