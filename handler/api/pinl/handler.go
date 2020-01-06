package pinl

import (
	"bytes"
	"context"
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/image"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkg/scrape"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/store"
)

// HandleList returns pinls.
func HandleList(pinls store.PinlStore, tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		ms, err := pinls.List(ctx, &store.PinlOpts{UserID: u.ID})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		mps := model.MustBeMorphables(ms)
		ts, err := tags.List(ctx, &store.TagOpts{Targets: mps})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := make([]interface{}, len(ms))
		for i, m := range ms {
			mt := (model.TagList)(ts).FindMorphable(m)
			resp[i] = Resp(m, mt)
		}
		response.JSON(w, resp)
	}
}

// HandleFind returns pinl and its relations.
func HandleFind(tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, has := request.PinlFrom(ctx)
		if !has {
			response.NotFound(w, nil)
			return
		}

		ts, err := tags.List(ctx, &store.TagOpts{Target: m})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, DetailResp(m, ts))
	}
}

// HandleCreate validates and create pinl from user input.
func HandleCreate(
	pinls store.PinlStore,
	tags store.TagStore,
	taggables store.TaggableStore,
	qm *queue.Manager,
	images store.ImageStore,
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

		ts, err := findOrCreateTags(ctx, tags, "", in.Tags)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = taggables.ReAssocTags(ctx, m, ts)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		err = qm.Enqueue(ctx, &model.Job{
			Name:     model.JobPinlCreated,
			TargetID: m.ID,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, DetailResp(m, ts))
	}
}

// HandleUpdate validates and updates pinl from user input.
func HandleUpdate(
	pinls store.PinlStore,
	tags store.TagStore,
	taggables store.TaggableStore,
	qm *queue.Manager,
	images store.ImageStore,
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

		ts, err := findOrCreateTags(ctx, tags, "", in.Tags)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = taggables.ReAssocTags(ctx, m, ts)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		err = qm.Enqueue(ctx, &model.Job{
			Name:     model.JobPinlUpdated,
			TargetID: m.ID,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, DetailResp(m, ts))
	}
}

// HandleDelete removes pinl and its relations.
func HandleDelete(
	pinls store.PinlStore,
	taggables store.TaggableStore,
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

		response.NoContent(w)
	}
}

func findOrCreateTags(ctx context.Context, tags store.TagStore, user string, tagNames []string) ([]model.Tag, error) {
	var ts []model.Tag
	for _, n := range tagNames {
		var t model.Tag
		found, err := tags.List(ctx, &store.TagOpts{Name: n, UserID: user})
		if err != nil {
			return nil, err
		}
		if len(found) == 0 {
			t2 := model.Tag{Name: n, UserID: user}
			err = tags.Create(ctx, &t2)
			if err != nil {
				return nil, err
			}
			t = t2
		} else {
			t = found[0]
		}
		ts = append(ts, t)
	}
	return ts, nil
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
	img, err := image.UploadFromReader(ctx, images, bytes.NewBuffer(ci))
	if err != nil {
		return err
	}

	m2.Title = card.Title()
	m2.Description = card.Description()
	m2.ImageID = img.ID
	*m = m2
	return nil
}
