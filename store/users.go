package store

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type Users struct {
	*Store
}

type UserOpts struct {
	ListOpts
	Login  string
	Role   field.NullValue
	Status field.NullValue
}

func NewUsers(s *Store) *Users {
	return &Users{s}
}

func (u *Users) table() string {
	return "users"
}

func (u *Users) List(ctx context.Context, opts *UserOpts) ([]*model.User, error) {
	if opts == nil {
		opts = &UserOpts{}
	}

	qb := u.RunnableBuilder(ctx).
		Select(u.columns()...).From(u.table())
	qb = u.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.User
	for rows.Next() {
		user, err := u.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, user)
	}
	return list, nil
}

func (u *Users) Count(ctx context.Context, opts *UserOpts) (int64, error) {
	if opts == nil {
		opts = &UserOpts{}
	}

	qb := u.RunnableBuilder(ctx).
		Select("count(*)").From(u.table())
	qb = u.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (u *Users) Find(ctx context.Context, id string) (*model.User, error) {
	qb := u.RunnableBuilder(ctx).
		Select(u.columns()...).From(u.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	user, err := u.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Users) FindLogin(ctx context.Context, login string) (*model.User, error) {
	qb := u.RunnableBuilder(ctx).
		Select(u.columns()...).From(u.table()).
		Where("login = ?", login).
		Where("status = ?", model.ActiveUser)
	row := qb.QueryRow()
	user, err := u.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Users) bindOpts(b squirrel.SelectBuilder, opts *UserOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if opts.Login != "" {
		b = b.Where("login = ?", opts.Login)
	}

	if opts.Status.Valid {
		if s, ok := opts.Status.Value().(model.UserStatus); ok {
			b = b.Where("status = ?", s)
		}
	}

	if opts.Role.Valid {
		if r, ok := opts.Role.Value().(model.UserRole); ok {
			b = b.Where("role = ?", r)
		}
	}

	return b
}

func (u *Users) columns() []string {
	return []string{
		"id",
		"login",
		"password",
		"name",
		"image_id",
		"hash",
		"role",
		"status",
		"last_seen",
		"created_at",
		"updated_at",
	}
}

func (u *Users) scan(row database.RowScanner) (*model.User, error) {
	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.Name,
		&user.ImageID,
		&user.Hash,
		&user.Role,
		&user.Status,
		&user.LastSeen,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *Users) Create(ctx context.Context, user *model.User) error {
	user2 := *user
	user2.ID = newID()
	user2.CreatedAt = timestamp()

	qb := u.RunnableBuilder(ctx).
		Insert(u.table()).
		Columns(
			"id",
			"login",
			"password",
			"name",
			"image_id",
			"hash",
			"role",
			"status",
			"last_seen",
			"created_at").
		Values(
			user2.ID,
			user2.Login,
			user2.Password,
			user2.Name,
			user2.ImageID,
			user2.Hash,
			user2.Role,
			user2.Status,
			user2.LastSeen,
			user2.CreatedAt)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*user = user2
	return nil
}

func (u *Users) Update(ctx context.Context, user *model.User) error {
	user2 := *user
	user2.UpdatedAt = timestamp()

	qb := u.RunnableBuilder(ctx).
		Update(u.table()).
		Set("login", user2.Login).
		Set("password", user2.Password).
		Set("name", user2.Name).
		Set("image_id", user2.ImageID).
		Set("hash", user2.Hash).
		Set("role", user2.Role).
		Set("status", user2.Status).
		Set("last_seen", user2.LastSeen).
		Set("updated_at", user2.UpdatedAt).
		Where("id = ?", user2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*user = user2
	return nil
}

func (u *Users) Delete(ctx context.Context, id string) (int64, error) {
	qb := u.RunnableBuilder(ctx).
		Delete(u.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
