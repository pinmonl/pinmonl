package model

import "github.com/pinmonl/pinmonl/model/field"

type Tag struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	UserID      string     `json:"userId"`
	ParentID    string     `json:"parentId"`
	Level       int        `json:"level"`
	Color       string     `json:"color"`
	BgColor     string     `json:"bgColor"`
	HasChildren bool       `json:"hasChildren"`
	CreatedAt   field.Time `json:"createdAt"`
	UpdatedAt   field.Time `json:"updatedAt"`

	Children *[]*Tag `json:"children,omitempty"`
}

func (t Tag) MorphKey() string  { return t.ID }
func (t Tag) MorphName() string { return "tag" }
