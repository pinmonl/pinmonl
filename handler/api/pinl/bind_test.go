package pinl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database/dbtest"
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

	paramName := "pinl"
	mockPinl := model.Pinl{}
	pinls := store.NewPinlStore(store.NewStore(db))
	pinls.Create(context.TODO(), &mockPinl)

	c := &chi.Context{}
	c.URLParams.Add(paramName, mockPinl.ID)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(
		context.WithValue(r.Context(), chi.RouteCtxKey, c),
	)

	BindByID(paramName, pinls)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
