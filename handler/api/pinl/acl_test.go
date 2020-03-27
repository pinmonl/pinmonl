package pinl

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestRequireOwner(t *testing.T) {
	mockUser := model.User{ID: "test-user-id"}
	mockPinl := model.Pinl{ID: "test-pinl-id", UserID: mockUser.ID}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(
		request.WithUser(
			request.WithPinl(r.Context(), mockPinl),
			mockUser,
		),
	)

	RequireOwner()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
