package session

import (
	"net/http"
)

// Store generalizes the session store.
type Store interface {
	Get(*http.Request) (*Values, error)
	Set(http.ResponseWriter, *Values) error
	Del(http.ResponseWriter, *http.Request) error
}

// Values contains the data for session payload.
type Values struct {
	UserID string
}
