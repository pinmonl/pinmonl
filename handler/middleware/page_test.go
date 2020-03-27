package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h := Pagination(10)

	h(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			p := PaginationFrom(r.Context())
			assert.NotNil(t, p)
			assert.Equal(t, int64(0), p.Offset)
			assert.Equal(t, int64(10), p.Limit)
		},
	)).ServeHTTP(w, r)

	r = httptest.NewRequest("GET", "/?page=10", nil)
	h(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			p := PaginationFrom(r.Context())
			assert.NotNil(t, p)
			assert.Equal(t, int64(90), p.Offset)
			assert.Equal(t, int64(10), p.Limit)
		},
	)).ServeHTTP(w, r)

	r = httptest.NewRequest("GET", "/?page=10&pageSize=100", nil)
	h(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			p := PaginationFrom(r.Context())
			assert.NotNil(t, p)
			assert.Equal(t, int64(900), p.Offset)
			assert.Equal(t, int64(100), p.Limit)
		},
	)).ServeHTTP(w, r)

	r = httptest.NewRequest("GET", "/?page=10&pageSize=200", nil)
	h(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			p := PaginationFrom(r.Context())
			assert.NotNil(t, p)
			assert.Equal(t, int64(1800), p.Offset)
			assert.Equal(t, int64(200), p.Limit)
		},
	)).ServeHTTP(w, r)
}
