package model

type Sharetag struct {
	ID          string       `json:"id"`
	ShareID     string       `json:"shareId"`
	TagID       string       `json:"tagId"`
	Kind        SharetagKind `json:"kind"`
	ParentID    string       `json:"parentId"`
	Level       int          `json:"level"`
	Status      ShareStatus  `json:"status"`
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
	var keys []string
	for _, s := range sl {
		keys = append(keys, s.ID)
	}
	return keys
}
