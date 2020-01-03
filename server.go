package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/handler/api"
	"github.com/pinmonl/pinmonl/queue"
)

func initHTTPHandler(cfg *config.Config, ss stores, qm *queue.Manager) http.Handler {
	api := newAPIServer(ss, qm)

	r := chi.NewRouter()
	r.Mount("/api", api.Handler())
	return r
}

func newAPIServer(ss stores, qm *queue.Manager) *api.Server {
	return api.NewServer(api.ServerOpts{
		QueueManager: qm,
		Store:        ss.store,
		Users:        ss.users,
		Pinls:        ss.pinls,
		Tags:         ss.tags,
		Taggables:    ss.taggables,
		Shares:       ss.shares,
		Sharetags:    ss.sharetags,
	})
}
