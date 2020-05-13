package pinl

import (
	"github.com/pinmonl/pinmonl/handler/api/pkg"
	"github.com/pinmonl/pinmonl/model"
)

// Body defines the response body of Pinl.
type Body struct {
	model.Pinl
	Tags []string   `json:"tags"`
	Pkgs []pkg.Body `json:"pkgs"`
}

// NewBody creates the response body.
func NewBody(p model.Pinl) Body {
	return Body{
		Pinl: p,
		Tags: make([]string, 0),
		Pkgs: make([]pkg.Body, 0),
	}
}

// WithTags sets value of Body.Tags.
func (b Body) WithTags(ts []model.Tag) Body {
	b.Tags = append(b.Tags, (model.TagList)(ts).PluckName()...)
	return b
}

// WithPkgs sets value of pkgs.
func (b Body) WithPkgs(ps []model.Pkg, statMap map[string][]model.Stat) Body {
	b.Pkgs = make([]pkg.Body, len(ps))
	for i, p := range ps {
		pBody := pkg.NewBody(p)
		if statMap != nil {
			pBody = pBody.WithStats(statMap[p.ID]...)
		}
		b.Pkgs[i] = pBody
	}
	return b
}
