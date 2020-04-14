package api

import (
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/session"
	"github.com/pinmonl/pinmonl/store"
)

// ServerOpts defines the options of server initiation.
type ServerOpts struct {
	SingleUser bool

	QueueManager  *queue.Manager
	CookieSession *session.CookieStore
	Pubsub        *pubsub.Server

	Store     store.Store
	Users     store.UserStore
	Pinls     store.PinlStore
	Tags      store.TagStore
	Taggables store.TaggableStore
	Shares    store.ShareStore
	Sharetags store.SharetagStore
	Images    store.ImageStore
	Monls     store.MonlStore
	Pkgs      store.PkgStore
	Stats     store.StatStore
}

// Server defines the api server.
type Server struct {
	singleUser bool

	qm     *queue.Manager
	cookie *session.CookieStore
	pubsub *pubsub.Server

	store     store.Store
	users     store.UserStore
	pinls     store.PinlStore
	tags      store.TagStore
	taggables store.TaggableStore
	shares    store.ShareStore
	sharetags store.SharetagStore
	images    store.ImageStore
	monls     store.MonlStore
	pkgs      store.PkgStore
	stats     store.StatStore
}

// NewServer creates api server.
func NewServer(opts ServerOpts) *Server {
	return &Server{
		singleUser: opts.SingleUser,

		qm:     opts.QueueManager,
		cookie: opts.CookieSession,
		pubsub: opts.Pubsub,

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
