package apibody

import "github.com/pinmonl/pinmonl/model"

// Sharing defines the response body of Sharing.
type Sharing struct {
	model.Share
	Owner    *User  `json:"owner,omitempty"`
	MustTags *[]Tag `json:"mustTags,omitempty"`
}

// NewSharing creates the response body.
func NewSharing(s model.Share) Sharing {
	return Sharing{Share: s}
}

// WithOwner sets value of owner.
func (s Sharing) WithOwner(u model.User) Sharing {
	ub := NewUser(u)
	s.Owner = &ub
	return s
}

// WithMustTags sets value of must tags.
func (s Sharing) WithMustTags(ts []model.Tag) Sharing {
	tbs := make([]Tag, len(ts))
	for i, t := range ts {
		tbs[i] = NewTag(t)
	}
	s.MustTags = &tbs
	return s
}
