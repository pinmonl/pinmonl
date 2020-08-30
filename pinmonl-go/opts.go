package pinmonl

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/pinmonl/pinmonl/model/field"
)

type ListOpts struct {
	Page int64
	Size int64
}

func (l ListOpts) AppendTo(val url.Values) {
	if l.Page > 0 {
		val.Add("page", strconv.FormatInt(l.Page, 10))
	}
	if l.Size > 0 {
		val.Add("page_size", strconv.FormatInt(l.Size, 10))
	} else if l.Size == -1 {
		val.Add("page_size", "0")
	}
}

type StatListOpts struct {
	ListOpts
	Kinds   []string
	Latest  field.NullBool
	Parents []string
	Pkgs    []string
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
	if len(s.Kinds) > 0 {
		qs.Add("kind", strings.Join(s.Kinds, ","))
	}
	if len(s.Parents) > 0 {
		qs.Add("parent", strings.Join(s.Parents, ","))
	}
	if len(s.Pkgs) > 0 {
		qs.Add("pkg", strings.Join(s.Pkgs, ","))
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
	URL string
}

func (p PkgListOpts) Encode() string {
	qs := url.Values{}
	p.ListOpts.AppendTo(qs)
	if p.URL != "" {
		qs.Add("url", p.URL)
	}
	return qs.Encode()
}
