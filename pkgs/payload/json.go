package payload

import (
	"encoding/json"
	"io"
)

func JSONMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func JSONEncode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func JSONDecode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
