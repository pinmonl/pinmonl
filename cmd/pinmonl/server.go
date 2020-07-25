package main

import (
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start client web server",
	Run: withApp(func(cmd *cobra.Command, args []string, app *application) {
		migrateCmd.Run(cmd, []string{"up"})

		wg := &sync.WaitGroup{}
		wg.Add(4)

		go func() {
			logrus.Debugf("listen on %s", app.cfg.Address)
			http.ListenAndServe(app.cfg.Address, app.handler)
			wg.Done()
		}()

		go func() {
			app.queue.Start()
			wg.Done()
		}()

		go func() {
			app.hub.Start()
			wg.Done()
		}()

		go func() {
			app.runner.Start()
			wg.Done()
		}()

		wg.Wait()
	}),
}
