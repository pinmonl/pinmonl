package store

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
)

type PinpkgOpts struct {
	ListOpts
}

type PinpkgStore interface {
	List(context.Context, *PinpkgOpts) ([]model.Pinpkg, error)
	ListPinls(context.Context, *PinpkgOpts) (map[string][]model.Pinl, error)
	ListPkgs(context.Context, *PinpkgOpts) (map[string][]model.Pkg, error)
}

func NewPinpkgStore(s Store) PinpkgStore {
	return &dbPinpkgStore{s}
}

type dbPinpkgStore struct {
	Store
}

func (s *dbPinpkgStore) List(ctx context.Context, opts *PinpkgOpts) ([]model.Pinpkg, error) {
	return nil, nil
}

func (s *dbPinpkgStore) ListPinls(ctx context.Context, opts *PinpkgOpts) (map[string][]model.Pinl, error) {
	return nil, nil
}

func (s *dbPinpkgStore) ListPkgs(ctx context.Context, opts *PinpkgOpts) (map[string][]model.Pkg, error) {
	return nil, nil
}
