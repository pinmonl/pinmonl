package model

import "github.com/pinmonl/pinmonl/model/field"

// User defines the user login account.
type User struct {
	ID        string     `json:"id"        db:"user_id"`
	Login     string     `json:"login"     db:"user_login"`
	Password  string     `json:"-"         db:"user_password"`
	Name      string     `json:"name"      db:"user_name"`
	ImageID   string     `json:"imageId"   db:"user_image_id"`
	Role      UserRole   `json:"role"      db:"user_role"`
	Hash      string     `json:"-"         db:"user_hash"`
	CreatedAt field.Time `json:"createdAt" db:"user_created_at"`
	UpdatedAt field.Time `json:"updatedAt" db:"user_updated_at"`
	LastLog   field.Time `json:"lastLog"   db:"user_last_log"`
}

// UserRole defines the nature of user account.
type UserRole int

const (
	// UserRoleNormal is the default value of user role.
	UserRoleNormal UserRole = iota
)
