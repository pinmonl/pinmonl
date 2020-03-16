package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// TaggableOpts defines the parameters for taggable filtering.
type TaggableOpts struct {
	ListOpts
	Target  model.Morphable
	Targets []model.Morphable
	Tags    []model.Tag
	TagIDs  []string

	joinTag bool
}

// TaggableStore defines the services of taggable.
type TaggableStore interface {
	List(context.Context, *TaggableOpts) ([]model.Taggable, error)
	ListTags(context.Context, *TaggableOpts) (map[string][]model.Tag, error)
	Create(context.Context, *model.Taggable) error
	Delete(context.Context, *model.Taggable) error
	AssocTag(context.Context, model.Morphable, model.Tag) error
	AssocTags(context.Context, model.Morphable, []model.Tag) error
	DissocTag(context.Context, model.Morphable, model.Tag) error
	DissocTags(context.Context, model.Morphable, []model.Tag) error
	ClearTags(context.Context, model.Morphable) error
	ReAssocTags(context.Context, model.Morphable, []model.Tag) error
}

// NewTaggableStore creates taggable store.
func NewTaggableStore(s Store) TaggableStore {
	return &dbTaggableStore{s}
}

type dbTaggableStore struct {
	Store
}

// List retrieves taggables by the filter parameters..
func (s *dbTaggableStore) List(ctx context.Context, opts *TaggableOpts) ([]model.Taggable, error) {
	e := s.Queryer(ctx)
	br, args := bindTaggableOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Taggable
	for rows.Next() {
		var m model.Taggable
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// ListTags retrieves tags by taggable relationship.
func (s *dbTaggableStore) ListTags(ctx context.Context, opts *TaggableOpts) (map[string][]model.Tag, error) {
	if opts == nil {
		opts = &TaggableOpts{}
	}
	opts.joinTag = true

	e := s.Queryer(ctx)
	br, args := bindTaggableOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make(map[string][]model.Tag)
	for rows.Next() {
		var m model.Taggable
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		k := m.TargetID
		list[k] = append(list[k], *m.Tag)
	}
	return list, nil
}

// Create inserts the fields of taggable.
func (s *dbTaggableStore) Create(ctx context.Context, m *model.Taggable) error {
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: taggableTB,
		Fields: map[string]interface{}{
			"tag_id":      ":taggable_tag_id",
			"target_id":   ":taggable_target_id",
			"target_name": ":taggable_target_name",
			"sort":        ":taggable_sort",
		},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

// Delete removes taggable by the relationship.
func (s *dbTaggableStore) Delete(ctx context.Context, m *model.Taggable) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From: taggableTB,
		Where: []string{
			"target_id = :taggable_target_id",
			"target_name = :taggable_target_name",
			"tag_id = :taggable_tag_id",
		},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

// AssocTag associates single tag to the target.
func (s *dbTaggableStore) AssocTag(ctx context.Context, target model.Morphable, tag model.Tag) error {
	return s.Create(ctx, &model.Taggable{
		TargetID:   target.MorphKey(),
		TargetName: target.MorphName(),
		TagID:      tag.ID,
	})
}

// AssocTags associates multiple tags to the target.
func (s *dbTaggableStore) AssocTags(ctx context.Context, target model.Morphable, tags []model.Tag) error {
	for i, tag := range tags {
		err := s.Create(ctx, &model.Taggable{
			TargetID:   target.MorphKey(),
			TargetName: target.MorphName(),
			TagID:      tag.ID,
			Sort:       int64(i),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// DissocTag dissociates single tag from the target.
func (s *dbTaggableStore) DissocTag(ctx context.Context, target model.Morphable, tag model.Tag) error {
	return s.Delete(ctx, &model.Taggable{
		TargetID:   target.MorphKey(),
		TargetName: target.MorphName(),
		TagID:      tag.ID,
	})
}

// DissocTags dissociates multiple tags from the target.
func (s *dbTaggableStore) DissocTags(ctx context.Context, target model.Morphable, tags []model.Tag) error {
	for _, tag := range tags {
		err := s.Delete(ctx, &model.Taggable{
			TargetID:   target.MorphKey(),
			TargetName: target.MorphName(),
			TagID:      tag.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ClearTags removes all tag relations from the target.
func (s *dbTaggableStore) ClearTags(ctx context.Context, target model.Morphable) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  taggableTB,
		Where: []string{"target_id = :target_id", "target_name = :target_name"},
	}
	args := map[string]interface{}{
		"target_id":   target.MorphKey(),
		"target_name": target.MorphName(),
	}
	_, err := e.NamedExec(br.String(), args)
	return err
}

// ReAssocTags rebuilds the tag relation of the target.
func (s *dbTaggableStore) ReAssocTags(ctx context.Context, target model.Morphable, tags []model.Tag) error {
	err := s.ClearTags(ctx, target)
	if err != nil {
		return err
	}
	return s.AssocTags(ctx, target, tags)
}

func bindTaggableOpts(opts *TaggableOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: taggableTB,
		Columns: database.NamespacedColumn(
			[]string{
				"tag_id AS taggable_tag_id",
				"target_id AS taggable_target_id",
				"target_name AS taggable_target_name",
				"sort AS taggable_sort",
			},
			taggableTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.Target != nil {
		opts.Targets = append(opts.Targets, opts.Target)
	}
	if opts.Targets != nil {
		tns := (model.MorphableList)(opts.Targets).Names()
		tks := (model.MorphableList)(opts.Targets).Keys()
		ks, ids := bindQueryIDs("target_ids", tks)
		args.Set("target_name", tns[0])
		args.AppendStringMap(ids)
		br.Where = append(br.Where,
			fmt.Sprintf("target_name = :target_name"),
			fmt.Sprintf("target_id IN (%s)", strings.Join(ks, ",")),
		)
	}

	if opts.Tags != nil {
		ids := make([]string, len(opts.Tags))
		for i, t := range opts.Tags {
			ids[i] = t.ID
		}
		opts.TagIDs = ids
	}
	if opts.TagIDs != nil {
		ks, ids := bindQueryIDs("tag_ids", opts.TagIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("tag_id IN (%s)", strings.Join(ks, ",")))
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
		br.Join = append(br.Join, fmt.Sprintf("INNER JOIN %s ON %[1]s.id = %s.tag_id", tagTB, taggableTB))
	}

	return br, args
}
