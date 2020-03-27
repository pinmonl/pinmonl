package cmd

import (
	"fmt"

	"github.com/pinmonl/pinmonl/pkg/generate"
	"github.com/urfave/cli"
)

// NewGenerate create Generate cmd.
func NewGenerate() Cmd {
	return Generate{}
}

// Generate defines the dependencies of generate command.
type Generate struct {
	//
}

// Command returns cli.Command of generate.
func (g Generate) Command() cli.Command {
	return cli.Command{
		Name:  "generate",
		Usage: "helper for generating key and certificate",
		Subcommands: []cli.Command{
			{
				Name: "key",
				Action: func(c *cli.Context) error {
					fmt.Println(generate.RandomString(c.Int("length")))
					return nil
				},
				Flags: []cli.Flag{
					cli.IntFlag{Name: "length, l, L", Value: 32},
				},
			},
		},
	}
}
