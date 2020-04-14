package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
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
	defer db.Close()
	mp := initMigrationPlan(db)

	ss := initStores(db)
	sess := initSessionStore(cfg)
	ml := initMonler(cfg, ss)
	ws := initWebSocketServer(sess)
	qm := initQueueManager(cfg)
	sched := initQueueScheduler(qm, ws, ml, ss)
	h := initHTTPHandler(cfg, ss, qm, sess, ws)

	app := initCmd(cfg, db, mp, h, ml, qm, ss, sched)
	err = app.Run(os.Args)
	if err != nil {
		logx.Fatal(err)
	}
}
