package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/session"
	"github.com/pinmonl/pinmonl/store"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	db, err := dbtest.Open()
	assert.Nil(t, err)
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	mockUser := model.User{Login: "user1"}
	users := store.NewUserStore(store.NewStore(db))
	users.Create(context.TODO(), &mockUser)
	sess := mockSession{User: mockUser}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(
		request.WithUser(r.Context(), mockUser),
	)

	Authenticate(sess, users)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

// mockSession simulates session.Store.
type mockSession struct {
	User model.User
}

func (m mockSession) Get(_ *http.Request) (*session.Values, error) {
	return &session.Values{UserID: m.User.ID}, nil
}

func (m mockSession) Set(_ http.ResponseWriter, _ *session.Values) (*session.Response, error) {
	return nil, nil
}

func (m mockSession) Del(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
