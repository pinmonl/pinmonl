package api

import (
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
	Monls     store.MonlStore
	Pkgs      store.PkgStore
	Stats     store.StatStore
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
	monls     store.MonlStore
	pkgs      store.PkgStore
	stats     store.StatStore
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
		monls:     opts.Monls,
		pkgs:      opts.Pkgs,
		stats:     opts.Stats,
	}
}
