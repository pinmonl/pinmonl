package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/handler/api"
	"github.com/pinmonl/pinmonl/handler/web"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/queue"
)

func initHTTPHandler(cfg *config.Config, ss stores, qm *queue.Manager, sess sessions, ws *pubsub.Server) http.Handler {
	api := newAPIServer(cfg, ss, qm, sess, ws)
	web := newWebServer(cfg, ss, sess)

	r := chi.NewRouter()
	r.Mount("/api", api.Handler())
	r.Mount("/ws", ws.Handler())
	r.Mount("/", web.Handler())
	return r
}

func newAPIServer(cfg *config.Config, ss stores, qm *queue.Manager, sess sessions, pubsub *pubsub.Server) *api.Server {
	return api.NewServer(api.ServerOpts{
		SingleUser: cfg.SingleUser,

		QueueManager:  qm,
		CookieSession: sess.cookie,
		Pubsub:        pubsub,

		Store:     ss.store,
		Users:     ss.users,
		Pinls:     ss.pinls,
		Tags:      ss.tags,
		Taggables: ss.taggables,
		Shares:    ss.shares,
		Sharetags: ss.sharetags,
		Images:    ss.images,
		Monls:     ss.monls,
		Pkgs:      ss.pkgs,
		Stats:     ss.stats,
	})
}

func newWebServer(cfg *config.Config, ss stores, sess sessions) *web.Server {
	return web.NewServer(web.ServerOpts{
		DevServer: cfg.HTTP.DevServer,
	})
}

func initWebSocketServer(sess sessions) *pubsub.Server {
	return pubsub.NewServer(&pubsub.ServerOpts{
		Cookie: sess.cookie,
	})
}
