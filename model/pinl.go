package model

import (
	"github.com/pinmonl/pinmonl/model/field"
)

// Pinl stores bookmark-liked record.
type Pinl struct {
	ID          string     `json:"id"          db:"pinl_id"`
	UserID      string     `json:"userId"      db:"pinl_user_id"`
	URL         string     `json:"url"         db:"pinl_url"`
	Title       string     `json:"title"       db:"pinl_title"`
	Description string     `json:"description" db:"pinl_description"`
	Readme      string     `json:"readme"      db:"pinl_readme"`
	ImageID     string     `json:"imageId"     db:"pinl_image_id"`
	CreatedAt   field.Time `json:"createdAt"   db:"pinl_created_at"`
	UpdatedAt   field.Time `json:"updatedAt"   db:"pinl_updated_at"`
}

// MorphKey returns the key of Pinl.
func (p Pinl) MorphKey() string { return p.ID }

// MorphName returns the name of Pinl.
func (p Pinl) MorphName() string { return "pinl" }
