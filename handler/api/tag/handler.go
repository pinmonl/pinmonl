package tag

import (
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// HandleList returns tags.
func HandleList(tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		ms, err := tags.List(ctx, &store.TagOpts{UserID: u.ID})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := make([]interface{}, len(ms))
		for i, m := range ms {
			resp[i] = NewBody(m)
		}
		response.JSON(w, resp)
	}
}

// HandleCreate validates and creates tag from user input.
func HandleCreate(tags store.TagStore) http.HandlerFunc {
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
		m := model.Tag{UserID: u.ID}
		in.Fill(&m)

		if m.ParentID != "" {
			pt := model.Tag{ID: m.ParentID}
			err := tags.Find(ctx, &pt)
			if err != nil {
				response.BadRequest(w, err)
				return
			}
			m.Level = pt.Level + 1
		}

		err = tags.Create(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, NewBody(m))
	}
}

// HandleUpdate validates and updates tag from user input.
func HandleUpdate(tags store.TagStore) http.HandlerFunc {
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
		m, _ := request.TagFrom(ctx)
		in.Fill(&m)

		if m.ParentID != "" {
			pt := model.Tag{ID: m.ParentID}
			err := tags.Find(ctx, &pt)
			if err != nil {
				response.BadRequest(w, err)
				return
			}
			m.Level = pt.Level + 1
		}

		err = tags.Update(r.Context(), &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, NewBody(m))
	}
}

// HandleDelete removes tag.
func HandleDelete(tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, _ := request.TagFrom(ctx)

		err := tags.Delete(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.NoContent(w)
	}
}

// HandlePageInfo returns the page info of Tag.
func HandlePageInfo(tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		count, err := tags.Count(ctx, &store.TagOpts{UserID: u.ID})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, response.NewPageInfo(count))
	}
}
