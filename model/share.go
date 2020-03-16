package model

import (
	"regexp"

	"github.com/pinmonl/pinmonl/model/field"
)

// Share contains information of a sharing list.
type Share struct {
	ID          string     `json:"id"          db:"share_id"`
	UserID      string     `json:"userId"      db:"share_user_id"`
	Name        string     `json:"name"        db:"share_name"`
	Description string     `json:"description" db:"share_description"`
	Readme      string     `json:"readme"      db:"share_readme"`
	ImageID     string     `json:"imageId"     db:"share_image_id"`
	CreatedAt   field.Time `json:"createdAt"   db:"share_created_at"`
	UpdatedAt   field.Time `json:"updatedAt"   db:"share_updated_at"`
}

// ShareNamePattern is the pattern of share name.
var ShareNamePattern = regexp.MustCompile("^[a-zA-Z0-9-_]+$")

// ShareList is slice of Share.
type ShareList []Share

// Keys gets ID slice from Shares.
func (sl ShareList) Keys() []string {
	out := make([]string, len(sl))
	for i, s := range sl {
		out[i] = s.ID
	}
	return out
}
