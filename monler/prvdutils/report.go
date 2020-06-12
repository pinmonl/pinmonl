package prvdutils

import (
	"errors"
	"io"
	"sync"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
)

var (
	ErrCursorInit = errors.New("monler report: cursor is not started yet")
)

type StaticReport struct {
	*pkguri.PkgURI
	*sync.Mutex
	stats  []*model.Stat
	tags   []*model.Stat
	cursor int
}

func NewStaticReport(pu *pkguri.PkgURI, stats, tags []*model.Stat) *StaticReport {
	return &StaticReport{
		PkgURI: pu,
		Mutex:  &sync.Mutex{},
		stats:  stats,
		tags:   tags,
		cursor: -1,
	}
}

func (s *StaticReport) URI() (*pkguri.PkgURI, error) {
	return s.PkgURI, nil
}

func (s *StaticReport) Stats() ([]*model.Stat, error) {
	return s.stats, nil
}

func (s *StaticReport) Next() bool {
	s.Lock()
	defer s.Unlock()
	if s.cursor+1 < len(s.tags) {
		s.cursor++
		return true
	}
	return false
}

func (s *StaticReport) Tag() (*model.Stat, error) {
	s.Lock()
	defer s.Unlock()
	if s.cursor < 0 {
		return nil, ErrCursorInit
	}
	if s.cursor >= len(s.tags) {
		return nil, io.EOF
	}
	return s.tags[s.cursor], nil
}

func (s *StaticReport) Close() error {
	return nil
}
