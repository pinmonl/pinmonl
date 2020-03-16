package pinl

import (
	"github.com/pinmonl/pinmonl/model"
)

// Body defines the response body of Pinl.
type Body struct {
	model.Pinl
	Tags []string `json:"tags"`
}

// NewBody creates the response body.
func NewBody(p model.Pinl) Body {
	return Body{
		Pinl: p,
		Tags: make([]string, 0),
	}
}

// WithTags sets value of Body.Tags.
func (b Body) WithTags(ts []model.Tag) Body {
	b.Tags = append(b.Tags, (model.TagList)(ts).PluckName()...)
	return b
}
