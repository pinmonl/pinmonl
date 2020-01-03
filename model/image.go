package model

import "github.com/pinmonl/pinmonl/model/field"

// Image contains the file content and information.
type Image struct {
	ID          string     `json:"id"          db:"id"`
	TargetID    string     `json:"-"           db:"target_id"`
	TargetName  string     `json:"-"           db:"target_name"`
	Kind        string     `json:"kind"        db:"kind"`
	Sort        int64      `json:"sort"        db:"sort"`
	Filename    string     `json:"filename"    db:"filename"`
	Content     []byte     `json:"content"     db:"content"`
	Description string     `json:"description" db:"description"`
	Size        int64      `json:"size"        db:"size"`
	CreatedAt   field.Time `json:"createdAt"   db:"created_at"`
	UpdatedAt   field.Time `json:"updatedAt"   db:"updated_at"`
}
