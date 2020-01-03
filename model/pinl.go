package model

import (
	"github.com/pinmonl/pinmonl/model/field"
)

// Pinl stores bookmark-liked record.
type Pinl struct {
	ID          string     `json:"id"          db:"id"`
	UserID      string     `json:"userId"      db:"user_id"`
	URL         string     `json:"url"         db:"url"`
	Title       string     `json:"title"       db:"title"`
	Description string     `json:"description" db:"description"`
	Readme      string     `json:"readme"      db:"readme"`
	ImageID     string     `json:"imageId"     db:"image_id"`
	CreatedAt   field.Time `json:"createdAt"   db:"created_at"`
	UpdatedAt   field.Time `json:"updatedAt"   db:"updated_at"`
}

// MorphKey returns the key of Pinl.
func (p Pinl) MorphKey() string { return p.ID }

// MorphName returns the name of Pinl.
func (p Pinl) MorphName() string { return "pinl" }

// Pinmon defines the connection between Pinl and Monl.
type Pinmon struct {
	PinlID string `json:"pinlId" db:"pinl_id"`
	MonlID string `json:"monlId" db:"monl_id"`
	UserID string `json:"userId" db:"user_id"`
	Sort   int64  `json:"sort"   db:"sort"`
}
