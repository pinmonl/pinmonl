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

	Dispatcher    *queue.Dispatcher
	CookieSession *session.CookieStore
	Pubsub        *pubsub.Server

	Store     store.Store
	Images    store.ImageStore
	Monls     store.MonlStore
	Monpkgs   store.MonpkgStore
	Pinls     store.PinlStore
	Pkgs      store.PkgStore
	Shares    store.ShareStore
	Sharetags store.SharetagStore
	Stats     store.StatStore
	Taggables store.TaggableStore
	Tags      store.TagStore
	Users     store.UserStore
}

// Server defines the api server.
type Server struct {
	singleUser bool

	dp     *queue.Dispatcher
	cookie *session.CookieStore
	pubsub *pubsub.Server

	store     store.Store
	images    store.ImageStore
	monls     store.MonlStore
	monpkgs   store.MonpkgStore
	pinls     store.PinlStore
	pkgs      store.PkgStore
	shares    store.ShareStore
	sharetags store.SharetagStore
	stats     store.StatStore
	taggables store.TaggableStore
	tags      store.TagStore
	users     store.UserStore
}

// NewServer creates api server.
func NewServer(opts ServerOpts) *Server {
	return &Server{
		singleUser: opts.SingleUser,

		dp:     opts.Dispatcher,
		cookie: opts.CookieSession,
		pubsub: opts.Pubsub,

		store:     opts.Store,
		images:    opts.Images,
		monls:     opts.Monls,
		monpkgs:   opts.Monpkgs,
		pinls:     opts.Pinls,
		pkgs:      opts.Pkgs,
		shares:    opts.Shares,
		sharetags: opts.Sharetags,
		stats:     opts.Stats,
		taggables: opts.Taggables,
		tags:      opts.Tags,
		users:     opts.Users,
	}
}
