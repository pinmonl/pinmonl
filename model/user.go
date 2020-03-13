package model

import "github.com/pinmonl/pinmonl/model/field"

// User defines the user login account.
type User struct {
	ID        string     `json:"id"        db:"id"`
	Login     string     `json:"login"     db:"login"`
	Password  string     `json:"-"         db:"password"`
	Name      string     `json:"name"      db:"name"`
	ImageID   string     `json:"imageId"   db:"image_id"`
	Role      UserRole   `json:"role"      db:"role"`
	Hash      string     `json:"hash"      db:"hash"`
	CreatedAt field.Time `json:"createdAt" db:"created_at"`
	UpdatedAt field.Time `json:"updatedAt" db:"updated_at"`
	LastLog   field.Time `json:"lastLog"   db:"last_log"`
}

// UserRole defines the nature of user account.
type UserRole int

const (
	// UserRoleNormal is the default value of user role.
	UserRoleNormal UserRole = iota
)
