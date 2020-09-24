package runner

import (
	"context"
	"sync"
	"time"

	"github.com/pinmonl/pinmonl/exchange"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
)

type ClientRunner struct {
	Queue    *queue.Manager
	Exchange *exchange.Manager
	Stores   *store.Stores

	ExchangeEnabled bool
}

func (c *ClientRunner) Start() error {
	ctx := context.TODO()
	if err := c.bootstrap(ctx); err != nil {
		return err
	}

	if err := c.start(ctx); err != nil {
		return err
	}
	return nil
}

func (c *ClientRunner) bootstrap(ctx context.Context) error {
	if !c.ExchangeEnabled {
		return nil
	}

	logrus.Debugln("runner: bootstrap")

	if err := c.bootstrapExchangeClients(ctx); err != nil {
		return err
	}
	return nil
}

func (c *ClientRunner) bootstrapExchangeClients(ctx context.Context) error {
	if c.Exchange.HasMachine() {
		logrus.Debugln("runner: bootstrap exchange machine alive")
		c.Exchange.MachineAlive()
	} else {
		logrus.Debugln("runner: bootstrap exchange machine signup")
		c.Exchange.MachineSignup()
	}

	if c.Exchange.HasUser() {
		logrus.Debugln("runner: bootstrap exchange user alive")
		c.Exchange.Alive()
	}

	if err := c.uploadUniqueURLs(ctx); err != nil {
		return err
	}

	return nil
}

func (c *ClientRunner) start(ctx context.Context) error {
	wg := sync.WaitGroup{}

	if c.ExchangeEnabled {
		wg.Add(1)
		go func() {
			c.keepExchangeAlive(ctx)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			c.regularUpdateMonls(ctx)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			c.regularUpdatePkgs(ctx)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func (c *ClientRunner) keepExchangeAlive(ctx context.Context) error {
	ticker := time.NewTicker(24 * time.Hour)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			logrus.Debugln("runner: cron exchange alive starts")

			c.Exchange.MachineAlive()
			c.uploadUniqueURLs(ctx)

			if c.Exchange.HasUser() {
				c.Exchange.Alive()
			}
		}
	}
	return nil
}

func (c *ClientRunner) regularUpdateMonls(ctx context.Context) error {
	interval := 3 * time.Hour
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()
	c.updateMonls(ctx, time.Now().Add(-1*interval))
	for {
		select {
		case <-ticker.C:
			before := time.Now().Add(-1 * interval)
			c.updateMonls(ctx, before)
		}
	}
}

func (c *ClientRunner) updateMonls(ctx context.Context, before time.Time) error {
	logrus.Debugln("runner: start monls update")
	expired, err := c.Stores.Monls.List(ctx, &store.MonlOpts{
		FetchedBefore: before,
	})
	if err != nil {
		return err
	}

	for _, monl := range expired {
		c.Queue.Add(job.NewFetchMonl(monl.ID))
	}
	logrus.Debugf("runner: %d monls updated", len(expired))
	return nil
}

func (c *ClientRunner) regularUpdatePkgs(ctx context.Context) error {
	interval := 1 * time.Hour
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()
	c.updatePkgs(ctx, time.Now().Add(-1*interval))
	for {
		select {
		case <-ticker.C:
			before := time.Now().Add(-1 * interval)
			c.updatePkgs(ctx, before)
		}
	}
}

func (c *ClientRunner) updatePkgs(ctx context.Context, before time.Time) error {
	logrus.Debugln("runner: start pkgs update")
	expired, err := c.Stores.Pkgs.List(ctx, &store.PkgOpts{
		FetchedBefore: before,
	})
	if err != nil {
		return err
	}

	for _, pkg := range expired {
		c.Queue.Add(job.NewFetchPkg(pkg.ID))
	}
	logrus.Debugf("runner: %d monls updated", len(expired))
	return nil
}

func (c *ClientRunner) uploadUniqueURLs(ctx context.Context) error {
	// logrus.Debugln("runner: upload unique urls")
	// monls, err := c.Stores.Monls.List(ctx, nil)
	// if err != nil {
	// 	return err
	// }

	// client := c.Exchange.MachineClient()
	// if err := client.PinlClear(); err != nil {
	// 	return err
	// }
	// for _, monl := range monls {
	// 	pinl := &pinmonl.Pinl{URL: monl.URL}
	// 	_, err := client.PinlCreate(pinl)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// logrus.Debugf("runner: total of %d unique urls uploaded", len(monls))
	return nil
}

var _ Runner = &ClientRunner{}
