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
	pkg    *model.Pkg
	stats  []*model.Stat
	tags   []*model.Stat
	cursor int
}

func NewStaticReport(pu *pkguri.PkgURI, pkg *model.Pkg, stats, tags []*model.Stat) *StaticReport {
	return &StaticReport{
		PkgURI: pu,
		Mutex:  &sync.Mutex{},
		pkg:    pkg,
		stats:  stats,
		tags:   tags,
		cursor: -1,
	}
}

func (s *StaticReport) URI() (*pkguri.PkgURI, error) {
	return s.PkgURI, nil
}

func (s *StaticReport) Pkg() (*model.Pkg, error) {
	return s.pkg, nil
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

type PageFunc func(page int64) (items []*model.Stat, total int64, hasNext bool, err error)

// PagesReport
type PagesReport struct {
	*pkguri.PkgURI
	*sync.Mutex
	pkg        *model.Pkg
	stats      []*model.Stat
	statPage   int64
	statFn     PageFunc
	tags       []*model.Stat
	tagPage    int64
	tagHasNext bool
	tagFn      PageFunc
	tagTotal   int64
	cursor     int
}

func NewPagesReport(pu *pkguri.PkgURI, pkg *model.Pkg, statFn, tagFn PageFunc) *PagesReport {
	report := &PagesReport{
		PkgURI:     pu,
		Mutex:      &sync.Mutex{},
		pkg:        pkg,
		statPage:   1,
		statFn:     statFn,
		tagPage:    1,
		tagFn:      tagFn,
		tagHasNext: true,
		tagTotal:   0,
		cursor:     -1,
	}
	return report
}

func (p *PagesReport) URI() (*pkguri.PkgURI, error) {
	return p.PkgURI, nil
}

func (p *PagesReport) Pkg() (*model.Pkg, error) {
	return p.pkg, nil
}

func (p *PagesReport) Stats() ([]*model.Stat, error) {
	if p.stats == nil {
		stats, _, _, err := p.statFn(p.statPage)
		if err != nil {
			return nil, err
		}
		p.stats = stats
	}
	return p.stats, nil
}

func (p *PagesReport) Next() bool {
	p.Lock()
	defer p.Unlock()

	if p.cursor+1 >= len(p.tags) {
		if p.tagHasNext {
			if err := p.fetchNextTags(); err != nil {
				return false
			}
		}
	}

	if p.cursor+1 < len(p.tags) {
		p.cursor++
		return true
	}

	return false
}

func (p *PagesReport) Tag() (*model.Stat, error) {
	p.Lock()
	defer p.Unlock()

	if p.cursor < 0 {
		return nil, ErrCursorInit
	}
	if p.cursor >= len(p.tags) {
		return nil, io.EOF
	}
	return p.tags[p.cursor], nil
}

func (p *PagesReport) fetchNextTags() error {
	tags, total, hasNext, err := p.tagFn(p.tagPage)
	if err != nil {
		return err
	}

	p.tags = tags
	p.tagTotal = total
	p.tagHasNext = hasNext
	p.tagPage++
	p.cursor = -1
	return nil
}

func (p *PagesReport) Close() error {
	return nil
}
