package main

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("please specify direction (up/down).")
			return
		}

		db := newDB(cfg)
		defer db.Close()

		action := args[0]
		switch action {
		case "up":
			err := db.Migrate.Up()
			if err != nil && err != migrate.ErrNoChange {
				catchErr(err)
			}
		case "down":
			err := db.Migrate.Down()
			if err != nil && err != migrate.ErrNoChange {
				catchErr(err)
			}
		}
	},
}
