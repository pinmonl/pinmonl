package model

type Sharetag struct {
	ID          string       `json:"id"`
	ShareID     string       `json:"shareId"`
	TagID       string       `json:"tagId"`
	Kind        SharetagKind `json:"kind"`
	ParentID    string       `json:"parentId"`
	Level       int          `json:"level"`
	Status      Status       `json:"status"`
	HasChildren bool         `json:"hasChildren"`

	Share    *Share        `json:"share,omitempty"`
	Tag      *Tag          `json:"tag,omitempty"`
	Children *SharetagList `json:"children,omitempty"`
}

func (s Sharetag) ViewTag() *Tag {
	if s.Tag == nil {
		return nil
	}
	tag := *s.Tag
	tag.ParentID = s.ParentID
	tag.Level = s.Level
	tag.HasChildren = s.HasChildren
	return &tag
}

type SharetagKind int

const (
	SharetagMust SharetagKind = iota
	SharetagAny
)

func IsValidSharetagKind(k SharetagKind) bool {
	checks := map[SharetagKind]int{
		SharetagMust: 0,
		SharetagAny:  0,
	}
	_, ok := checks[k]
	return ok
}

type SharetagList []*Sharetag

func (sl SharetagList) Keys() []string {
	keys := make([]string, len(sl))
	for i := range sl {
		keys[i] = sl[i].ID
	}
	return keys
}

func (sl SharetagList) Tags() TagList {
	tags := make([]*Tag, len(sl))
	for i := range sl {
		tags[i] = sl[i].Tag
	}
	return tags
}

func (sl SharetagList) ViewTags() TagList {
	tags := make([]*Tag, len(sl))
	for i := range sl {
		tags[i] = sl[i].ViewTag()
	}
	return tags
}
