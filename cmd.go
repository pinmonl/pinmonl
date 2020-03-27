package main

import (
	"net/http"

	"github.com/pinmonl/pinmonl/cmd"
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/monl"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/urfave/cli"
)

func initCmd(
	cfg *config.Config,
	db *database.DB,
	mp *database.MigrationPlan,
	h http.Handler,
	ml *monl.Monl,
	qm *queue.Manager,
	ss stores,
) *cli.App {
	cmds := cmd.Cmds{
		cmd.NewClient(cfg.Client.Host),
		cmd.NewGenerate(),
		cmd.NewMigration(mp),
		cmd.NewServer(cfg.HTTP.Endpoint, cfg.SingleUser, h, qm, ss.users, mp),
	}

	return &cli.App{
		Name:     "pinmonl",
		Version:  "0.1.0",
		Commands: cmds.Commands(),
	}
}
