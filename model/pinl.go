package model

import "github.com/pinmonl/pinmonl/model/field"

type Pinl struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	MonlID      string     `json:"monlId"`
	URL         string     `json:"url"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ImageID     string     `json:"imageId"`
	Status      Status     `json:"status"`
	CreatedAt   field.Time `json:"createdAt"`
	UpdatedAt   field.Time `json:"updatedAt"`

	Tags     *TagList  `json:"-"`
	TagNames *[]string `json:"tags,omitempty"`
}

func (p Pinl) MorphKey() string  { return p.ID }
func (p Pinl) MorphName() string { return "pinl" }

func (p *Pinl) SetTagNames(tags TagList) {
	tn := tags.Names()
	p.TagNames = &tn
}

type PinlList []*Pinl

func (pl PinlList) Keys() []string {
	keys := make([]string, len(pl))
	for i := range pl {
		keys[i] = pl[i].ID
	}
	return keys
}

func (pl PinlList) Morphables() MorphableList {
	list := make([]Morphable, len(pl))
	for i := range pl {
		list[i] = pl[i]
	}
	return list
}

func (pl PinlList) SetTagNames(tagMap map[string]TagList) {
	for i := range pl {
		k := pl[i].ID
		pl[i].SetTagNames(tagMap[k])
	}
}
