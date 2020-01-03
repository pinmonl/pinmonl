package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pinmonl/pinmonl/validate"
)

// JSON writes response body in json format.
func JSON(w http.ResponseWriter, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.Encode(wrap(v))
	w.Write(buf.Bytes())
}

func wrap(v interface{}) interface{} {
	switch v.(type) {
	case validate.Errors:
		es := v.(validate.Errors)
		return Body{"validations": es.Errors()}
	case error:
		return Body{"error": fmt.Sprintf("%s", v)}
	case Body:
		return v
	default:
		return v
	}
}
