package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/exchange"
	"github.com/pinmonl/pinmonl/handler/web"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/generate"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/runner"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type application struct {
	cfg      *config
	db       *database.DB
	configs  *store.Configs
	exchange *exchange.Manager
	handler  http.Handler
	queue    *queue.Manager
	stores   *store.Stores
	runner   runner.Runner
	hub      pubsub.Pubsuber
}

func (a *application) bootstrapDefaultUser(ctx context.Context) error {
	// Create default user if not found.
	userId := a.configs.GetUserDefaultUserID()
	if userId == "" {
		user := model.User{Hash: generate.UserHash()}
		if err := a.stores.Users.Create(ctx, &user); err != nil {
			return err
		}
		userId = user.ID
		a.configs.SetUserDefaultUserID(userId)
		a.configs.Save()

		a.handler = newHandler(a.cfg, a.db, a.stores, a.queue, a.exchange, a.hub, a.configs)
	}

	return nil
}

func withApp(fn func(*cobra.Command, []string, *application)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		cfg, err := unmarshalConfig()
		catchErr(err)

		setupLogger(cfg)

		db := newDB(cfg)
		configs := newConfigStore(cfg)
		stores := store.NewStores(db)
		exm := newExchange(cfg, configs)
		hub := newPubsubHub(cfg, stores)
		qm := newQueue(cfg, db, stores, exm, hub)
		runner := newRunner(cfg, stores, qm, exm)
		handler := newHandler(cfg, db, stores, qm, exm, hub, configs)

		app := &application{
			cfg:      cfg,
			db:       db,
			configs:  configs,
			exchange: exm,
			handler:  handler,
			queue:    qm,
			stores:   stores,
			runner:   runner,
			hub:      hub,
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

func newDB(cfg *config) *database.DB {
	db, err := database.NewDB(
		cfg.DB.Driver,
		cfg.DB.DSN,
	)
	catchErr(err)
	return db
}

func newQueue(cfg *config, db *database.DB, stores *store.Stores, exm *exchange.Manager, hub pubsub.Pubsuber) *queue.Manager {
	qm, err := queue.NewManager(
		db,
		cfg.Queue.Job,
		cfg.Queue.Worker,
	)
	catchErr(err)
	qm = qm.Stores(stores).ExchangeManager(exm).Pubsuber(hub)
	return qm
}

func newExchange(cfg *config, configs *store.Configs) *exchange.Manager {
	exm, err := exchange.NewManager(
		configs,
		cfg.Exchange.Address,
	)
	catchErr(err)
	return exm
}

func newRunner(cfg *config, stores *store.Stores, qm *queue.Manager, exm *exchange.Manager) runner.Runner {
	r := &runner.ClientRunner{
		Queue:    qm,
		Exchange: exm,
		Stores:   stores,

		ExchangeEnabled: cfg.Exchange.Enabled,
	}
	return r
}

func newPubsubHub(cfg *config, stores *store.Stores) pubsub.Pubsuber {
	return pubsub.NewHub(
		[]byte(cfg.JWT.Secret),
		cfg.JWT.Expire,
		cfg.JWT.Issuer,
		stores.Users,
	)
}

func newHandler(cfg *config, db *database.DB, stores *store.Stores, qm *queue.Manager, exm *exchange.Manager, hub pubsub.Pubsuber, configs *store.Configs) http.Handler {
	web := &web.Server{
		Txer:        db,
		TokenSecret: []byte(cfg.JWT.Secret),
		TokenExpire: cfg.JWT.Expire,
		TokenIssuer: cfg.JWT.Issuer,
		Queue:       qm,
		Exchange:    exm,
		Pubsub:      hub,

		ExchangeEnabled: cfg.Exchange.Enabled,
		DevServer:       cfg.Web.DevServer,

		Images:    stores.Images,
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

	if cfg.DefaultUser {
		web.DefaultUserID = configs.GetUserDefaultUserID()
	}

	r := chi.NewRouter()
	r.Handle("/ws", hub.ServeWs())
	r.Mount("/", web.Handler())
	return r
}

func newConfigStore(_ *config) *store.Configs {
	configs := store.NewConfigs()
	return configs
}
