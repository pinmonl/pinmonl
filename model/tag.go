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

	Children *TagList `json:"children,omitempty"`
}

func (t Tag) MorphKey() string  { return t.ID }
func (t Tag) MorphName() string { return "tag" }

type TagList []*Tag

func (tl TagList) Keys() []string {
	keys := make([]string, len(tl))
	for i := range tl {
		keys[i] = tl[i].ID
	}
	return keys
}

func (tl TagList) Names() []string {
	names := make([]string, len(tl))
	for i := range tl {
		names[i] = tl[i].Name
	}
	return names
}
