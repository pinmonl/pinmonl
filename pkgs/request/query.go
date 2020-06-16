package request

import (
	"net/http"
	"strings"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type PinlQuery struct {
	Query string
	Tags  []string
}

func ParsePinlQuery(r *http.Request) (*PinlQuery, error) {
	query := PinlQuery{
		Query: r.URL.Query().Get("q"),
		Tags:  getQueryCsv(r, "tags"),
	}
	return &query, nil
}

type TagQuery struct {
	Query string
}

func ParseTagQuery(r *http.Request) (*TagQuery, error) {
	query := TagQuery{
		Query: r.URL.Query().Get("q"),
	}
	return &query, nil
}

type StatQuery struct {
	Kind   field.NullValue
	Latest field.NullBool
}

func ParseStatQuery(r *http.Request) (*StatQuery, error) {
	query := StatQuery{
		Latest: getQueryBool(r, "latest"),
	}
	if kindq := r.URL.Query().Get("kind"); kindq != "" {
		query.Kind = field.NewNullValue(model.StatKind(kindq))
	}
	return &query, nil
}

func getQueryCsv(r *http.Request, paramName string) []string {
	out := make([]string, 0)
	qv := r.URL.Query().Get(paramName)
	if qv == "" {
		return out
	}

	for _, val := range strings.Split(qv, ",") {
		if tv := strings.TrimSpace(val); tv != "" {
			out = append(out, tv)
		}
	}
	return out
}

func getQueryBool(r *http.Request, paramName string) field.NullBool {
	qv := r.URL.Query().Get(paramName)
	if qv == "" {
		return field.NullBool{}
	}

	qv = strings.ToLower(qv)
	if qv == "1" || qv == "true" {
		return field.NewNullBool(true)
	}
	return field.NewNullBool(false)
}
