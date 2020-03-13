package model

import (
	"regexp"

	"github.com/pinmonl/pinmonl/model/field"
)

// Tag defines the structure of tag.
type Tag struct {
	ID        string     `json:"id"        db:"id"`
	Name      string     `json:"name"      db:"name"`
	UserID    string     `json:"userId"    db:"user_id"`
	ParentID  string     `json:"parentId"  db:"parent_id"`
	Sort      int64      `json:"sort"      db:"sort"`
	Level     int64      `json:"level"     db:"level"`
	CreatedAt field.Time `json:"createdAt" db:"created_at"`
	UpdatedAt field.Time `json:"updatedAt" db:"updated_at"`

	TargetID   string `json:"-" db:"target_id"`
	TargetName string `json:"-" db:"target_name"`
}

// TagNamePattern is the pattern of tag name.
var TagNamePattern = regexp.MustCompile("^[a-zA-Z0-9-_\\.: ()\\[\\]]+$")

// TagList is slice of Tag.
type TagList []Tag

// FindMorphable filters Tags by Morphable.
func (list TagList) FindMorphable(target Morphable) TagList {
	out := make([]Tag, 0)
	for _, t := range list {
		if t.TargetID == target.MorphKey() && t.TargetName == target.MorphName() {
			out = append(out, t)
		}
	}
	return out
}

// PluckName returns slice of tag name.
func (list TagList) PluckName() []string {
	out := make([]string, len(list))
	for i, t := range list {
		out[i] = t.Name
	}
	return out
}

// Taggable defines the connection between tag and morphable record.
type Taggable struct {
	TagID      string `json:"tagId"    db:"tag_id"`
	TargetID   string `json:"targetId" db:"target_id"`
	TargetName string `json:"-"        db:"target_name"`
	Sort       int64  `json:"sort"     db:"sort"`
}
