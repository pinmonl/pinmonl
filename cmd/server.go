package cmd

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/pinmonl/pinmonl/queue"
	"github.com/urfave/cli"
)

// Server defines the dependencies of server command.
type Server struct {
	Endpoint     string
	Handler      http.Handler
	QueueManager *queue.Manager
}

// Command returns cli.Command of server.
func (s Server) Command() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "run HTTP server",
		Action: func(c *cli.Context) error {
			wg := &sync.WaitGroup{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				err := http.ListenAndServe(s.Endpoint, s.Handler)
				fmt.Printf("HTTP server error: %s\n", err.Error())
			}()

			wg.Add(1)
			ctx := context.Background()
			go func() {
				defer wg.Done()
				s.QueueManager.Start(ctx)
			}()

			fmt.Printf("HTTP server is running at %s\n", s.Endpoint)
			wg.Wait()
			return nil
		},
	}
}
