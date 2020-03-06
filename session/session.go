package session

import (
	"encoding/json"
	"net/http"
)

// Store generalizes the session store.
type Store interface {
	Get(*http.Request) (*Values, error)
	Set(http.ResponseWriter, *Values) (*Response, error)
	Del(http.ResponseWriter, *http.Request) error
}

// Values contains the data for session payload.
type Values struct {
	UserID string
}

// Response defines the payload of session data.
type Response struct {
	Data map[string]interface{}
}

// MarshalJSON implement json.Marshaler.
func (r *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Data)
}
