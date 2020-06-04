package model

import (
	"github.com/pinmonl/pinmonl/model/field"
)

type User struct {
	ID        string     `json:"id"`
	Login     string     `json:"login"`
	Password  string     `json:"-"`
	Name      string     `json:"name"`
	ImageID   string     `json:"imageId"`
	Hash      string     `json:"-"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
	LastSeen  field.Time `json:"lastSeen"`
	CreatedAt field.Time `json:"createdAt"`
	UpdatedAt field.Time `json:"updatedAt"`
}

func (u User) MorphKey() string  { return u.ID }
func (u User) MorphName() string { return "user" }

type UserRole int

const (
	UserRoleNotSet UserRole = iota
	NormalUser
)

type UserStatus int

const (
	UserStatusNotSet UserStatus = iota
	ActiveUser
	ExpiredUser
)
