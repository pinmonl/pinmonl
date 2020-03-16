package tag

import "github.com/pinmonl/pinmonl/model"

// Body defines the response body of Tag.
type Body struct {
	model.Tag
}

// NewBody creates response body.
func NewBody(t model.Tag) Body {
	return Body{Tag: t}
}
