package model

// Sharetag defines the connection between share and tag.
type Sharetag struct {
	*Share
	*Tag
	ShareID  string       `json:"shareId"  db:"sharetag_share_id"`
	TagID    string       `json:"tagId"    db:"sharetag_tag_id"`
	Kind     SharetagKind `json:"kind"     db:"sharetag_kind"`
	ParentID string       `json:"parentId" db:"sharetag_parent_id"`
	Sort     int64        `json:"sort"     db:"sharetag_sort"`
	Level    int64        `json:"level"    db:"sharetag_level"`
}

// SharetagKind categories the group of Sharetag.
type SharetagKind int

const (
	// SharetagKindEmpty indicates the zero value of SharetagKind.
	SharetagKindEmpty SharetagKind = iota
	// SharetagKindMust defines the key of must-exist kind.
	SharetagKindMust
	// SharetagKindAny defines the key of any kind.
	SharetagKindAny
)
