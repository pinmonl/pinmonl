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
) *cli.App {
	server := cmd.Server{
		Endpoint:     cfg.HTTP.Endpoint,
		Handler:      h,
		QueueManager: qm,
	}
	migration := cmd.Migration{
		MigrationPlan: mp,
	}
	generate := cmd.Generate{}

	return &cli.App{
		Name:    "pinmonl",
		Version: "0.1.0",
		Commands: []cli.Command{
			server.Command(),
			migration.Command(),
			generate.Command(),
		},
	}
}
