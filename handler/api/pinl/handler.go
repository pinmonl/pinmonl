package pinl

import (
	"bytes"
	"context"
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/apibody"
	"github.com/pinmonl/pinmonl/handler/api/apiutils"
	"github.com/pinmonl/pinmonl/handler/api/image"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/handler/middleware"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkg/scrape"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/pubsub/message"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
)

// HandleList returns pinls.
func HandleList(pinls store.PinlStore, taggables store.TaggableStore, monpkgs store.MonpkgStore, stats store.StatStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		p := middleware.PaginationFrom(ctx)
		u, _ := request.UserFrom(ctx)
		ms, err := pinls.List(ctx, &store.PinlOpts{UserID: u.ID, ListOpts: *p})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		tsm := map[string][]model.Tag{}
		if len(ms) > 0 {
			mps := model.MustBeMorphables(ms)
			tsm, err = taggables.ListTags(ctx, &store.TaggableOpts{Targets: mps})
			if err != nil {
				response.InternalError(w, err)
				return
			}
		}

		pkgMap, statMap, err := apiutils.ListPinlStats(ctx, monpkgs, stats, ms...)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := make([]apibody.Pinl, len(ms))
		for i, m := range ms {
			resp[i] = apibody.NewPinl(m).
				WithTags(tsm[m.ID]).
				WithPkgs(pkgMap[m.MonlID], statMap)
		}
		response.JSON(w, resp)
	}
}

// HandleFind returns pinl and its relations.
func HandleFind(taggables store.TaggableStore, monpkgs store.MonpkgStore, pkgs store.PkgStore, stats store.StatStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, has := request.PinlFrom(ctx)
		if !has {
			response.NotFound(w, nil)
			return
		}

		ts, err := taggables.ListTags(ctx, &store.TaggableOpts{Target: m})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		pkgMap, statMap, err := apiutils.ListPinlStats(ctx, monpkgs, stats, m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := apibody.NewPinl(m).
			WithTags(ts[m.ID]).
			WithPkgs(pkgMap[m.MonlID], statMap)
		response.JSON(w, resp)
	}
}

// HandleCreate validates and create pinl from user input.
func HandleCreate(
	pinls store.PinlStore,
	tags store.TagStore,
	taggables store.TaggableStore,
	dp *queue.Dispatcher,
	images store.ImageStore,
	pkgs store.PkgStore,
	stats store.StatStore,
	pubsub *pubsub.Server,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		if err = in.Validate(); err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		m := model.Pinl{UserID: u.ID}
		err = in.Fill(&m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = fillCardIfEmpty(ctx, images, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = pinls.Create(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		ts, err := apiutils.FindOrCreateTagsByName(ctx, tags, u, in.Tags)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = taggables.ReAssocTags(ctx, m, ts)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := apibody.NewPinl(m).WithTags(ts)
		go func() {
			pubsub.Publish(message.NewPinlCreateMessage(resp))
			dp.SyncPinl(m)
		}()
		response.JSON(w, resp)
	}
}

// HandleUpdate validates and updates pinl from user input.
func HandleUpdate(
	pinls store.PinlStore,
	tags store.TagStore,
	taggables store.TaggableStore,
	dp *queue.Dispatcher,
	images store.ImageStore,
	pkgs store.PkgStore,
	monpkgs store.MonpkgStore,
	stats store.StatStore,
	pubsub *pubsub.Server,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		if err = in.Validate(); err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		m, _ := request.PinlFrom(ctx)
		err = in.Fill(&m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = fillCardIfEmpty(ctx, images, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = pinls.Update(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		ts, err := apiutils.FindOrCreateTagsByName(ctx, tags, u, in.Tags)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = taggables.ReAssocTags(ctx, m, ts)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		pkgMap, statMap, err := apiutils.ListPinlStats(ctx, monpkgs, stats, m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := apibody.NewPinl(m).
			WithTags(ts).
			WithPkgs(pkgMap[m.MonlID], statMap)
		go func() {
			pubsub.Publish(message.NewPinlUpdateMessage(resp))
			dp.SyncPinl(m)
		}()
		response.JSON(w, resp)
	}
}

// HandleDelete removes pinl and its relations.
func HandleDelete(
	pinls store.PinlStore,
	taggables store.TaggableStore,
	pubsub *pubsub.Server,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, _ := request.PinlFrom(ctx)

		err := pinls.Delete(ctx, &m)
		if err != nil {
			response.InternalError(w, nil)
			return
		}

		err = taggables.ClearTags(ctx, m)
		if err != nil {
			response.InternalError(w, nil)
			return
		}

		resp := apibody.NewPinl(m)
		pubsub.Publish(message.NewPinlDeleteMessage(resp))
		response.NoContent(w)
	}
}

// HandlePageInfo returns the page info of Pinl.
func HandlePageInfo(pinls store.PinlStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		count, err := pinls.Count(ctx, &store.PinlOpts{UserID: u.ID})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, response.NewPageInfo(count))
	}
}

func fillCardIfEmpty(ctx context.Context, images store.ImageStore, m *model.Pinl) error {
	if m.Title != "" {
		return nil
	}

	resp, err := scrape.Get(m.URL)
	if err != nil {
		return err
	}
	card, err := resp.Card()
	if err != nil {
		return err
	}
	ci, err := card.Image()
	if err != nil {
		return err
	}

	m2 := *m
	m2.Title = card.Title()
	m2.Description = card.Description()
	if ci != nil {
		img, err := image.UploadFromReader(ctx, images, bytes.NewBuffer(ci))
		if err != nil {
			return err
		}
		m2.ImageID = img.ID
	}
	*m = m2
	return nil
}
