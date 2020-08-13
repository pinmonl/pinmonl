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
	NoTag field.NullBool
}

func ParsePinlQuery(r *http.Request) (*PinlQuery, error) {
	query := PinlQuery{
		Query: r.URL.Query().Get("q"),
		Tags:  QueryCsv(r, "tag"),
		NoTag: QueryBool(r, "notag"),
	}
	return &query, nil
}

type TagQuery struct {
	Query     string
	Names     []string
	ParentIDs []string
}

func ParseTagQuery(r *http.Request) (*TagQuery, error) {
	query := TagQuery{
		Query: r.URL.Query().Get("q"),
		Names: QueryCsv(r, "name"),
	}
	if parentq := r.URL.Query()["parent"]; len(parentq) > 0 {
		if parentq[0] == "" {
			query.ParentIDs = []string{""}
		} else {
			query.ParentIDs = QueryCsv(r, "parent")
		}
	}
	return &query, nil
}

type StatQuery struct {
	PkgIDs    []string
	Kinds     []model.StatKind
	Latest    field.NullBool
	ParentIDs []string
}

func ParseStatQuery(r *http.Request) (*StatQuery, error) {
	query := StatQuery{
		PkgIDs:    QueryCsv(r, "pkg"),
		Latest:    QueryBool(r, "latest"),
		ParentIDs: QueryCsv(r, "parent"),
	}
	for _, kindstr := range QueryCsv(r, "kind") {
		query.Kinds = append(query.Kinds, model.StatKind(kindstr))
	}
	return &query, nil
}

func QueryCsv(r *http.Request, paramName string) []string {
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

func QueryBool(r *http.Request, paramName string) field.NullBool {
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
