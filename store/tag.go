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
	Target   model.Morphable
	Targets  []model.Morphable
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
	rows, err := e.NamedQuery(br.String(), args)
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
	rows, err := e.NamedQuery(br.String(), args)
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
	br.Where = []string{"id = :id"}
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
	br.Where = []string{"user_id = :user_id", "name = :name"}
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
	stmt := database.InsertBuilder{
		Into: tagTB,
		Fields: map[string]interface{}{
			"id":         nil,
			"name":       nil,
			"user_id":    nil,
			"parent_id":  nil,
			"sort":       nil,
			"level":      nil,
			"created_at": nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m2)
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
	stmt := database.UpdateBuilder{
		From: tagTB,
		Fields: map[string]interface{}{
			"name":       nil,
			"user_id":    nil,
			"parent_id":  nil,
			"sort":       nil,
			"level":      nil,
			"updated_at": nil,
		},
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes tag by id.
func (s *dbTagStore) Delete(ctx context.Context, m *model.Tag) error {
	e := s.Execer(ctx)
	stmt := database.DeleteBuilder{
		From:  tagTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindTagOpts(opts *TagOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{
		From: tagTB,
		Columns: database.NamespacedColumn(
			[]string{"id", "name", "user_id", "parent_id", "sort", "level", "created_at", "updated_at"},
			tagTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := make(map[string]interface{})

	if opts.IDs != nil {
		ks := make([]string, len(opts.IDs))
		for i, id := range opts.IDs {
			k := fmt.Sprintf("id%d", i)
			args[k] = id
			ks[i] = ":" + k
		}
		br.Where = append(br.Where, fmt.Sprintf("id IN (%s)", strings.Join(ks, ", ")))
	}

	if opts.Target != nil {
		opts.Targets = append(opts.Targets, opts.Target)
	}
	if opts.Targets != nil && len(opts.Targets) > 0 {
		br.Columns = append(br.Columns, "b.target_id")
		br.Columns = append(br.Columns, "b.target_name")
		br.Join = append(br.Join, fmt.Sprintf(`INNER JOIN %s AS b ON %s.id = b.tag_id`, taggableTB, tagTB))
		br.Where = append(br.Where, "b.target_name = :target_name")
		args["target_name"] = opts.Targets[0].MorphName()

		tids := make([]string, len(opts.Targets))
		for i, t := range opts.Targets {
			tids[i] = t.MorphKey()
		}
		ks, ids := bindQueryIDs("target_ids", tids)
		for k, id := range ids {
			args[k] = id
		}
		br.Where = append(br.Where, fmt.Sprintf("b.target_id IN (%s)", strings.Join(ks, ",")))
	}

	if opts.Name != "" {
		opts.Names = append(opts.Names, opts.Name)
	}
	if opts.Names != nil {
		ks := make([]string, len(opts.Names))
		for i, id := range opts.Names {
			k := fmt.Sprintf("name%d", i)
			args[k] = id
			ks[i] = ":" + k
		}
		br.Where = append(br.Where, fmt.Sprintf("name IN (%s)", strings.Join(ks, ", ")))
	}

	if opts.UserID != "" {
		br.Where = append(br.Where, "user_id = :user_id")
		args["user_id"] = opts.UserID
	}

	return br, args
}
