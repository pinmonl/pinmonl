package apibody

import "github.com/pinmonl/pinmonl/model"

// Tag defines the response body of Tag.
type Tag struct {
	model.Tag
}

// NewTag creates response body.
func NewTag(t model.Tag) Tag {
	return Tag{Tag: t}
}
