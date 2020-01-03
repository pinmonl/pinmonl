package request

import (
	"encoding/json"
	"io"
)

// JSON decodes body and parses to interface.
func JSON(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(v)
	return err
}
