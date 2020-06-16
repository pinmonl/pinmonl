package main

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/handler/server"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

type stores struct {
	Monls     *store.Monls
	Monpkgs   *store.Monpkgs
	Pinls     *store.Pinls
	Pkgs      *store.Pkgs
	Sharepins *store.Sharepins
	Shares    *store.Shares
	Sharetags *store.Sharetags
	Stats     *store.Stats
	Taggables *store.Taggables
	Tags      *store.Tags
	Users     *store.Users
}

func newStores(db *database.DB) stores {
	s := store.NewStore(db)
	return stores{
		Monls:     store.NewMonls(s),
		Monpkgs:   store.NewMonpkgs(s),
		Pinls:     store.NewPinls(s),
		Pkgs:      store.NewPkgs(s),
		Sharepins: store.NewSharepins(s),
		Shares:    store.NewShares(s),
		Sharetags: store.NewSharetags(s),
		Stats:     store.NewStats(s),
		Taggables: store.NewTaggables(s),
		Tags:      store.NewTags(s),
		Users:     store.NewUsers(s),
	}
}

func newRouter(cfg *config, db *database.DB, s stores, queue *queue.Manager) http.Handler {
	r := chi.NewRouter()
	server := server.NewServer(&server.ServerOpts{
		Txer:        db,
		TokenSecret: []byte(cfg.JWT.Secret),
		TokenExpire: cfg.JWT.Expire,
		TokenIssuer: cfg.JWT.Issuer,
		Queue:       queue,

		Monls:     s.Monls,
		Monpkgs:   s.Monpkgs,
		Pinls:     s.Pinls,
		Pkgs:      s.Pkgs,
		Sharepins: s.Sharepins,
		Shares:    s.Shares,
		Sharetags: s.Sharetags,
		Stats:     s.Stats,
		Taggables: s.Taggables,
		Tags:      s.Tags,
		Users:     s.Users,
	})
	r.Mount("/", server.Handler())
	return r
}

func newDB(cfg *config) *database.DB {
	db, err := database.NewDB(&database.DBOpts{
		Driver: cfg.DB.Driver,
		DSN:    cfg.DB.DSN,
	})
	catchErr(err)
	return db
}

func newQueue(cfg *config, db *database.DB) *queue.Manager {
	qm, err := queue.NewManager(
		db,
		cfg.Queue.Job,
		cfg.Queue.Worker,
	)
	catchErr(err)
	return qm
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start pinmonl server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Run(migrateCmd, []string{"up"})

		db := newDB(cfg)
		defer db.Close()

		queue := newQueue(cfg, db)
		stores := newStores(db)
		router := newRouter(cfg, db, stores, queue)

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			queue.Start()
			wg.Done()
		}()

		go func() {
			logrus.Printf("listen on %s", cfg.Address)
			http.ListenAndServe(cfg.Address, router)
			wg.Done()
		}()

		wg.Wait()
	},
}
