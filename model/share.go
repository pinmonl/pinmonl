package model

import (
	"regexp"

	"github.com/pinmonl/pinmonl/model/field"
)

// Share contains information of a sharing list.
type Share struct {
	ID          string     `json:"id"          db:"id"`
	UserID      string     `json:"userId"      db:"user_id"`
	Name        string     `json:"name"        db:"name"`
	Description string     `json:"description" db:"description"`
	Readme      string     `json:"readme"      db:"readme"`
	ImageID     string     `json:"imageId"     db:"image_id"`
	CreatedAt   field.Time `json:"createdAt"   db:"created_at"`
	UpdatedAt   field.Time `json:"updatedAt"   db:"updated_at"`
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

// ShareTag defines the connection between share and tag.
type ShareTag struct {
	ShareID  string `json:"shareId"  db:"share_id"`
	TagID    string `json:"tagId"    db:"tag_id"`
	Kind     string `json:"kind"     db:"kind"`
	ParentID string `json:"parentId" db:"parent_id"`
	Sort     int64  `json:"sort"     db:"sort"`
}

// ShareTagKind categories the group of ShareTag.
type ShareTagKind string

const (
	// MustTag defines the key of must-exist kind.
	MustTag ShareTagKind = "must"
	// AnyTag defines the key of any kind.
	AnyTag ShareTagKind = "any"
)
