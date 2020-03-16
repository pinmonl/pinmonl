package share

import "github.com/pinmonl/pinmonl/model"

// Body defines the response body of Share.
type Body struct {
	model.Share
	MustTags []string  `json:"mustTags"`
	AnyTags  *[]string `json:"anyTags,omitempty"`
}

// NewBody creates the response body.
func NewBody(s model.Share) Body {
	return Body{
		Share:    s,
		MustTags: make([]string, 0),
	}
}

// WithMustTags sets value of Body.MustTags.
func (b Body) WithMustTags(mts []model.Tag) Body {
	b.MustTags = append(b.MustTags, (model.TagList)(mts).PluckName()...)
	return b
}

// WithAnyTags sets value of Body.AnyTags.
func (b Body) WithAnyTags(ats []model.Tag) Body {
	tns := make([]string, 0)
	tns = append(tns, (model.TagList)(ats).PluckName()...)
	b.AnyTags = &tns
	return b
}
