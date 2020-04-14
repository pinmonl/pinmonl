package main

import (
	"net/http"

	"github.com/pinmonl/pinmonl/cmd"
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/urfave/cli"
)

func initCmd(
	cfg *config.Config,
	db *database.DB,
	mp *database.MigrationPlan,
	h http.Handler,
	ml *monler.Repository,
	qm *queue.Manager,
	ss stores,
	sched *queue.Scheduler,
) *cli.App {
	cmds := cmd.Cmds{
		cmd.NewClient(cfg.Client.Host),
		cmd.NewGenerate(),
		cmd.NewMigration(mp),
		cmd.NewServer(cfg.HTTP.Endpoint, cfg.SingleUser, h, qm, ss.users, mp, sched),
	}

	return &cli.App{
		Name:     "pinmonl",
		Version:  "dev",
		Commands: cmds.Commands(),
	}
}
