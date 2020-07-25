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

func (sl SharetagList) GetKind(kind SharetagKind) SharetagList {
	list := make([]*Sharetag, 0)
	for i := range sl {
		if sl[i].Kind == kind {
			list = append(list, sl[i])
		}
	}
	return list
}

func (sl SharetagList) TagsByShare() map[string]TagList {
	out := make(map[string]TagList)
	for i := range sl {
		k := sl[i].ShareID
		out[k] = append(out[k], sl[i].Tag)
	}
	return out
}

func (sl SharetagList) ByTagID() map[string]*Sharetag {
	out := make(map[string]*Sharetag)
	for i := range sl {
		k := sl[i].TagID
		out[k] = sl[i]
	}
	return out
}
