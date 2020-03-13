package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	db, err := dbtest.Open()
	assert.Nil(t, err)
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	st := store.NewStore(db)
	users := store.NewUserStore(st)
	t.Run("Create", testHandleCreate(users))
	t.Run("GetMe", testHandleGetMe(users))
	t.Run("UpdateMe", testHandleUpdateMe(users))
}

func testHandleCreate(users store.UserStore) func(*testing.T) {
	return func(t *testing.T) {
		var (
			w    *httptest.ResponseRecorder
			r    *http.Request
			body = &bytes.Buffer{}
		)

		w = httptest.NewRecorder()
		body.WriteString(`{"login":"user1","password":"pw"}`)
		r = httptest.NewRequest("POST", "/", body)
		HandleCreate(users)(w, r)
		assert.Equal(t, 400, w.Code)

		w = httptest.NewRecorder()
		body.Reset()
		body.WriteString(`{"email":"email","password":"pw"}`)
		r = httptest.NewRequest("POST", "/", body)
		HandleCreate(users)(w, r)
		assert.Equal(t, 400, w.Code)

		w = httptest.NewRecorder()
		body.Reset()
		body.WriteString(`{"email":"email@email.com","password":"pw"}`)
		r = httptest.NewRequest("POST", "/", body)
		HandleCreate(users)(w, r)
		assert.Equal(t, 400, w.Code)

		w = httptest.NewRecorder()
		body.Reset()
		body.WriteString(`{"name":"user name 1","password":"pw"}`)
		r = httptest.NewRequest("POST", "/", body)
		HandleCreate(users)(w, r)
		assert.Equal(t, 400, w.Code, w.Body)

		w = httptest.NewRecorder()
		body.Reset()
		body.WriteString(`{"login":"user1","email":"user1@email.com","name":"user name 1","password":"pw12345678"}`)
		r = httptest.NewRequest("POST", "/", body)
		HandleCreate(users)(w, r)
		assert.Equal(t, 200, w.Code, w.Body)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/", body)
		HandleCreate(users)(w, r)
		assert.Equal(t, 400, w.Code, w.Body)
	}
}

func testHandleGetMe(users store.UserStore) func(*testing.T) {
	return func(t *testing.T) {
		user := model.User{Login: "user1"}
		users.FindLogin(context.TODO(), &user)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r = r.WithContext(
			request.WithUser(r.Context(), user),
		)

		HandleGetMe()(w, r)
		resp := make(map[string]string)
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, user.ID, resp["id"], w.Body)
	}
}

func testHandleUpdateMe(users store.UserStore) func(*testing.T) {
	return func(t *testing.T) {
		user := model.User{Login: "user1"}
		users.FindLogin(context.TODO(), &user)

		body := &bytes.Buffer{}
		body.WriteString(`{"login":"user1b","email":"user1b@email.com","password":"87654321pw","name":"user name 1b"}`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/", body)
		r = r.WithContext(
			request.WithUser(r.Context(), user),
		)

		HandleUpdateMe(users)(w, r)
		resp := make(map[string]string)
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, 200, w.Code, w.Body)
		assert.Equal(t, "user1b", resp["login"])
		assert.Equal(t, "user1b@email.com", resp["email"])
		assert.Equal(t, "user name 1b", resp["name"])
	}
}
