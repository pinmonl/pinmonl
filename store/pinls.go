package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type Pinls struct {
	*Store
}

type PinlOpts struct {
	ListOpts
	IDs     []string
	UserID  string
	UserIDs []string
	MonlIDs []string
	Query   string
	Status  field.NullValue
	URL     string

	TagIDs          []string
	TagNames        []string
	TagNamePatterns []string
	NoTag           field.NullBool

	Orders []PinlOrder
}

type PinlOrder int

const (
	PinlOrderByLatest PinlOrder = iota
)

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
	list := make([]*model.Pinl, 0)
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

	o2 := *opts
	o2.Orders = nil

	qb := p.RunnableBuilder(ctx).
		Select("count(*)").From(p.table())
	qb = p.bindOpts(qb, &o2)
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

func (p Pinls) bindOpts(b squirrel.SelectBuilder, opts *PinlOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.IDs) > 0 {
		b = b.Where(squirrel.Eq{"id": opts.IDs})
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

	if opts.Query != "" {
		b = b.Where(squirrel.Or{
			squirrel.Expr("title like ?", "%"+opts.Query+"%"),
			squirrel.Expr("description like ?", "%"+opts.Query+"%"),
			squirrel.Expr("url like ?", "%"+opts.Query+"%"),
		})
	}

	if opts.Status.Valid {
		if s, ok := opts.Status.Value().(model.Status); ok {
			b = b.Where("status = ?", s)
		}
	}

	if opts.URL != "" {
		b = b.Where("url = ?", opts.URL)
	}

	if len(opts.TagIDs) > 0 {
		sq := p.Builder().Select("1").
			From(Taggables{}.table()).
			Where("target_id = "+p.table()+".id").
			Where("target_name = ?", model.Pinl{}.MorphName()).
			Where(squirrel.Eq{"tag_id": opts.TagIDs}).
			GroupBy("target_id").
			Having("COUNT( DISTINCT tag_id ) >= ?", len(opts.TagIDs)).
			Prefix("EXISTS (").
			Suffix(")")
		b = b.Where(sq)
	}

	if len(opts.TagNames) > 0 {
		sq := p.Builder().Select("1").
			From(Taggables{}.table()).
			Join(fmt.Sprintf("%s ON %[1]s.id = %s.tag_id", Tags{}.table(), Taggables{}.table())).
			Where("target_id = "+p.table()+".id").
			Where("target_name = ?", model.Pinl{}.MorphName()).
			Where(squirrel.Eq{"name": opts.TagNames}).
			GroupBy("target_id").
			Having("COUNT( DISTINCT tag_id ) >= ?", len(opts.TagNames)).
			Prefix("EXISTS (").
			Suffix(")")
		b = b.Where(sq)
	}

	if len(opts.TagNamePatterns) > 0 {
		sq := p.Builder().Select("1").
			From(Taggables{}.table()).
			Join(fmt.Sprintf("%s ON %[1]s.id = %s.tag_id", Tags{}.table(), Taggables{}.table())).
			Where("target_id = "+p.table()+".id").
			Where("target_name = ?", model.Pinl{}.MorphName()).
			GroupBy("target_id").
			Having("COUNT( DISTINCT tag_id ) >= ?", len(opts.TagNames)).
			Prefix("EXISTS (").
			Suffix(")")
		for _, tagName := range opts.TagNamePatterns {
			b = b.Where(sq.Where("name like ?", tagName))
		}
	}

	if opts.NoTag.Valid {
		sq := p.Builder().Select("1").
			From(Taggables{}.table()).
			Join(fmt.Sprintf("%s ON %[1]s.id = %s.tag_id", Tags{}.table(), Taggables{}.table())).
			Where("target_id = "+p.table()+".id").
			Where("target_name = ?", model.Pinl{}.MorphName()).
			GroupBy("target_id").
			Suffix(")")

		if opts.NoTag.Value() {
			sq = sq.Prefix("NOT EXISTS (")
		} else {
			sq = sq.Prefix("EXISTS (")
		}
		b = b.Where(sq)
	}

	for _, order := range opts.Orders {
		switch order {
		case PinlOrderByLatest:
			b = b.OrderBy("created_at DESC")
		}
	}

	return b
}

func (p Pinls) columns() []string {
	return []string{
		p.table() + ".id",
		p.table() + ".user_id",
		p.table() + ".monl_id",
		p.table() + ".url",
		p.table() + ".title",
		p.table() + ".description",
		p.table() + ".image_id",
		p.table() + ".status",
		p.table() + ".created_at",
		p.table() + ".updated_at",
	}
}

func (p Pinls) scanColumns(pinl *model.Pinl) []interface{} {
	return []interface{}{
		&pinl.ID,
		&pinl.UserID,
		&pinl.MonlID,
		&pinl.URL,
		&pinl.Title,
		&pinl.Description,
		&pinl.ImageID,
		&pinl.Status,
		&pinl.CreatedAt,
		&pinl.UpdatedAt,
	}
}

func (p Pinls) scan(row database.RowScanner) (*model.Pinl, error) {
	var pinl model.Pinl
	err := row.Scan(p.scanColumns(&pinl)...)
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
