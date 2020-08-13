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
			catchErr(app.migrateUp())
		case "down":
			catchErr(app.migrateDown())
		}
	}),
}

func (a *application) migrateUp() error {
	err := a.db.Migrate.Up()
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (a *application) migrateDown() error {
	err := a.db.Migrate.Down()
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}
