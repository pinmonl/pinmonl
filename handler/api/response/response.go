package response

import (
	"net/http"
)

// Body defines the top level response body structure.
type Body map[string]interface{}

// BadRequest returns 400 HTTP status code.
func BadRequest(w http.ResponseWriter, v interface{}) {
	w.WriteHeader(400)
	JSON(w, v)
}

// NotFound returns 404 HTTP status code.
func NotFound(w http.ResponseWriter, v interface{}) {
	w.WriteHeader(404)
	JSON(w, v)
}

// InternalError returns 500 HTTP status code.
func InternalError(w http.ResponseWriter, v interface{}) {
	w.WriteHeader(500)
	JSON(w, v)
}

// NoContent returns 204 HTTP status code.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(204)
}
