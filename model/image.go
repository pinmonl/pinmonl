package model

import "github.com/pinmonl/pinmonl/model/field"

// Image contains the file content and information.
type Image struct {
	ID          string     `json:"id"          db:"image_id"`
	TargetID    string     `json:"-"           db:"image_target_id"`
	TargetName  string     `json:"-"           db:"image_target_name"`
	ContentType string     `json:"contentType" db:"image_content_type"`
	Sort        int64      `json:"sort"        db:"image_sort"`
	Filename    string     `json:"filename"    db:"image_filename"`
	Content     []byte     `json:"content"     db:"image_content"`
	Description string     `json:"description" db:"image_description"`
	Size        int64      `json:"size"        db:"image_size"`
	CreatedAt   field.Time `json:"createdAt"   db:"image_created_at"`
	UpdatedAt   field.Time `json:"updatedAt"   db:"image_updated_at"`
}
