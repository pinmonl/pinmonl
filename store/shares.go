package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

type Shares struct {
	*Store
}

type ShareOpts struct {
	ListOpts
	UserID  string
	UserIDs []string
	Slug    string
}

func NewShares(s *Store) *Shares {
	return &Shares{s}
}

func (s *Shares) table() string {
	return "shares"
}

func (s *Shares) List(ctx context.Context, opts *ShareOpts) ([]*model.Share, error) {
	if opts == nil {
		opts = &ShareOpts{}
	}

	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table())
	qb = s.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Share
	for rows.Next() {
		share, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, share)
	}
	return list, nil
}

func (s *Shares) Count(ctx context.Context, opts *ShareOpts) (int64, error) {
	if opts == nil {
		opts = &ShareOpts{}
	}

	qb := s.RunnableBuilder(ctx).
		Select("count(*)").From(s.table())
	qb = s.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Shares) Find(ctx context.Context, id string) (*model.Share, error) {
	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	share, err := s.scan(row)
	if err != nil {
		return nil, err
	}
	return share, nil
}

func (s *Shares) columns() []string {
	return []string{
		"id",
		"user_id",
		"slug",
		"name",
		"description",
		"image_id",
		"created_at",
		"updated_at",
	}
}

func (s *Shares) bindOpts(b squirrel.SelectBuilder, opts *ShareOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	return b
}

func (s *Shares) scan(row database.RowScanner) (*model.Share, error) {
	var share model.Share
	err := row.Scan(
		&share.ID,
		&share.UserID,
		&share.Slug,
		&share.Name,
		&share.Description,
		&share.ImageID,
		&share.CreatedAt,
		&share.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (s *Shares) Create(ctx context.Context, share *model.Share) error {
	share2 := *share
	share2.ID = newID()
	share2.CreatedAt = timestamp()

	qb := s.RunnableBuilder(ctx).
		Insert(s.table()).
		Columns(
			"id",
			"user_id",
			"slug",
			"name",
			"description",
			"image_id",
			"created_at").
		Values(
			share2.ID,
			share2.UserID,
			share2.Slug,
			share2.Name,
			share2.Description,
			share2.ImageID,
			share2.CreatedAt)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*share = share2
	return nil
}

func (s *Shares) Update(ctx context.Context, share *model.Share) error {
	share2 := *share
	share2.UpdatedAt = timestamp()

	qb := s.RunnableBuilder(ctx).
		Update(s.table()).
		Set("user_id", share2.UserID).
		Set("slug", share2.Slug).
		Set("name", share2.Name).
		Set("description", share2.Description).
		Set("image_id", share2.ImageID).
		Set("updated_at", share2.UpdatedAt).
		Where("id = ?", share2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*share = share2
	return nil
}

func (s *Shares) Delete(ctx context.Context, id string) (int64, error) {
	qb := s.RunnableBuilder(ctx).
		Delete(s.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
