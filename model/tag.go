package model

import (
	"regexp"

	"github.com/pinmonl/pinmonl/model/field"
)

// Tag defines the structure of tag.
type Tag struct {
	ID        string     `json:"id"        db:"tag_id"`
	Name      string     `json:"name"      db:"tag_name"`
	UserID    string     `json:"userId"    db:"tag_user_id"`
	ParentID  string     `json:"parentId"  db:"tag_parent_id"`
	Sort      int64      `json:"sort"      db:"tag_sort"`
	Level     int64      `json:"level"     db:"tag_level"`
	Color     string     `json:"color"     db:"tag_color"`
	Bgcolor   string     `json:"bgcolor"   db:"tag_bgcolor"`
	CreatedAt field.Time `json:"createdAt" db:"tag_created_at"`
	UpdatedAt field.Time `json:"updatedAt" db:"tag_updated_at"`
}

// TagNamePattern is the pattern of tag name.
var TagNamePattern = regexp.MustCompile("^[a-zA-Z0-9-_\\.: ()\\[\\]]+$")

// TagList is slice of Tag.
type TagList []Tag

// PluckName returns slice of tag name.
func (list TagList) PluckName() []string {
	var out []string
	for _, t := range list {
		out = append(out, t.Name)
	}
	return out
}
