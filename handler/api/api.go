package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/api/image"
	"github.com/pinmonl/pinmonl/handler/api/pinl"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/handler/api/share"
	"github.com/pinmonl/pinmonl/handler/api/tag"
	"github.com/pinmonl/pinmonl/handler/api/user"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/session"
	"github.com/pinmonl/pinmonl/store"
)

// ServerOpts defines the options of server initiation.
type ServerOpts struct {
	QueueManager  *queue.Manager
	CookieSession *session.CookieStore

	Store     store.Store
	Users     store.UserStore
	Pinls     store.PinlStore
	Tags      store.TagStore
	Taggables store.TaggableStore
	Shares    store.ShareStore
	Sharetags store.ShareTagStore
	Images    store.ImageStore
}

// Server defines the api server.
type Server struct {
	qm     *queue.Manager
	cookie *session.CookieStore

	store     store.Store
	users     store.UserStore
	pinls     store.PinlStore
	tags      store.TagStore
	taggables store.TaggableStore
	shares    store.ShareStore
	sharetags store.ShareTagStore
	images    store.ImageStore
}

// NewServer creates api server.
func NewServer(opts ServerOpts) *Server {
	return &Server{
		qm:     opts.QueueManager,
		cookie: opts.CookieSession,

		store:     opts.Store,
		users:     opts.Users,
		pinls:     opts.Pinls,
		tags:      opts.Tags,
		taggables: opts.Taggables,
		shares:    opts.Shares,
		sharetags: opts.Sharetags,
		images:    opts.Images,
	}
}

// Handler returns the handlers of api.
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(enableTransaction(s.store))
	r.Use(user.Authenticate(s.cookie, s.users))

	r.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(user.Guest())
			r.Post("/register", user.HandleRegister(s.cookie, s.users))
			r.Post("/login", user.HandleLogin(s.cookie, s.users))
		})
		r.Group(func(r chi.Router) {
			r.Use(user.Authorize())
			r.Post("/logout", user.HandleLogout(s.cookie))
		})
	})

	r.Route("/pinl", func(r chi.Router) {
		r.Use(user.Authorize())
		r.Get("/", pinl.HandleList(s.pinls, s.tags))
		r.Post("/", pinl.HandleCreate(s.pinls, s.tags, s.taggables, s.qm, s.images))
		r.Route("/{pinl}", func(r chi.Router) {
			r.Use(pinl.BindByID("pinl", s.pinls))
			r.Use(pinl.RequireOwner())
			r.Get("/", pinl.HandleFind(s.tags))
			r.Put("/", pinl.HandleUpdate(s.pinls, s.tags, s.taggables, s.qm, s.images))
			r.Delete("/", pinl.HandleDelete(s.pinls, s.taggables))
		})
	})

	r.Route("/tag", func(r chi.Router) {
		r.Use(user.Authorize())
		r.Get("/", tag.HandleList(s.tags))
		r.Post("/", tag.HandleCreate(s.tags))
		r.Route("/{tag}", func(r chi.Router) {
			r.Use(tag.BindByID("tag", s.tags))
			r.Use(tag.RequireOwner())
			r.Put("/", tag.HandleUpdate(s.tags))
			r.Delete("/", tag.HandleDelete(s.tags))
		})
	})

	r.Route("/share", func(r chi.Router) {
		r.Use(user.Authorize())
		r.Get("/", share.HandleList(s.shares))
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

	return r
}

func enableTransaction(s store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := response.NewTxWriter(w)

			ctx, err := s.BeginTx(r.Context())
			if err != nil {
				logx.Fatalf("api: fails to start transaction, err: %s", err)
			}
			next.ServeHTTP(ww, r.WithContext(ctx))

			if ww.Fails() {
				err = s.Rollback(ctx)
				if err != nil {
					logx.Fatalf("api: fails to rollback transaction, err: %s", err)
				}
			} else {
				err = s.Commit(ctx)
				if err != nil {
					logx.Fatalf("api: fails to commit transaction, err: %s", err)
				}
			}
		}
		return http.HandlerFunc(fn)
	}
}
