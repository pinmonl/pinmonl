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
	Run: withApp(func(cmd *cobra.Command, args []string, app *application) {
		if len(args) == 0 {
			fmt.Println("please specify direction (up/down).")
			return
		}

		action := args[0]
		switch action {
		case "up":
			err := app.db.Migrate.Up()
			if err != nil && err != migrate.ErrNoChange {
				catchErr(err)
			}
		case "down":
			err := app.db.Migrate.Down()
			if err != nil && err != migrate.ErrNoChange {
				catchErr(err)
			}
		}
	}),
}
