package response

import (
	"encoding/json"
	"net/http"
)

type Body map[string]interface{}

func Error(err error) Body {
	return Body{"error": err.Error()}
}

func JSON(w http.ResponseWriter, v interface{}, code int) error {
	if code > 0 {
		w.WriteHeader(code)
	}
	enc := json.NewEncoder(w)
	switch v.(type) {
	case error:
		return enc.Encode(Error(v.(error)))
	case nil:
		return nil
	default:
		return enc.Encode(v)
	}
}

func IsError(code int) bool {
	if 0 < code && code < 200 {
		return true
	}
	if 400 <= code {
		return true
	}
	return false
}

type PageInfo struct {
	TotalCount int64 `json:"totalCount"`
	Page       int64 `json:"page"`
	PageSize   int64 `json:"pageSize"`
}

func ListJSON(w http.ResponseWriter, v interface{}, info *PageInfo, code int) error {
	body := struct {
		*PageInfo
		Data interface{} `json:"data"`
	}{info, v}
	return JSON(w, body, code)
}
