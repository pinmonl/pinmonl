package queue

import (
	"context"
	"fmt"
	"time"

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
	monls store.MonlStore
	pkgs  store.PkgStore
	pinls store.PinlStore
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

	switch job.Detail.Name {
	case model.JobPinlCreated:
		err = w.processPinlCreated(ctx, job)
	case model.JobPinlUpdated:
		err = w.processPinlCreated(ctx, job)
	case model.JobPkgCreated:
		err = w.processPkgCreated(ctx, job)
	case model.JobPkgRegularUpdate:
		err = w.processPkgRegularUpdate(ctx, job)
	default:
		defer w.store.Commit(ctx)
		return w.jobStopped(ctx, job, fmt.Errorf("queue worker: job name is not defined"))
	}

	if err != nil {
		logx.Errorf("queue worker: err: %s", err)
		w.store.Rollback(ctx)
		ctx, _ = w.store.EndTx(ctx)
		ctx, _ = w.store.BeginTx(ctx)
		defer w.store.Commit(ctx)
		return w.jobStopped(ctx, job, err)
	}
	defer w.store.Commit(ctx)
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

	tm := model.Monl{URL: t.URL}
	err = w.monls.Create(ctx, &tm)
	if err != nil {
		return err
	}

	vendors := w.monl.GuessURL(t.URL)
	for _, v := range vendors {
		logx.Debugf("queue worker: working on monl vendor - %s", v.Name())
		ml, err := w.pkgs.List(ctx, &store.PkgOpts{
			MonlID: tm.ID,
			Vendor: v.Name(),
		})
		if err != nil {
			return err
		}
		if len(ml) > 0 {
			continue
		}

		m := model.Pkg{MonlID: tm.ID, Vendor: v.Name()}
		err = w.pkgs.Create(ctx, &m)
		if err != nil {
			return err
		}

		err = w.manager.Enqueue(ctx, &model.Job{
			Name:     model.JobPkgCreated,
			TargetID: m.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// processPkgCreated handles the hook after pkg created.
func (w *worker) processPkgCreated(ctx context.Context, job *Job) error {
	t := model.Pkg{ID: job.Detail.TargetID}
	err := w.pkgs.Find(ctx, &t)
	if err != nil {
		return err
	}
	tm := model.Monl{ID: t.MonlID}
	err = w.monls.Find(ctx, &tm)
	if err != nil {
		return err
	}

	v, err := w.monl.Get(t.Vendor)
	if err != nil {
		return err
	}

	r, err := v.Load(ctx, tm.URL)
	if err != nil {
		return err
	}
	defer r.Close()

	t.URL = r.RawURL()
	t.VendorURI = r.URI()
	err = w.pkgs.Update(ctx, &t)
	if err != nil {
		return err
	}

	err = w.parseMonlReport(ctx, t, r)
	if err != nil {
		return err
	}

	schd := w.makeScheduleJobForPkg(*job)
	err = w.scheduleJob(ctx, &schd)
	if err != nil {
		return err
	}
	return nil
}

// processPkgRegularUpdate handles regular update of pkg.
func (w *worker) processPkgRegularUpdate(ctx context.Context, job *Job) error {
	t := model.Pkg{ID: job.Detail.TargetID}
	err := w.pkgs.Find(ctx, &t)
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

	schd := w.makeScheduleJobForPkg(*job)
	err = w.scheduleJob(ctx, &schd)
	if err != nil {
		return err
	}
	return nil
}

// parseMonlReport parses monl.Report into store.
func (w *worker) parseMonlReport(ctx context.Context, target model.Pkg, report monl.Report) error {
	if report.Latest() != nil {
		// Save stat of release kind
		saved, err := w.stats.List(ctx, &store.StatOpts{
			PkgID: target.ID,
			Kind:  report.Latest().Group(),
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
				PkgID:      target.ID,
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
	}

	// Save stats of popularity
	sl, err := w.stats.List(ctx, &store.StatOpts{
		PkgID:      target.ID,
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
			PkgID:      target.ID,
			RecordedAt: (field.Time)(sp.Date()),
			Kind:       sp.Group(),
			Value:      sp.Value(),
			IsLatest:   true,
			Manifest:   sp.Manifest(),
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

func (w *worker) scheduleJob(ctx context.Context, job *Job) error {
	return w.manager.scheduleJob(ctx, job)
}

func (w *worker) makeScheduleJobForPkg(orig Job) Job {
	j := Job{}
	j.Detail.Status = model.JobStatusPending
	j.Detail.Name = model.JobPkgRegularUpdate
	j.Detail.TargetID = orig.Detail.TargetID

	var at time.Time
	intv := time.Hour
	if orig.Detail.ScheduledAt.Time().IsZero() {
		at = time.Now()
		at = at.Truncate(time.Hour)
		at = at.Add(intv + time.Hour)
	} else {
		at = orig.Detail.ScheduledAt.Time()
		at = at.Add(intv)
	}
	j.Detail.ScheduledAt = field.Time(at)
	return j
}
