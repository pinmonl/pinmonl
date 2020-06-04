package model

type Taggable struct {
	ID           string `json:"id"`
	TagID        string `json:"tagId"`
	TaggableID   string `json:"taggableId"`
	TaggableName string `json:"taggableName"`

	Tag      *Tag      `json:"tag,omitempty"`
	Taggable Morphable `json:"-"`
}
