package model

import (
	"github.com/pinmonl/pinmonl/model/field"
)

// Monl stores the package information.
type Monl struct {
	ID          string     `json:"id"          db:"id"`
	URL         string     `json:"url"         db:"url"`
	Title       string     `json:"title"       db:"title"`
	Description string     `json:"description" db:"description"`
	Readme      string     `json:"readme"      db:"readme"`
	ImageID     string     `json:"imageId"     db:"image_id"`
	CreatedAt   field.Time `json:"createdAt"   db:"created_at"`
	UpdatedAt   field.Time `json:"updatedAt"   db:"updated_at"`
}

// MorphKey returns the key of Monl.
func (m Monl) MorphKey() string { return m.ID }

// MorphName returns the name of Monl.
func (m Monl) MorphName() string { return "monl" }

// MonlList is slice of Monl.
type MonlList []Monl

// Keys gets ID slice from Monls.
func (ml MonlList) Keys() []string {
	out := make([]string, len(ml))
	for i, m := range ml {
		out[i] = m.ID
	}
	return out
}
