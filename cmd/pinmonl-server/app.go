package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/cmd/pinmonl-server/version"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/handler/server"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider/git"
	"github.com/pinmonl/pinmonl/monler/provider/github"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type application struct {
	cfg       *config
	db        *database.DB
	handler   http.Handler
	queue     *queue.Manager
	stores    *store.Stores
	scheduler *queue.Scheduler
}

func withApp(fn func(*cobra.Command, []string, *application)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg, err := unmarshalConfig()
		catchErr(err)

		setupLogger(cfg)
		setupMonler(cfg)

		db := newDB(cfg)
		stores := store.NewStores(db)
		qm := newQueue(cfg, db, stores)
		sch := newScheduler(cfg, stores, qm)
		handler := newHandler(cfg, db, stores, qm)

		app := &application{
			cfg:       cfg,
			db:        db,
			handler:   handler,
			queue:     qm,
			stores:    stores,
			scheduler: sch,
		}

		defer func() {
			db.Close()
		}()

		fn(cmd, args, app)
	}
}

func setupLogger(cfg *config) {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	switch cfg.Verbose {
	case 3:
		logrus.SetLevel(logrus.TraceLevel)
	case 2:
		logrus.SetLevel(logrus.DebugLevel)
	case 1:
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func setupMonler(cfg *config) {
	if gitPvd, err := git.NewProvider(); err == nil {
		monler.Register(gitPvd.ProviderName(), gitPvd)
	}
	if githubPvd, err := github.NewProvider(); err == nil {
		monler.Register(githubPvd.ProviderName(), githubPvd)
		github.AddToken(cfg.Github.Tokens)
	}
}

func newDB(cfg *config) *database.DB {
	db, err := database.NewDB(
		cfg.DB.Driver,
		cfg.DB.DSN,
	)
	catchErr(err)
	return db
}

func newQueue(cfg *config, db *database.DB, stores *store.Stores) *queue.Manager {
	qm, err := queue.NewManager(
		db,
		stores,
		cfg.Queue.Job,
		cfg.Queue.Worker,
	)
	catchErr(err)
	return qm
}

func newScheduler(cfg *config, stores *store.Stores, qm *queue.Manager) *queue.Scheduler {
	sch := &queue.Scheduler{
		Queue:  qm,
		Stores: stores,
	}
	return sch
}

func newHandler(cfg *config, db *database.DB, stores *store.Stores, qm *queue.Manager) http.Handler {
	server := &server.Server{
		Txer:        db,
		TokenSecret: []byte(cfg.JWT.Secret),
		TokenExpire: cfg.JWT.Expire,
		TokenIssuer: cfg.JWT.Issuer,
		Queue:       qm,
		Version:     version.Version,

		Monls:     stores.Monls,
		Monpkgs:   stores.Monpkgs,
		Pinls:     stores.Pinls,
		Pkgs:      stores.Pkgs,
		Sharepins: stores.Sharepins,
		Shares:    stores.Shares,
		Sharetags: stores.Sharetags,
		Stats:     stores.Stats,
		Taggables: stores.Taggables,
		Tags:      stores.Tags,
		Users:     stores.Users,
	}

	r := chi.NewRouter()
	r.Mount("/", server.Handler())
	return r
}
