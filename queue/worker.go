package queue

import (
	"context"
	"fmt"

	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monl"
	"github.com/pinmonl/pinmonl/store"
)

type worker struct {
	manager *Manager
	monl    *monl.Monl

	store store.Store
	pinls store.PinlStore
	monls store.MonlStore
	stats store.StatStore
}

// run is goroutine which listens to job queue.
func (w *worker) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case job := <-w.manager.jobQueue:
			w.process(ctx, job)
		}
	}
	return ctx.Err()
}

// process handles job by the job name.
func (w *worker) process(ctx context.Context, job *Job) error {
	err := w.manager.jobStarted(ctx, job)
	if err != nil {
		return err
	}

	ctx, err = w.store.BeginTx(ctx)
	if err != nil {
		logx.Fatal("queue worker: fails to start tx, err: %s", err)
		return err
	}
	defer func() {
		if err != nil {
			w.store.Rollback(ctx)
		} else {
			w.store.Commit(ctx)
		}
	}()

	switch job.Detail.Name {
	case model.JobPinlCreated:
		err = w.processPinlCreated(ctx, job)
	case model.JobPinlUpdated:
		err = w.processPinlCreated(ctx, job)
	case model.JobMonlCreated:
		err = w.processMonlCreated(ctx, job)
	case model.JobMonlRegularUpdate:
		err = w.processMonlRegularUpdate(ctx, job)
	default:
		return w.jobStopped(ctx, job, fmt.Errorf("queue worker: job name is not defined"))
	}

	if err != nil {
		logx.Errorf("queue worker: err: %s", err)
		return w.jobStopped(ctx, job, err)
	}
	return w.manager.jobCompleted(ctx, job)
}

// processPinlCreated handles the hook after pinl created.
func (w *worker) processPinlCreated(ctx context.Context, job *Job) error {
	t := model.Pinl{ID: job.Detail.TargetID}
	err := w.pinls.Find(ctx, &t)
	if err != nil {
		return err
	}

	tml, err := w.monls.List(ctx, &store.MonlOpts{URL: t.URL})
	if err != nil {
		return err
	}
	if len(tml) > 0 {
		return nil
	}

	vendors := w.monl.GuessURL(t.URL)
	for _, v := range vendors {
		logx.Debugf("queue worker: working on monl vendor - %s", v.Name())
		ml, err := w.monls.List(ctx, &store.MonlOpts{
			URL:    t.URL,
			Vendor: v.Name(),
		})
		if err != nil {
			return err
		}
		if len(ml) > 0 {
			continue
		}

		m := model.Monl{URL: t.URL, Vendor: v.Name()}
		err = w.monls.Create(ctx, &m)
		if err != nil {
			return err
		}

		err = w.manager.Enqueue(ctx, &model.Job{
			Name:     model.JobMonlCreated,
			TargetID: m.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// processMonlCreated handles the hook after monl created.
func (w *worker) processMonlCreated(ctx context.Context, job *Job) error {
	t := model.Monl{ID: job.Detail.TargetID}
	err := w.monls.Find(ctx, &t)
	if err != nil {
		return err
	}

	v, err := w.monl.Get(t.Vendor)
	if err != nil {
		return err
	}

	r, err := v.Load(ctx, t.URL)
	if err != nil {
		return err
	}
	defer r.Close()

	t.URL = r.RawURL()
	t.VendorURI = r.URI()
	err = w.monls.Update(ctx, &t)
	if err != nil {
		return err
	}

	err = w.parseMonlReport(ctx, t, r)
	if err != nil {
		return err
	}

	return nil
}

// processMonlRegularUpdate handles regular update of monl.
func (w *worker) processMonlRegularUpdate(ctx context.Context, job *Job) error {
	t := model.Monl{ID: job.Detail.TargetID}
	err := w.monls.Find(ctx, &t)
	if err != nil {
		return err
	}

	v, err := w.monl.Get(t.Vendor)
	if err != nil {
		return err
	}

	r, err := v.Load(ctx, t.URL)
	if err != nil {
		return err
	}
	defer r.Close()

	err = w.parseMonlReport(ctx, t, r)
	if err != nil {
		logx.Debugf("queue worker: fails to parse Monl report, err(%s)", err)
		return err
	}
	return nil
}

// parseMonlReport parses monl.Report into store.
func (w *worker) parseMonlReport(ctx context.Context, target model.Monl, report monl.Report) error {
	// Save stat of release kind
	saved, err := w.stats.List(ctx, &store.StatOpts{
		MonlID: target.ID,
		Kind:   report.Latest().Group(),
	})
	if err != nil {
		return err
	}
	for report.Next() {
		rs := report.Stat()
		logx.WithField("value", rs.Value()).Debugf("queue worker: stat")
		isLatest := rs.Value() == report.Latest().Value()
		sm := model.StatList(saved).
			FindKind(rs.Group()).
			FindValue(rs.Value())
		if len(sm) > 0 {
			so := sm[0]
			if isLatest {
				break
			}
			if !so.IsLatest {
				break
			}
			so.IsLatest = false
			err = w.stats.Update(ctx, &so)
			if err != nil {
				return err
			}
			break
		}

		s := model.Stat{
			MonlID:     target.ID,
			RecordedAt: (field.Time)(rs.Date()),
			Kind:       rs.Group(),
			Value:      rs.Value(),
			IsLatest:   isLatest,
			Manifest:   rs.Manifest(),
		}
		err = w.stats.Create(ctx, &s)
		if err != nil {
			return err
		}
	}

	// Save stats of popularity
	sl, err := w.stats.List(ctx, &store.StatOpts{
		MonlID:     target.ID,
		WithLatest: true,
	})
	for _, sp := range report.Popularity() {
		sm := model.StatList(sl).FindKind(sp.Group())
		if len(sm) > 0 {
			so := sm[0]
			so.IsLatest = false
			err = w.stats.Update(ctx, &so)
			if err != nil {
				return err
			}
		}
		s := model.Stat{
			MonlID:     target.ID,
			RecordedAt: (field.Time)(sp.Date()),
			Kind:       sp.Group(),
			Value:      sp.Value(),
			IsLatest:   true,
		}
		err = w.stats.Create(ctx, &s)
		if err != nil {
			return err
		}
	}

	return nil
}

// jobStopped wraps error into the payload and tells manager for status update.
func (w *worker) jobStopped(ctx context.Context, job *Job, msg error) error {
	job.Error = msg
	err := w.manager.jobStopped(ctx, job)
	if err != nil {
		return err
	}
	return msg
}
