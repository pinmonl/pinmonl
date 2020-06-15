package model

type Taggable struct {
	ID         string `json:"id"`
	TagID      string `json:"tagId"`
	TargetID   string `json:"targetId"`
	TargetName string `json:"targetName"`

	Tag  *Tag  `json:"tag,omitempty"`
	Pinl *Pinl `json:"pinl,omitempty"`
}

type TaggableList []*Taggable

func (tl TaggableList) Tags() []*Tag {
	tags := make([]*Tag, len(tl))
	for i := range tl {
		tags[i] = tl[i].Tag
	}
	return tags
}

func (tl TaggableList) TargetKeys() []string {
	keys := make([]string, 0)
	for _, t := range tl {
		keys = append(keys, t.TargetID)
	}
	return keys
}
