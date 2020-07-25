package model

import (
	"github.com/pinmonl/pinmonl/model/field"
)

type Share struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ImageID     string     `json:"imageId"`
	Status      Status     `json:"status"`
	CreatedAt   field.Time `json:"createdAt"`
	UpdatedAt   field.Time `json:"updatedAt"`

	User         *User     `json:"user,omitempty"`
	MustTagNames *[]string `json:"mustTags,omitempty"`
	AnyTagNames  *[]string `json:"anyTags,omitempty"`
}

func (s Share) MorphKey() string  { return s.ID }
func (s Share) MorphName() string { return "share" }

func (s *Share) SetMustTagNames(tags TagList) {
	tn := tags.Names()
	s.MustTagNames = &tn
}

func (s *Share) SetAnyTagNames(tags TagList) {
	tn := tags.Names()
	s.AnyTagNames = &tn
}

type ShareList []*Share

func (sl ShareList) Keys() []string {
	keys := make([]string, len(sl))
	for i := range sl {
		keys[i] = sl[i].ID
	}
	return keys
}

func (sl ShareList) SetMustTagNames(tagMap map[string]TagList) {
	for i := range sl {
		k := sl[i].ID
		sl[i].SetMustTagNames(tagMap[k])
	}
}

func (sl ShareList) SetAnyTagNames(tagMap map[string]TagList) {
	for i := range sl {
		k := sl[i].ID
		tags := append(make([]*Tag, 0), tagMap[k]...)
		sl[i].SetAnyTagNames(tags)
	}
}
