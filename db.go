package main

import (
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/migrations"
)

func initDatabase(cfg *config.Config) (*database.DB, error) {
	return database.Open(
		cfg.DB.Driver,
		cfg.DB.DSN,
	)
}

func initMigrationPlan(db *database.DB) *database.MigrationPlan {
	src := database.PackrMigrationSource{
		Box: migrations.PackrBox(),
		Dir: db.DriverName(),
	}
	return database.NewMigrationPlan(db.DB, src)
}
