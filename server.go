package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/handler/api"
	"github.com/pinmonl/pinmonl/handler/web"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/queue"
)

func initHTTPHandler(cfg *config.Config, ss stores, qm *queue.Manager, dp *queue.Dispatcher, sess sessions, ws *pubsub.Server) http.Handler {
	api := newAPIServer(cfg, ss, dp, sess, ws)
	web := newWebServer(cfg, ss, sess)

	r := chi.NewRouter()
	r.Mount("/api", api.Handler())
	r.Mount("/ws", ws.Handler())
	r.Mount("/", web.Handler())
	return r
}

func newAPIServer(cfg *config.Config, ss stores, dp *queue.Dispatcher, sess sessions, pubsub *pubsub.Server) *api.Server {
	return api.NewServer(api.ServerOpts{
		SingleUser: cfg.SingleUser,

		Dispatcher:    dp,
		CookieSession: sess.cookie,
		Pubsub:        pubsub,

		Store:     ss.store,
		Images:    ss.images,
		Monls:     ss.monls,
		Monpkgs:   ss.monpkgs,
		Pinls:     ss.pinls,
		Pkgs:      ss.pkgs,
		Shares:    ss.shares,
		Sharetags: ss.sharetags,
		Stats:     ss.stats,
		Taggables: ss.taggables,
		Tags:      ss.tags,
		Users:     ss.users,
	})
}

func newWebServer(cfg *config.Config, ss stores, sess sessions) *web.Server {
	return web.NewServer(web.ServerOpts{
		DevServer: cfg.HTTP.DevServer,
	})
}

func initWebSocketServer(cfg *config.Config, sess sessions, ss stores) *pubsub.Server {
	ws, err := pubsub.NewServer(&pubsub.ServerOpts{
		SingleUser: cfg.SingleUser,
		Cookie:     sess.cookie,
		Users:      ss.users,
	})
	if err != nil {
		logx.Panic(err)
	}
	return ws
}
