package request

import (
	"encoding/json"
	"net/http"
)

func JSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
