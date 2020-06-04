package model

type Sharetag struct {
	ID          string       `json:"id"`
	ShareID     string       `json:"shareId"`
	TagID       string       `json:"tagId"`
	Kind        SharetagKind `json:"kind"`
	ParentID    string       `json:"parentId"`
	Level       int          `json:"level"`
	HasChildren bool         `json:"hasChildren"`

	Share *Share `json:"share,omitempty"`
	Tag   *Tag   `json:"tag,omitempty"`
}

type SharetagKind int

const (
	SharetagKindNotSet SharetagKind = iota
	SharetagMust
	SharetagAny
)
