package cmd

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
	"github.com/urfave/cli"
)

// NewServer creates Server cmd.
func NewServer(endpoint string, singleUser bool, h http.Handler, qm *queue.Manager, users store.UserStore, mp *database.MigrationPlan) Cmd {
	return Server{
		Endpoint:      endpoint,
		Handler:       h,
		QueueManager:  qm,
		SingleUser:    singleUser,
		Users:         users,
		MigrationPlan: mp,
	}
}

// Server defines the dependencies of server command.
type Server struct {
	Endpoint      string
	Handler       http.Handler
	QueueManager  *queue.Manager
	MigrationPlan *database.MigrationPlan

	SingleUser bool
	Users      store.UserStore
}

// Command returns cli.Command of server.
func (s Server) Command() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "run HTTP server",
		Action: func(c *cli.Context) error {
			if !s.MigrationPlan.HasMigrationTable() {
				if err := s.MigrationPlan.Install(); err != nil {
					return err
				}
				if err := s.MigrationPlan.Up(); err != nil {
					return err
				}
			}
			if s.SingleUser {
				if err := s.initSingleUser(); err != nil {
					return err
				}
			}

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

func (s Server) initSingleUser() error {
	ctx := context.Background()
	users, err := s.Users.List(ctx, &store.UserOpts{ListOpts: store.ListOpts{Limit: 1}})
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return nil
	}
	user := model.User{}
	err = s.Users.Create(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}
