package tag

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestRequireOwner(t *testing.T) {
	mockUser := model.User{ID: "user-test-id"}
	mockTag := model.Tag{ID: "tag-test-id", UserID: "user-test-id"}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(
		request.WithTag(
			request.WithUser(r.Context(), mockUser),
			mockTag,
		),
	)
	RequireOwner()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
