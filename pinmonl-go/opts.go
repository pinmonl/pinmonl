package pinmonl

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/pinmonl/pinmonl/model/field"
)

type ListOpts struct {
	Page int
	Size int
}

func (l ListOpts) AppendTo(val url.Values) {
	if l.Page > 0 {
		val.Add("page", strconv.Itoa(l.Page))
	}
	if l.Size > 0 {
		val.Add("page_size", strconv.Itoa(l.Size))
	}
}

type StatLatestListOpts struct {
	ListOpts
	Kind string
}

func (s StatLatestListOpts) Encode() string {
	qs := url.Values{}
	s.ListOpts.AppendTo(qs)
	if s.Kind != "" {
		qs.Add("kind", s.Kind)
	}
	return qs.Encode()
}

type StatListOpts struct {
	ListOpts
	Kind   string
	Latest field.NullBool
}

func (s StatListOpts) Encode() string {
	qs := url.Values{}
	s.ListOpts.AppendTo(qs)
	if s.Latest.Valid {
		vl := "0"
		if s.Latest.Value() {
			vl = "1"
		}
		qs.Add("latest", vl)
	}
	if s.Kind != "" {
		qs.Add("kind", s.Kind)
	}
	return qs.Encode()
}

type PinlListOpts struct {
	ListOpts
	Query string
	Tags  []string
}

func (p PinlListOpts) Encode() string {
	qs := url.Values{}
	p.ListOpts.AppendTo(qs)
	if p.Query != "" {
		qs.Add("q", p.Query)
	}
	if len(p.Tags) > 0 {
		qs.Add("tags", strings.Join(p.Tags, ","))
	}
	return qs.Encode()
}

type TagListOpts struct {
	ListOpts
	Query string
}

func (t TagListOpts) Encode() string {
	qs := url.Values{}
	t.ListOpts.AppendTo(qs)
	if t.Query != "" {
		qs.Add("q", t.Query)
	}
	return qs.Encode()
}

type PkgListOpts struct {
	ListOpts
}

func (p PkgListOpts) Encode() string {
	qs := url.Values{}
	p.ListOpts.AppendTo(qs)
	return qs.Encode()
}
