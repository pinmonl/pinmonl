package tag

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
	"github.com/stretchr/testify/assert"
)

func TestBindByID(t *testing.T) {
	db, err := dbtest.Open()
	assert.Nil(t, err)
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	paramName := "tag"
	mockTag := model.Tag{Name: "tag-test-id"}
	tags := store.NewTagStore(store.NewStore(db))
	tags.Create(context.TODO(), &mockTag)

	c := &chi.Context{}
	c.URLParams.Add(paramName, mockTag.ID)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(
		context.WithValue(
			request.WithTag(r.Context(), mockTag),
			chi.RouteCtxKey,
			c,
		),
	)

	BindByID(paramName, tags)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
