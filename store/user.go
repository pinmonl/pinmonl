package store

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
)

type Users struct {
	*Store
}

func NewUsers(s *Store) *Users {
	return &Users{s}
}

func (u *Users) Table() string {
	return "users"
}

func (u *Users) Columns() []string {
	return []string{"id", "login", "password", "name", "image_id",
		"hash", "role", "status", "created_at", "updated_at"}
}

func (u *Users) Create(ctx context.Context, data *model.User) error {
	return nil
}
