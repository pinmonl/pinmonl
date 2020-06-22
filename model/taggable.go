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

func (tl TaggableList) Tags() TagList {
	tags := make([]*Tag, len(tl))
	for i := range tl {
		tags[i] = tl[i].Tag
	}
	return tags
}

func (tl TaggableList) TargetKeys() []string {
	keys := make([]string, len(tl))
	for i := range tl {
		keys[i] = tl[i].TargetID
	}
	return keys
}

func (tl TaggableList) TagsByTarget() map[string]TagList {
	out := make(map[string]TagList)
	for _, tg := range tl {
		k := tg.TargetID
		out[k] = append(out[k], tg.Tag)
	}
	return out
}
