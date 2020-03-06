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
	br.From = taggableTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
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
	e := s.Queryer(ctx)
	br, args := bindTaggableOpts(opts)
	br.From = taggableTB
	br.Join = []string{fmt.Sprintf("INNER JOIN %s ON %s.tag_id = %s.id", tagTB, taggableTB, tagTB)}
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make(map[string][]model.Tag)
	var r struct {
		model.Tag
		model.Taggable
	}
	for rows.Next() {
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		k := r.Taggable.TargetID
		list[k] = append(list[k], r.Tag)
	}
	return list, nil
}

// Create inserts the fields of taggable.
func (s *dbTaggableStore) Create(ctx context.Context, m *model.Taggable) error {
	e := s.Execer(ctx)
	stmt := database.InsertBuilder{
		Into: taggableTB,
		Fields: map[string]interface{}{
			"tag_id":      nil,
			"target_id":   nil,
			"target_name": nil,
			"sort":        nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

// Delete removes taggable by the relationship.
func (s *dbTaggableStore) Delete(ctx context.Context, m *model.Taggable) error {
	e := s.Execer(ctx)
	stmt := database.DeleteBuilder{
		From: taggableTB,
		Where: []string{
			"target_id = :target_id",
			"target_name = :target_name",
			"tag_id = :tag_id",
		},
	}.String()
	_, err := e.NamedExec(stmt, m)
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
	stmt := database.DeleteBuilder{
		From:  taggableTB,
		Where: []string{"target_id = :target_id", "target_name = :target_name"},
	}.String()
	args := map[string]interface{}{
		"target_id":   target.MorphKey(),
		"target_name": target.MorphName(),
	}
	_, err := e.NamedExec(stmt, args)
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

func bindTaggableOpts(opts *TaggableOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	args := make(map[string]interface{})
	if opts.Target != nil {
		br.Where = append(br.Where, "target_id = :target_id")
		br.Where = append(br.Where, "target_name = :target_name")
		args["target_id"] = opts.Target.MorphKey()
		args["target_name"] = opts.Target.MorphName()
	}

	if opts.Tags != nil {
		ids := make([]string, len(opts.Tags))
		for i, t := range opts.Tags {
			ids[i] = t.ID
		}
	}
	if opts.TagIDs != nil {
		ks, ids := bindQueryIDs("tag_ids", opts.TagIDs)
		for k, id := range ids {
			args[k] = id
		}
		br.Where = append(br.Where, fmt.Sprintf("tag_id IN (%s)", strings.Join(ks, ",")))
	}

	return br, args
}
