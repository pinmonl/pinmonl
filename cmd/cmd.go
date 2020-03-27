package cmd

import "github.com/urfave/cli"

// Cmd defines the command of app.
type Cmd interface {
	Command() cli.Command
}

// Cmds is slice of Cmd.
type Cmds []Cmd

// Commands returns slice of cli.Command.
func (cs Cmds) Commands() []cli.Command {
	out := make([]cli.Command, len(cs))
	for i, c := range cs {
		out[i] = c.Command()
	}
	return out
}
