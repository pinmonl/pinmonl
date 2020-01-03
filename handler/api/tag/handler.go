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
		ms, err := tags.List(r.Context(), &store.TagOpts{})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := make([]interface{}, len(ms))
		for i, m := range ms {
			resp[i] = Resp(m)
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

		var m model.Tag
		in.Fill(&m)
		err = tags.Create(r.Context(), &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, Resp(m))
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
		err = tags.Update(r.Context(), &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, Resp(m))
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
