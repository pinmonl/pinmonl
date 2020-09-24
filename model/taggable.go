package model

type Taggable struct {
	ID          string `json:"id"`
	TagID       string `json:"tagId"`
	TargetID    string `json:"targetId"`
	TargetName  string `json:"targetName"`
	Value       string `json:"value"`
	ValueType   int    `json:"-"`
	ValuePrefix string `json:"-"`
	ValueSuffix string `json:"-"`

	Tag  *Tag  `json:"tag,omitempty"`
	Pinl *Pinl `json:"pinl,omitempty"`
}

func (t Taggable) Pivot() *TagPivot {
	return &TagPivot{
		Value:     t.Value,
		Prefix:    t.ValuePrefix,
		Suffix:    t.ValueSuffix,
		ValueType: t.ValueType,
	}
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

func (tl TaggableList) ByTarget() map[string]TaggableList {
	out := make(map[string]TaggableList)
	for _, tg := range tl {
		k := tg.TargetID
		out[k] = append(out[k], tg)
	}
	return out
}

type TagPivot struct {
	Prefix    string `json:"prefix"`
	Suffix    string `json:"suffix"`
	Value     string `json:"value"`
	ValueType int    `json:"valueType"`
}

type TagPivotList []*TagPivot
