package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// TagOpts defines the parameters for tag filtering.
type TagOpts struct {
	ListOpts
	IDs      []string
	ParentID string
	UserID   string
	Name     string
	Names    []string
}

// TagStore defines the services of tag.
type TagStore interface {
	List(context.Context, *TagOpts) ([]model.Tag, error)
	Count(context.Context, *TagOpts) (int64, error)
	Find(context.Context, *model.Tag) error
	FindByName(context.Context, *model.Tag) error
	Create(context.Context, *model.Tag) error
	Update(context.Context, *model.Tag) error
	Delete(context.Context, *model.Tag) error
}

// NewTagStore creates tag store.
func NewTagStore(s Store) TagStore {
	return &dbTagStore{s}
}

type dbTagStore struct {
	Store
}

// List retrieves tags by the filter parameters.
func (s *dbTagStore) List(ctx context.Context, opts *TagOpts) ([]model.Tag, error) {
	e := s.Queryer(ctx)
	br, args := bindTagOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Tag
	for rows.Next() {
		var m model.Tag
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Count counts number of tags with the filter parameters.
func (s *dbTagStore) Count(ctx context.Context, opts *TagOpts) (int64, error) {
	e := s.Queryer(ctx)
	br, args := bindTagOpts(opts)
	br.Columns = []string{"COUNT(*) as count"}
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, sql.ErrNoRows
	}
	var count int64
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Find retrieves tag by id.
func (s *dbTagStore) Find(ctx context.Context, m *model.Tag) error {
	e := s.Queryer(ctx)
	br, _ := bindTagOpts(nil)
	br.Where = []string{"id = :tag_id"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	var m2 model.Tag
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// FindByName retrieves tag by user and tag name.
func (s *dbTagStore) FindByName(ctx context.Context, m *model.Tag) error {
	e := s.Queryer(ctx)
	br, _ := bindTagOpts(nil)
	br.Where = []string{"user_id = :tag_user_id", "name = :tag_name"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	var m2 model.Tag
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Create inserts the fields of tag with generated id.
func (s *dbTagStore) Create(ctx context.Context, m *model.Tag) error {
	m2 := *m
	m2.ID = newUID()
	m2.CreatedAt = timestamp()
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: tagTB,
		Fields: map[string]interface{}{
			"id":         ":tag_id",
			"name":       ":tag_name",
			"user_id":    ":tag_user_id",
			"parent_id":  ":tag_parent_id",
			"sort":       ":tag_sort",
			"level":      ":tag_level",
			"color":      ":tag_color",
			"bgcolor":    ":tag_bgcolor",
			"created_at": ":tag_created_at",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates the fields of tag by id.
func (s *dbTagStore) Update(ctx context.Context, m *model.Tag) error {
	m2 := *m
	m2.UpdatedAt = timestamp()
	e := s.Execer(ctx)
	br := database.UpdateBuilder{
		From: tagTB,
		Fields: map[string]interface{}{
			"name":       ":tag_name",
			"user_id":    ":tag_user_id",
			"parent_id":  ":tag_parent_id",
			"sort":       ":tag_sort",
			"level":      ":tag_level",
			"color":      ":tag_color",
			"bgcolor":    ":tag_bgcolor",
			"updated_at": ":tag_updated_at",
		},
		Where: []string{"id = :tag_id"},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes tag by id.
func (s *dbTagStore) Delete(ctx context.Context, m *model.Tag) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  tagTB,
		Where: []string{"id = :tag_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindTagOpts(opts *TagOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: tagTB,
		Columns: database.NamespacedColumn(
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
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.IDs != nil {
		ks, ids := bindQueryIDs("ids", opts.IDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("id IN (%s)", strings.Join(ks, ", ")))
	}

	if opts.Name != "" {
		opts.Names = append(opts.Names, opts.Name)
	}
	if opts.Names != nil {
		ks, names := bindQueryIDs("names", opts.Names)
		args.AppendStringMap(names)
		br.Where = append(br.Where, fmt.Sprintf("name IN (%s)", strings.Join(ks, ",")))
	}

	if opts.UserID != "" {
		br.Where = append(br.Where, "user_id = :user_id")
		args["user_id"] = opts.UserID
	}

	return br, args
}
