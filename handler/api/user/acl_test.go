package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	mockUser := model.User{ID: "user-test-id"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(
		request.WithUser(r.Context(), mockUser),
	)

	Authorize()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMatchUser(t *testing.T) {
	mockUser := model.User{ID: "test-user-id"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(
		request.WithUser(r.Context(), mockUser),
	)

	assert.True(t, MatchUser(w, r, mockUser.ID))
}

func TestGuest(t *testing.T) {
	mockUser := model.User{ID: "test-user-id"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h := Guest()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r = r.WithContext(
		request.WithUser(r.Context(), mockUser),
	)
	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
