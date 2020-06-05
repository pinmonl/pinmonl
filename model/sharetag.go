package model

type Sharetag struct {
	ID          string       `json:"id"`
	ShareID     string       `json:"shareId"`
	TagID       string       `json:"tagId"`
	Kind        SharetagKind `json:"kind"`
	ParentID    string       `json:"parentId"`
	Level       int          `json:"level"`
	HasChildren bool         `json:"hasChildren"`

	Share    *Share        `json:"share,omitempty"`
	Tag      *Tag          `json:"tag,omitempty"`
	Children *SharetagList `json:"children,omitempty"`
}

type SharetagKind int

const (
	SharetagMust SharetagKind = iota
	SharetagAny
)

type SharetagList []*Sharetag

func (sl SharetagList) Keys() []string {
	var keys []string
	for _, s := range sl {
		keys = append(keys, s.ID)
	}
	return keys
}
