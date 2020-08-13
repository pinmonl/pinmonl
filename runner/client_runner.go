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

	wg := sync.WaitGroup{}

	if c.ExchangeEnabled {
		wg.Add(1)
		go func() {
			c.keepExchangeAlive(ctx)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func (c *ClientRunner) bootstrap(ctx context.Context) error {
	if !c.ExchangeEnabled {
		return nil
	}

	logrus.Debugln("runner: bootstrap")

	if err := c.resumeMonls(ctx); err != nil {
		return err
	}
	if err := c.bootstrapExchangeClients(ctx); err != nil {
		return err
	}
	return nil
}

func (c *ClientRunner) resumeMonls(ctx context.Context) error {
	mList, err := c.Stores.Monls.List(ctx, &store.MonlOpts{
		FetchedBefore: time.Now().Add(-1 * 8 * time.Hour),
	})
	if err != nil {
		return err
	}

	logrus.Debugf("runner: resume monl n=%d", len(mList))

	for i := range mList {
		c.Queue.Add(job.NewFetchMonl(mList[i].ID))
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
