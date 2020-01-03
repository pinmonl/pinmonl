package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// ShareTagOpts defines the parameters for share tag filtering.
type ShareTagOpts struct {
	ListOpts
	ShareID  string
	ShareIDs []string
	Kind     model.ShareTagKind
	TagID    string
	TagIDs   []string
}

// ShareTagStore defines the services of share tag.
type ShareTagStore interface {
	List(context.Context, *ShareTagOpts) ([]model.ShareTag, error)
	ListTags(context.Context, *ShareTagOpts) (map[string][]model.Tag, error)
	Create(context.Context, *model.ShareTag) error
	Delete(context.Context, *model.ShareTag) error
	AssocTag(context.Context, model.Share, model.ShareTagKind, model.Tag) error
	AssocTags(context.Context, model.Share, model.ShareTagKind, []model.Tag) error
	DissocTag(context.Context, model.Share, model.ShareTagKind, model.Tag) error
	DissocTags(context.Context, model.Share, model.ShareTagKind, []model.Tag) error
	ClearByKind(context.Context, model.Share, model.ShareTagKind) error
	ReAssocTags(context.Context, model.Share, model.ShareTagKind, []model.Tag) error
}

// NewShareTagStore creates share tag store.
func NewShareTagStore(s Store) ShareTagStore {
	return &dbShareTagStore{s}
}

type dbShareTagStore struct {
	Store
}

// List retrieves share tags by the filter parameters.
func (s *dbShareTagStore) List(ctx context.Context, opts *ShareTagOpts) ([]model.ShareTag, error) {
	e := s.Exter(ctx)
	br, args := bindShareTagOpts(opts)
	br.From = shareTagTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.ShareTag
	for rows.Next() {
		var m model.ShareTag
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// ListTags retrieves tags by share tag relationship.
func (s *dbShareTagStore) ListTags(ctx context.Context, opts *ShareTagOpts) (map[string][]model.Tag, error) {
	e := s.Exter(ctx)
	br, args := bindShareTagOpts(opts)
	br.From = shareTagTB
	br.Join = []string{fmt.Sprintf("INNER JOIN %s ON %s.tag_id = %s.id", tagTB, shareTagTB, tagTB)}
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make(map[string][]model.Tag)
	var r struct {
		model.Tag
		model.ShareTag
	}
	for rows.Next() {
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		k := r.ShareTag.ShareID
		list[k] = append(list[k], r.Tag)
	}
	return list, nil
}

// Create inserts the fields of share tag.
func (s *dbShareTagStore) Create(ctx context.Context, m *model.ShareTag) error {
	e := s.Exter(ctx)
	stmt := database.InsertBuilder{
		Into: shareTagTB,
		Fields: map[string]interface{}{
			"share_id": nil,
			"tag_id":   nil,
			"kind":     nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

// Delete removes share tag by the relationship.
func (s *dbShareTagStore) Delete(ctx context.Context, m *model.ShareTag) error {
	e := s.Exter(ctx)
	stmt := database.DeleteBuilder{
		From:  shareTagTB,
		Where: []string{"share_id = :share_id", "tag_id = :tag_id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func (s *dbShareTagStore) AssocTag(ctx context.Context, share model.Share, kind model.ShareTagKind, tag model.Tag) error {
	return s.Create(ctx, &model.ShareTag{
		ShareID: share.ID,
		TagID:   tag.ID,
		Kind:    string(kind),
	})
}

func (s *dbShareTagStore) AssocTags(ctx context.Context, share model.Share, kind model.ShareTagKind, tags []model.Tag) error {
	for _, t := range tags {
		err := s.Create(ctx, &model.ShareTag{
			ShareID: share.ID,
			TagID:   t.ID,
			Kind:    string(kind),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *dbShareTagStore) DissocTag(ctx context.Context, share model.Share, kind model.ShareTagKind, tag model.Tag) error {
	return s.Delete(ctx, &model.ShareTag{
		ShareID: share.ID,
		TagID:   tag.ID,
	})
}

func (s *dbShareTagStore) DissocTags(ctx context.Context, share model.Share, kind model.ShareTagKind, tags []model.Tag) error {
	for _, t := range tags {
		err := s.Delete(ctx, &model.ShareTag{
			ShareID: share.ID,
			TagID:   t.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *dbShareTagStore) ClearByKind(ctx context.Context, share model.Share, kind model.ShareTagKind) error {
	e := s.Exter(ctx)
	stmt := database.DeleteBuilder{
		From:  shareTagTB,
		Where: []string{"share_id = :share_id", "kind = :kind"},
	}.String()
	args := map[string]interface{}{
		"share_id": share.ID,
		"kind":     kind,
	}
	_, err := e.NamedExec(stmt, args)
	return err
}

func (s *dbShareTagStore) ReAssocTags(ctx context.Context, share model.Share, kind model.ShareTagKind, tags []model.Tag) error {
	err := s.ClearByKind(ctx, share, kind)
	if err != nil {
		return err
	}
	return s.AssocTags(ctx, share, kind, tags)
}

func bindShareTagOpts(opts *ShareTagOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	br = bindListOpts(opts.ListOpts)
	args := make(map[string]interface{})

	if opts.ShareID != "" {
		opts.ShareIDs = append(opts.ShareIDs, opts.ShareID)
	}
	if opts.ShareIDs != nil {
		ks := make([]string, len(opts.ShareIDs))
		for i, id := range opts.ShareIDs {
			k := fmt.Sprintf("share_id%d", i)
			ks[i] = ":" + k
			args[k] = id
		}
		br.Where = append(br.Where, fmt.Sprintf("share_id in (%s)", strings.Join(ks, ", ")))
	}

	if opts.TagID != "" {
		opts.TagIDs = append(opts.TagIDs, opts.TagID)
	}
	if opts.TagIDs != nil {
		ks := make([]string, len(opts.TagIDs))
		for i, id := range opts.TagIDs {
			k := fmt.Sprintf("tag_id%d", i)
			ks[i] = ":" + k
			args[k] = id
		}
		br.Where = append(br.Where, fmt.Sprintf("tag_id in (%s)", strings.Join(ks, ", ")))
	}

	if opts.Kind != "" {
		br.Where = append(br.Where, "kind = :kind")
		args["kind"] = opts.Kind
	}

	return br, args
}
