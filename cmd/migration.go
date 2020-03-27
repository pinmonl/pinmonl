package cmd

import (
	"fmt"

	"github.com/pinmonl/pinmonl/database"
	"github.com/urfave/cli"
)

// NewMigration creates Migration cmd.
func NewMigration(mp *database.MigrationPlan) Cmd {
	return Migration{
		MigrationPlan: mp,
	}
}

// Migration defines the dependencies of migration command.
type Migration struct {
	MigrationPlan *database.MigrationPlan
}

// Command returns cli.Command of migration.
func (m Migration) Command() cli.Command {
	mp := m.MigrationPlan
	return cli.Command{
		Name:  "migration",
		Usage: "manage database version",
		Subcommands: []cli.Command{
			{
				Name:  "install",
				Usage: "install migration",
				Action: func(c *cli.Context) error {
					err := mp.Install()
					if err != nil {
						return err
					}
					fmt.Println("migration installed.")
					return nil
				},
			},
			{
				Name:  "up",
				Usage: "run migration",
				Action: func(c *cli.Context) error {
					limit := c.Int("limit")
					err := mp.UpTo(limit)
					if err != nil {
						return err
					}
					fmt.Println("migration is up.")
					return nil
				},
				Flags: []cli.Flag{
					cli.IntFlag{Name: "limit", Value: -1},
				},
			},
			{
				Name:  "down",
				Usage: "rollback migration",
				Action: func(c *cli.Context) error {
					limit := c.Int("limit")
					err := mp.DownTo(limit)
					if err != nil {
						return err
					}
					fmt.Println("migration is down.")
					return nil
				},
				Flags: []cli.Flag{
					cli.IntFlag{Name: "limit", Value: -1},
				},
			},
			{
				Name:  "reset",
				Usage: "reset migration",
				Action: func(c *cli.Context) error {
					err := mp.Down()
					if err != nil {
						return err
					}
					err = mp.Up()
					if err != nil {
						return err
					}
					fmt.Println("migration is reset.")
					return nil
				},
			},
			{
				Name:  "version",
				Usage: "show migration version",
				Action: func(c *cli.Context) error {
					rs := mp.Records()
					if len(rs) == 0 {
						fmt.Println("no migration")
						return nil
					}
					for i := len(rs) - 1; i >= 0; i-- {
						fmt.Printf("%q migrated.\n", rs[i].Name)
					}
					return nil
				},
			},
		},
	}
}
