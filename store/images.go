package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

type Images struct {
	*Store
}

type ImageOpts struct {
	ListOpts
	Targets model.MorphableList
}

func NewImages(s *Store) *Images {
	return &Images{s}
}

func (i *Images) table() string {
	return "images"
}

func (i *Images) List(ctx context.Context, opts *ImageOpts) ([]*model.Image, error) {
	if opts == nil {
		opts = &ImageOpts{}
	}

	qb := i.RunnableBuilder(ctx).
		Select(i.columns()...).From(i.table())
	qb = i.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Image
	for rows.Next() {
		image, err := i.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, image)
	}
	return list, nil
}

func (i *Images) Count(ctx context.Context, opts *ImageOpts) (int64, error) {
	if opts == nil {
		opts = &ImageOpts{}
	}

	qb := i.RunnableBuilder(ctx).
		Select("count(*)").From(i.table())
	qb = i.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (i *Images) Find(ctx context.Context, id string) (*model.Image, error) {
	qb := i.RunnableBuilder(ctx).
		Select(i.columns()...).From(i.table()).
		Where("id = ?", id).
		Limit(1)
	row := qb.QueryRow()
	image, err := i.scan(row)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (i *Images) columns() []string {
	return []string{
		"id",
		"target_id",
		"target_name",
		"content",
		"description",
		"size",
		"content_type",
		"created_at",
		"updated_at",
	}
}

func (i *Images) bindOpts(b squirrel.SelectBuilder, opts *ImageOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.Targets) > 0 && !opts.Targets.IsMixed() {
		b = b.Where("target_name = ?", opts.Targets.MorphName()).
			Where(squirrel.Eq{"target_id": opts.Targets.MorphKeys()})
	}

	return b
}

func (i *Images) scan(row database.RowScanner) (*model.Image, error) {
	var image model.Image
	err := row.Scan(
		&image.ID,
		&image.TargetID,
		&image.TargetName,
		&image.Content,
		&image.Description,
		&image.Size,
		&image.ContentType,
		&image.CreatedAt,
		&image.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (i *Images) Create(ctx context.Context, image *model.Image) error {
	image2 := *image
	image2.ID = newID()
	image2.CreatedAt = timestamp()

	qb := i.RunnableBuilder(ctx).
		Insert(i.table()).
		Columns(
			"id",
			"target_id",
			"target_name",
			"content",
			"description",
			"size",
			"content_type",
			"created_at").
		Values(
			image2.ID,
			image2.TargetID,
			image2.TargetName,
			image2.Content,
			image2.Description,
			image2.Size,
			image2.ContentType,
			image2.CreatedAt)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*image = image2
	return nil
}

func (i *Images) Update(ctx context.Context, image *model.Image) error {
	image2 := *image
	image2.UpdatedAt = timestamp()

	qb := i.RunnableBuilder(ctx).
		Update(i.table()).
		Set("target_id", image2.TargetID).
		Set("target_name", image2.TargetName).
		Set("content", image2.Content).
		Set("description", image2.Description).
		Set("size", image2.Size).
		Set("content_type", image2.ContentType).
		Set("updated_at", image2.UpdatedAt).
		Where("id = ?", image2.ID).
		Limit(1)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*image = image2
	return nil
}

func (i *Images) Delete(ctx context.Context, id string) (int64, error) {
	qb := i.RunnableBuilder(ctx).
		Delete(i.table()).
		Where("id = ?", id).
		Limit(1)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
