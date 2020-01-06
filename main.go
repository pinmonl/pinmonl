package main

import (
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/logx"
)

func main() {
	cfg := config.Read()
	initLogger(cfg)

	db, err := initDatabase(cfg)
	if err != nil {
		logx.Fatal(err)
	}
	mp := initMigrationPlan(db)

	ss := initStores(db)
	sess := initSessionStore(cfg)
	ml := initMonl(cfg, ss)
	qm := initQueueManager(cfg, ss, ml)
	h := initHTTPHandler(cfg, ss, qm, sess)

	app := initCli(cfg, db, mp, h, ml, qm)
	err = app.Run(os.Args)
	if err != nil {
		logx.Fatal(err)
	}
}
