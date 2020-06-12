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
}

func (p Pinl) MorphKey() string  { return p.ID }
func (p Pinl) MorphName() string { return "pinl" }

type PinlList []*Pinl

func (pl PinlList) Keys() []string {
	var keys []string
	for _, p := range pl {
		keys = append(keys, p.ID)
	}
	return keys
}
