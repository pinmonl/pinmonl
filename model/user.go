package model

import "github.com/pinmonl/pinmonl/model/field"

// User defines the user login account.
type User struct {
	ID        string     `json:"id"        db:"id"`
	Login     string     `json:"login"     db:"login"`
	Password  string     `json:"-"         db:"password"`
	Name      string     `json:"name"      db:"name"`
	Email     string     `json:"email"     db:"email"`
	ImageID   string     `json:"imageId"   db:"image_id"`
	CreatedAt field.Time `json:"createdAt" db:"created_at"`
	UpdatedAt field.Time `json:"updatedAt" db:"updated_at"`
}
