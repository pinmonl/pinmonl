package model

import "github.com/pinmonl/pinmonl/model/field"

type Image struct {
	ID          string     `json:"id"`
	TargetID    string     `json:"targetId"`
	TargetName  string     `json:"targetName"`
	Content     []byte     `json:"content"`
	Description string     `json:"description"`
	Size        int        `json:"size"`
	ContentType string     `json:"contentType"`
	CreatedAt   field.Time `json:"createdAt"`
	UpdatedAt   field.Time `json:"updatedAt"`

	Target Morphable `json:"-"`
}

func (i Image) MorphKey() string  { return i.ID }
func (i Image) MorphName() string { return "image" }

type ImageList []*Image

func (il ImageList) Keys() []string {
	keys := make([]string, len(il))
	for i := range il {
		keys[i] = il[i].ID
	}
	return keys
}
