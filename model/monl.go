package model

import "github.com/pinmonl/pinmonl/model/field"

type Monl struct {
	ID        string     `json:"id"`
	URL       string     `json:"url"`
	FetchedAt field.Time `json:"fetchedAt"`
	CreatedAt field.Time `json:"createdAt"`
	UpdatedAt field.Time `json:"updatedAt"`
}

func (m Monl) MorphKey() string  { return m.ID }
func (m Monl) MorphName() string { return "monl" }

type MonlList []*Monl

func (ml MonlList) Keys() []string {
	keys := make([]string, len(ml))
	for i := range ml {
		keys[i] = ml[i].ID
	}
	return keys
}
