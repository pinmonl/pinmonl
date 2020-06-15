package store

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

type Pinls struct {
	*Store
}

type PinlOpts struct {
	ListOpts
	UserID  string
	UserIDs []string
	MonlIDs []string
}

func NewPinls(s *Store) *Pinls {
	return &Pinls{s}
}

func (p Pinls) table() string {
	return "pinls"
}

func (p *Pinls) List(ctx context.Context, opts *PinlOpts) (model.PinlList, error) {
	if opts == nil {
		opts = &PinlOpts{}
	}

	qb := p.RunnableBuilder(ctx).
		Select(p.columns()...).From(p.table())
	qb = p.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Pinl
	for rows.Next() {
		pinl, err := p.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, pinl)
	}
	return list, nil
}

func (p *Pinls) Count(ctx context.Context, opts *PinlOpts) (int64, error) {
	if opts == nil {
		opts = &PinlOpts{}
	}

	qb := p.RunnableBuilder(ctx).
		Select("count(*)").From(p.table())
	qb = p.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *Pinls) Find(ctx context.Context, id string) (*model.Pinl, error) {
	qb := p.RunnableBuilder(ctx).
		Select(p.columns()...).From(p.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	pinl, err := p.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return pinl, nil
}

func (p Pinls) columns() []string {
	return []string{
		"id",
		"user_id",
		"monl_id",
		"url",
		"title",
		"description",
		"image_id",
		"status",
		"created_at",
		"updated_at",
	}
}

func (p Pinls) bindOpts(b squirrel.SelectBuilder, opts *PinlOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if opts.UserID != "" {
		opts.UserIDs = append(opts.UserIDs, opts.UserID)
	}
	if len(opts.UserIDs) > 0 {
		b = b.Where(squirrel.Eq{"user_id": opts.UserIDs})
	}

	if len(opts.MonlIDs) > 0 {
		b = b.Where(squirrel.Eq{"monl_id": opts.MonlIDs})
	}

	return b
}

func (p Pinls) scan(row database.RowScanner) (*model.Pinl, error) {
	var pinl model.Pinl
	err := row.Scan(
		&pinl.ID,
		&pinl.UserID,
		&pinl.MonlID,
		&pinl.URL,
		&pinl.Title,
		&pinl.Description,
		&pinl.ImageID,
		&pinl.Status,
		&pinl.CreatedAt,
		&pinl.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &pinl, nil
}

func (p *Pinls) Create(ctx context.Context, pinl *model.Pinl) error {
	pinl2 := *pinl
	pinl2.ID = newID()
	pinl2.CreatedAt = timestamp()
	pinl2.UpdatedAt = timestamp()

	qb := p.RunnableBuilder(ctx).
		Insert(p.table()).
		Columns(
			"id",
			"user_id",
			"monl_id",
			"url",
			"title",
			"description",
			"image_id",
			"status",
			"created_at",
			"updated_at").
		Values(
			pinl2.ID,
			pinl2.UserID,
			pinl2.MonlID,
			pinl2.URL,
			pinl2.Title,
			pinl2.Description,
			pinl2.ImageID,
			pinl2.Status,
			pinl2.CreatedAt,
			pinl2.UpdatedAt)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*pinl = pinl2
	return nil
}

func (p *Pinls) Update(ctx context.Context, pinl *model.Pinl) error {
	pinl2 := *pinl
	pinl2.UpdatedAt = timestamp()

	qb := p.RunnableBuilder(ctx).
		Update(p.table()).
		Set("user_id", pinl2.UserID).
		Set("monl_id", pinl2.MonlID).
		Set("url", pinl2.URL).
		Set("title", pinl2.Title).
		Set("description", pinl2.Description).
		Set("image_id", pinl2.ImageID).
		Set("status", pinl2.Status).
		Set("updated_at", pinl2.UpdatedAt).
		Where("id = ?", pinl2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*pinl = pinl2
	return nil
}

func (p *Pinls) Delete(ctx context.Context, id string) (int64, error) {
	qb := p.RunnableBuilder(ctx).
		Delete(p.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
