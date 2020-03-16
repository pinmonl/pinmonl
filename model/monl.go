package model

import (
	"github.com/pinmonl/pinmonl/model/field"
)

// Monl stores the package information.
type Monl struct {
	ID          string     `json:"id"          db:"monl_id"`
	URL         string     `json:"url"         db:"monl_url"`
	Title       string     `json:"title"       db:"monl_title"`
	Description string     `json:"description" db:"monl_description"`
	Readme      string     `json:"readme"      db:"monl_readme"`
	ImageID     string     `json:"imageId"     db:"monl_image_id"`
	CreatedAt   field.Time `json:"createdAt"   db:"monl_created_at"`
	UpdatedAt   field.Time `json:"updatedAt"   db:"monl_updated_at"`
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
