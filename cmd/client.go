package cmd

import (
	"fmt"

	"github.com/pinmonl/pinmonl/cmd/tui"
	"github.com/pinmonl/pinmonl/pmapi"
	"github.com/urfave/cli"
)

// NewClient creates Client cmd.
func NewClient(host string) Cmd {
	return Client{
		host: host,
	}
}

// Client manages and start TUI.
type Client struct {
	host string
}

// Command returns cli.Command of Client.
func (cl Client) Command() cli.Command {
	return cli.Command{
		Name:  "client",
		Usage: "text-based client",
		Action: func(c *cli.Context) error {
			if err := cl.ping(); err != nil {
				fmt.Println("Failed to connect. Please check the host URL.")
				return nil
			}

			app := tui.NewApp(cl.host, c.Bool("debug"))

			if err := app.Run(); err != nil {
				return err
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "debug, d",
				Usage: "print debug message",
			},
		},
	}
}

func (cl Client) ping() error {
	pmc := pmapi.NewClient(cl.host, nil)
	return pmc.Ping()
}
