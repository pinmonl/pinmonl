package payload

import (
	"encoding/json"
	"io"
)

// UnmarshalJSON reads json payload from reader.
func UnmarshalJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// MarshalJSON writes json payload into writer.
func MarshalJSON(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}
