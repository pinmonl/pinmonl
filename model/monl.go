package model

import "github.com/pinmonl/pinmonl/model/field"

type Monl struct {
	ID        string     `json:"id"`
	URL       string     `json:"url"`
	CreatedAt field.Time `json:"createdAt"`
	UpdatedAt field.Time `json:"updatedAt"`
}

func (m Monl) MorphKey() string  { return m.ID }
func (m Monl) MorphName() string { return "monl" }

type MonlList []*Monl

func (ml MonlList) Keys() []string {
	keys := make([]string, 0)
	for _, m := range ml {
		keys = append(keys, m.ID)
	}
	return keys
}
