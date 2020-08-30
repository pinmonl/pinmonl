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
	Short: "start pinmonl server",
	Run: withApp(func(cmd *cobra.Command, args []string, app *application) {
		catchErr(app.migrateUp())

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			app.queue.Start()
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			logrus.Printf("listen on %s", app.cfg.Address)
			http.ListenAndServe(app.cfg.Address, app.handler)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			app.runner.Start()
			wg.Done()
		}()

		wg.Wait()
	}),
}
