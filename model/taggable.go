package model

// Taggable defines the connection between tag and morphable record.
type Taggable struct {
	*Tag
	*Pinl
	TagID      string `json:"tagId"      db:"taggable_tag_id"`
	TargetID   string `json:"targetId"   db:"taggable_target_id"`
	TargetName string `json:"targetName" db:"taggable_target_name"`
	Sort       int64  `json:"sort"       db:"taggable_sort"`
}
