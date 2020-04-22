package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// PinpkgOpts defines the parameters for pinpkg filtering.
type PinpkgOpts struct {
	ListOpts
	PinlIDs []string
	PkgIDs  []string

	joinPkg  bool
	joinPinl bool
}

// PinpkgStore defines the service of pinpkg.
type PinpkgStore interface {
	List(context.Context, *PinpkgOpts) ([]model.Pinpkg, error)
	ListPinls(context.Context, *PinpkgOpts) (map[string][]model.Pinl, error)
	ListPkgs(context.Context, *PinpkgOpts) (map[string][]model.Pkg, error)
}

// NewPinpkgStore creates pinpkg store.
func NewPinpkgStore(s Store) PinpkgStore {
	return &dbPinpkgStore{s}
}

type dbPinpkgStore struct {
	Store
}

func (s *dbPinpkgStore) List(ctx context.Context, opts *PinpkgOpts) ([]model.Pinpkg, error) {
	e := s.Queryer(ctx)
	br, args := bindPinpkgOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pinpkg
	for rows.Next() {
		var m model.Pinpkg
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (s *dbPinpkgStore) ListPinls(ctx context.Context, opts *PinpkgOpts) (map[string][]model.Pinl, error) {
	return nil, nil
}

func (s *dbPinpkgStore) ListPkgs(ctx context.Context, opts *PinpkgOpts) (map[string][]model.Pkg, error) {
	return nil, nil
}

func bindPinpkgOpts(opts *PinpkgOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: pinpkgTB,
		Columns: database.NamespacedColumn(
			[]string{
				"pinl_id AS pinpkg_pinl_id",
				"pkg_id AS pinpkg_pkg_id",
				"sort AS pinpkg_sort",
			},
			pinpkgTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if len(opts.PinlIDs) > 0 {
		ks, ids := bindQueryIDs("pinl_ids", opts.PinlIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("%s.pinl_id IN (%s)", pinpkgTB, strings.Join(ks, ",")))
	}

	if len(opts.PkgIDs) > 0 {
		ks, ids := bindQueryIDs("pkg_ids", opts.PkgIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("%s.pkg_id IN (%s)", pinpkgTB, strings.Join(ks, ",")))
	}

	if opts.joinPinl {
		br.Columns = append(br.Columns, database.NamespacedColumn(
			[]string{
				"id AS pinl_id",
				"user_id AS pinl_user_id",
				"monl_id AS pinl_monl_id",
				"url AS pinl_url",
				"title AS pinl_title",
				"description AS pinl_description",
				"readme AS pinl_readme",
				"image_id AS pinl_image_id",
				"created_at AS pinl_created_at",
				"updated_at AS pinl_updated_at",
			},
			pinlTB,
		)...)
		br.Join = append(br.Join, fmt.Sprintf("%s ON %[1]s.id = %s.pinl_id", pinlTB, pinpkgTB))
	}
	if opts.joinPkg {
		br.Columns = append(br.Columns, database.NamespacedColumn(
			[]string{
				"id AS pkg_id",
				"url AS pkg_url",
				"provider AS pkg_provider",
				"provider_host AS pkg_provider_host",
				"provider_uri AS pkg_provider_uri",
				"title AS pkg_title",
				"description AS pkg_description",
				"readme AS pkg_readme",
				"image_id AS pkg_image_id",
				"labels AS pkg_labels",
				"created_at AS pkg_created_at",
				"updated_at AS pkg_updated_at",
			},
			pkgTB,
		)...)
		br.Join = append(br.Join, fmt.Sprintf("%s ON %[1]s.id = %s.pkg_id", pkgTB, pinpkgTB))
	}

	return br, args
}
