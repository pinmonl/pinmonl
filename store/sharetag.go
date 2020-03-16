package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// SharetagOpts defines the parameters for share tag filtering.
type SharetagOpts struct {
	ListOpts
	ShareID  string
	ShareIDs []string
	Kind     model.SharetagKind
	TagID    string
	TagIDs   []string

	joinTag   bool
	joinShare bool
}

// SharetagStore defines the services of share tag.
type SharetagStore interface {
	List(context.Context, *SharetagOpts) ([]model.Sharetag, error)
	ListTags(context.Context, *SharetagOpts) (map[string][]model.Tag, error)
	ListShares(context.Context, *SharetagOpts) (map[string][]model.Share, error)
	Create(context.Context, *model.Sharetag) error
	Delete(context.Context, *model.Sharetag) error
	AssocTag(context.Context, model.Share, model.SharetagKind, model.Tag) error
	AssocTags(context.Context, model.Share, model.SharetagKind, []model.Tag) error
	DissocTag(context.Context, model.Share, model.SharetagKind, model.Tag) error
	DissocTags(context.Context, model.Share, model.SharetagKind, []model.Tag) error
	ClearByKind(context.Context, model.Share, model.SharetagKind) error
	ReAssocTags(context.Context, model.Share, model.SharetagKind, []model.Tag) error
}

// NewSharetagStore creates share tag store.
func NewSharetagStore(s Store) SharetagStore {
	return &dbSharetagStore{s}
}

type dbSharetagStore struct {
	Store
}

// List retrieves share tags by the filter parameters.
func (s *dbSharetagStore) List(ctx context.Context, opts *SharetagOpts) ([]model.Sharetag, error) {
	e := s.Queryer(ctx)
	br, args := bindSharetagOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Sharetag
	for rows.Next() {
		var m model.Sharetag
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// ListTags retrieves tags by share tag relationship.
func (s *dbSharetagStore) ListTags(ctx context.Context, opts *SharetagOpts) (map[string][]model.Tag, error) {
	if opts == nil {
		opts = &SharetagOpts{}
	}
	opts.joinTag = true

	e := s.Queryer(ctx)
	br, args := bindSharetagOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make(map[string][]model.Tag)
	for rows.Next() {
		var m model.Sharetag
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		k := m.ShareID
		list[k] = append(list[k], *m.Tag)
	}
	return list, nil
}

// ListShares retrieves shares by share tag relationship.
func (s *dbSharetagStore) ListShares(ctx context.Context, opts *SharetagOpts) (map[string][]model.Share, error) {
	if opts == nil {
		opts = &SharetagOpts{}
	}
	opts.joinShare = true

	e := s.Queryer(ctx)
	br, args := bindSharetagOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make(map[string][]model.Share)
	for rows.Next() {
		var m model.Sharetag
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		k := m.TagID
		list[k] = append(list[k], *m.Share)
	}
	return list, nil
}

// Create inserts the fields of share tag.
func (s *dbSharetagStore) Create(ctx context.Context, m *model.Sharetag) error {
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: sharetagTB,
		Fields: map[string]interface{}{
			"share_id":  ":sharetag_share_id",
			"tag_id":    ":sharetag_tag_id",
			"kind":      ":sharetag_kind",
			"parent_id": ":sharetag_parent_id",
			"sort":      ":sharetag_sort",
			"level":     ":sharetag_level",
		},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

// Delete removes share tag by the relationship.
func (s *dbSharetagStore) Delete(ctx context.Context, m *model.Sharetag) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  sharetagTB,
		Where: []string{"share_id = :sharetag_share_id", "tag_id = :sharetag_tag_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func (s *dbSharetagStore) AssocTag(ctx context.Context, share model.Share, kind model.SharetagKind, tag model.Tag) error {
	return s.Create(ctx, &model.Sharetag{
		ShareID:  share.ID,
		TagID:    tag.ID,
		Kind:     kind,
		ParentID: tag.ParentID,
	})
}

func (s *dbSharetagStore) AssocTags(ctx context.Context, share model.Share, kind model.SharetagKind, tags []model.Tag) error {
	for i, t := range tags {
		err := s.Create(ctx, &model.Sharetag{
			ShareID:  share.ID,
			TagID:    t.ID,
			Kind:     kind,
			ParentID: t.ParentID,
			Sort:     int64(i),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *dbSharetagStore) DissocTag(ctx context.Context, share model.Share, kind model.SharetagKind, tag model.Tag) error {
	return s.Delete(ctx, &model.Sharetag{
		ShareID: share.ID,
		TagID:   tag.ID,
	})
}

func (s *dbSharetagStore) DissocTags(ctx context.Context, share model.Share, kind model.SharetagKind, tags []model.Tag) error {
	for _, t := range tags {
		err := s.Delete(ctx, &model.Sharetag{
			ShareID: share.ID,
			TagID:   t.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *dbSharetagStore) ClearByKind(ctx context.Context, share model.Share, kind model.SharetagKind) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  sharetagTB,
		Where: []string{"share_id = :share_id", "kind = :kind"},
	}
	args := map[string]interface{}{
		"share_id": share.ID,
		"kind":     kind,
	}
	_, err := e.NamedExec(br.String(), args)
	return err
}

func (s *dbSharetagStore) ReAssocTags(ctx context.Context, share model.Share, kind model.SharetagKind, tags []model.Tag) error {
	err := s.ClearByKind(ctx, share, kind)
	if err != nil {
		return err
	}
	return s.AssocTags(ctx, share, kind, tags)
}

func bindSharetagOpts(opts *SharetagOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: sharetagTB,
		Columns: database.NamespacedColumn(
			[]string{
				"share_id AS sharetag_share_id",
				"tag_id AS sharetag_tag_id",
				"kind AS sharetag_kind",
				"parent_id AS sharetag_parent_id",
				"sort AS sharetag_sort",
				"level AS sharetag_level",
			},
			sharetagTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.ShareID != "" {
		opts.ShareIDs = append(opts.ShareIDs, opts.ShareID)
	}
	if opts.ShareIDs != nil {
		ks, ids := bindQueryIDs("share_ids", opts.ShareIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("share_id in (%s)", strings.Join(ks, ",")))
	}

	if opts.TagID != "" {
		opts.TagIDs = append(opts.TagIDs, opts.TagID)
	}
	if opts.TagIDs != nil {
		ks, ids := bindQueryIDs("tag_ids", opts.TagIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("tag_id in (%s)", strings.Join(ks, ",")))
	}

	if opts.Kind != model.SharetagKindEmpty {
		br.Where = append(br.Where, "kind = :kind")
		args["kind"] = opts.Kind
	}

	if opts.joinTag {
		br.Columns = append(br.Columns, database.NamespacedColumn(
			[]string{
				"id AS tag_id",
				"name AS tag_name",
				"user_id AS tag_user_id",
				"parent_id AS tag_parent_id",
				"sort AS tag_sort",
				"level AS tag_level",
				"color AS tag_color",
				"bgcolor AS tag_bgcolor",
				"created_at AS tag_created_at",
				"updated_at AS tag_updated_at",
			},
			tagTB,
		)...)
		br.Join = append(br.Join, fmt.Sprintf("INNER JOIN %s ON %[1]s.id = %s.tag_id", tagTB, sharetagTB))
	}
	if opts.joinShare {
		br.Columns = append(br.Columns, database.NamespacedColumn(
			[]string{
				"id AS share_id",
				"user_id AS share_user_id",
				"name AS share_name",
				"description AS share_description",
				"readme AS share_readme",
				"image_id AS share_image_id",
				"created_at AS share_created_at",
				"updated_at AS share_updated_at",
			},
			shareTB,
		)...)
		br.Join = append(br.Join, fmt.Sprintf("INNER JOIN %s ON %[1]s.id = %s.share_id", shareTB, sharetagTB))
	}

	return br, args
}
