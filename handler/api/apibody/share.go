package apibody

import "github.com/pinmonl/pinmonl/model"

// Share defines the response body of Share.
type Share struct {
	model.Share
	MustTags []string  `json:"mustTags"`
	AnyTags  *[]string `json:"anyTags,omitempty"`
}

// NewShare creates the response body.
func NewShare(s model.Share) Share {
	return Share{
		Share:    s,
		MustTags: make([]string, 0),
	}
}

// WithMustTags sets value of MustTags.
func (s Share) WithMustTags(mts []model.Tag) Share {
	s.MustTags = append(s.MustTags, (model.TagList)(mts).PluckName()...)
	return s
}

// WithAnyTags sets value of AnyTags.
func (s Share) WithAnyTags(ats []model.Tag) Share {
	tns := make([]string, 0)
	tns = append(tns, (model.TagList)(ats).PluckName()...)
	s.AnyTags = &tns
	return s
}
