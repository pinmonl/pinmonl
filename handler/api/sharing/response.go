package sharing

import (
	"github.com/pinmonl/pinmonl/handler/api/tag"
	"github.com/pinmonl/pinmonl/handler/api/user"
	"github.com/pinmonl/pinmonl/model"
)

// Body defines the response body of Share.
type Body struct {
	model.Share
	Owner    *user.Body  `json:"owner,omitempty"`
	MustTags *[]tag.Body `json:"mustTags,omitempty"`
}

// NewBody creates the response body.
func NewBody(s model.Share) Body {
	b := Body{
		Share: s,
	}
	return b
}

// WithOwner sets value of owner.
func (b Body) WithOwner(u model.User) Body {
	ub := user.NewBody(u)
	b.Owner = &ub
	return b
}

// WithMustTags sets value of must tags.
func (b Body) WithMustTags(ts []model.Tag) Body {
	tbs := make([]tag.Body, len(ts))
	for i, t := range ts {
		tbs[i] = tag.NewBody(t)
	}
	b.MustTags = &tbs
	return b
}
