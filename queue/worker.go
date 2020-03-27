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

	store    store.Store
	monls    store.MonlStore
	pkgs     store.PkgStore
	pinls    store.PinlStore
	stats    store.StatStore
	substats store.SubstatStore
}

// run is goroutine which listens to job queue.
func (w *worker) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
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
		logx.Fatalf("queue worker: fails to start tx, err: %s", err)
		return err
	}

	switch job.Detail.Name {
	case model.JobPinlCreated:
	case model.JobPinlUpdated:
	case model.JobPkgCreated:
	case model.JobPkgRegularUpdate:
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
