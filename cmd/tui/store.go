package tui

import (
	"context"
	"net/http"
	"strings"

	"github.com/pinmonl/pinmonl/pmapi"
)

type Store struct {
	client   *pmapi.Client
	pageSize int64

	pinls pinlSlice
	tags  tagSlice

	pinlCache pinlSlice
	tagCache  tagSlice

	pinlSearch string
	tagSearch  string
	tagFilter  tagSlice
}

func NewStore(endpoint string) *Store {
	httpClient := &http.Client{}
	client := pmapi.NewClient(endpoint, httpClient)
	return &Store{
		client:   client,
		pageSize: 20,
	}
}

func (s *Store) GetPinls(ctx context.Context) (pinlSlice, error) {
	if s.pinls == nil {
		info, err := s.client.GetPinlPageInfo(ctx)
		if err != nil {
			return nil, err
		}
		for i := int64(1); i <= info.Count; i += s.pageSize {
			page := i/s.pageSize + 1
			pinls, err := s.client.ListPinls(ctx, pmapi.PageOpts{
				Page:     page,
				PageSize: s.pageSize,
			})
			debugln("Get pinls: ", page)
			if err != nil {
				return nil, err
			}
			s.pinls = append(s.pinls, pinls...)
		}
	}
	if s.pinlCache == nil {
		cache := s.pinls
		if s.pinlSearch != "" {
			filtered := pinlSlice{}
			for _, p := range cache {
				if strings.Contains(p.Title, s.pinlSearch) {
					filtered = append(filtered, p)
				}
			}
			cache = filtered
		}
		if len(s.tagFilter) > 0 {
			filtered := pinlSlice{}
		filter:
			for _, p := range cache {
				for _, t := range s.tagFilter {
					hasTag := false
					for _, pt := range p.Tags {
						if t.Name == pt {
							hasTag = true
							break
						}
					}
					if !hasTag {
						continue filter
					}
				}
				filtered = append(filtered, p)
			}
			cache = filtered
		}
		s.pinlCache = cache
	}
	return s.pinlCache, nil
}

func (s *Store) CreatePinl(ctx context.Context, in *pmapi.Pinl) error {
	err := s.client.CreatePinl(ctx, in)
	if err != nil {
		return err
	}
	err = s.syncTags(ctx, in.Tags)
	if err != nil {
		return err
	}
	s.pinls = append(s.pinls, *in)
	s.pinlCache = nil
	return nil
}

func (s *Store) UpdatePinl(ctx context.Context, in *pmapi.Pinl) error {
	err := s.client.UpdatePinl(ctx, in)
	if err != nil {
		return err
	}
	err = s.syncTags(ctx, in.Tags)
	if err != nil {
		return err
	}
	s.pinls = s.pinls.replaceItem(*in)
	s.pinlCache = nil
	return nil
}

func (s *Store) DeletePinl(ctx context.Context, in *pmapi.Pinl) error {
	err := s.client.DeletePinl(ctx, in)
	if err != nil {
		return err
	}
	s.pinls = s.pinls.removeItem(*in)
	s.pinlCache = nil
	return nil
}

func (s *Store) clearPinls() {
	s.pinls = nil
	s.pinlCache = nil
}

func (s *Store) SetPinlSearch(search string) {
	s.pinlSearch = search
	s.pinlCache = nil
}

func (s *Store) syncTags(ctx context.Context, tagNames []string) error {
	flush := false
	for _, name := range tagNames {
		if s.tags.findName(name) == -1 {
			flush = true
			tag := pmapi.Tag{}
			tag.Name = name
			err := s.client.FindTagByName(ctx, &tag)
			if err != nil {
				return err
			}
			s.tags = append(s.tags, tag)
		}
	}
	if flush {
		s.tagCache = nil
	}
	return nil
}

func (s *Store) GetTags(ctx context.Context) (tagSlice, error) {
	if s.tags == nil {
		info, err := s.client.GetTagPageInfo(ctx)
		if err != nil {
			return nil, err
		}
		for i := int64(1); i <= info.Count; i += s.pageSize {
			page := i/s.pageSize + 1
			tags, err := s.client.ListTags(ctx, pmapi.PageOpts{
				Page:     page,
				PageSize: s.pageSize,
			})
			debugln("Get tags: ", page)
			if err != nil {
				return nil, err
			}
			s.tags = append(s.tags, tags...)
		}
	}
	if s.tagCache == nil {
		cache := s.tags
		if s.tagSearch != "" {
			filtered := tagSlice{}
			for _, tag := range cache {
				if strings.Contains(tag.Name, s.tagSearch) {
					filtered = append(filtered, tag)
				}
			}
			cache = filtered
		}
		s.tagCache = cache
	}
	return s.tagCache, nil
}

func (s *Store) CreateTag(ctx context.Context, in *pmapi.Tag) error {
	err := s.client.CreateTag(ctx, in)
	if err != nil {
		return err
	}
	s.tags = append(s.tags, *in)
	s.tagCache = nil
	return nil
}

func (s *Store) UpdateTag(ctx context.Context, in *pmapi.Tag) error {
	err := s.client.UpdateTag(ctx, in)
	if err != nil {
		return err
	}
	s.tags = s.tags.replaceItem(*in)
	s.tagCache = nil
	return nil
}

func (s *Store) DeleteTag(ctx context.Context, in *pmapi.Tag) error {
	err := s.client.DeleteTag(ctx, in)
	if err != nil {
		return err
	}
	s.tags = s.tags.removeItem(*in)
	s.deleteTagFromPinls(*in)
	s.tagCache = nil
	s.pinlCache = nil
	return nil
}

func (s *Store) deleteTagFromPinls(tag pmapi.Tag) error {
	for i, p := range s.pinls {
		for k, t := range p.Tags {
			if t == tag.Name {
				p.Tags = append(p.Tags[:k], p.Tags[k+1:]...)
				s.pinls[i] = p
				break
			}
		}
	}
	return nil
}

func (s *Store) clearTags() {
	s.tags = nil
	s.tagCache = nil
}

func (s *Store) SetTagSearch(search string) {
	s.tagSearch = search
	s.tagCache = nil
}

func (s *Store) AddTagFilter(tag pmapi.Tag) {
	s.tagFilter = append(s.tagFilter, tag)
	s.tagCache = nil
	s.pinlCache = nil
}

func (s *Store) DelTagFilter(tag pmapi.Tag) {
	s.tagFilter = s.tagFilter.removeItem(tag)
	s.tagCache = nil
	s.pinlCache = nil
}

func (s *Store) InTagFilter(tag pmapi.Tag) bool {
	return s.tagFilter.findIndex(tag) > -1
}

func (s *Store) ClearTagFilter(tag pmapi.Tag) {
	s.tagFilter = nil
	s.pinlCache = nil
}

type pinlSlice []pmapi.Pinl

func (ps pinlSlice) findIndex(item pmapi.Pinl) int {
	for i, p := range ps {
		if p.ID == item.ID {
			return i
		}
	}
	return -1
}

func (ps pinlSlice) replace(i int, item pmapi.Pinl) pinlSlice {
	return append(append(ps[:i], item), ps[i+1:]...)
}

func (ps pinlSlice) replaceItem(item pmapi.Pinl) pinlSlice {
	if i := ps.findIndex(item); i > -1 {
		return ps.replace(i, item)
	}
	return ps
}

func (ps pinlSlice) remove(i int) pinlSlice {
	return append(ps[:i], ps[i+1:]...)
}

func (ps pinlSlice) removeItem(item pmapi.Pinl) pinlSlice {
	if i := ps.findIndex(item); i > -1 {
		return ps.remove(i)
	}
	return ps
}

type tagSlice []pmapi.Tag

func (ts tagSlice) findIndex(item pmapi.Tag) int {
	return ts.findID(item.ID)
}

func (ts tagSlice) replace(i int, item pmapi.Tag) tagSlice {
	return append(append(ts[:i], item), ts[i+1:]...)
}

func (ts tagSlice) replaceItem(item pmapi.Tag) tagSlice {
	if i := ts.findIndex(item); i > -1 {
		return ts.replace(i, item)
	}
	return ts
}

func (ts tagSlice) remove(i int) tagSlice {
	return append(ts[:i], ts[i+1:]...)
}

func (ts tagSlice) removeItem(item pmapi.Tag) tagSlice {
	if i := ts.findIndex(item); i > -1 {
		return ts.remove(i)
	}
	return ts
}

func (ts tagSlice) findName(name string) int {
	for i, t := range ts {
		if t.Name == name {
			return i
		}
	}
	return -1
}

func (ts tagSlice) findID(id string) int {
	for i, t := range ts {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func (ts tagSlice) byParent() map[string]tagSlice {
	tm := map[string]tagSlice{}
	for _, t := range ts {
		k := t.ParentID
		tm[k] = append(tm[k], t)
	}
	return tm
}

func (ts tagSlice) Children(id string) tagSlice {
	out := tagSlice{}
	for _, t := range ts {
		if t.ParentID == id {
			out = append(out, t)
		}
	}
	return out
}

func (ts tagSlice) HasChildren(id string) bool {
	for _, t := range ts {
		if t.ParentID == id {
			return true
		}
	}
	return false
}

func (ts tagSlice) Keys() []string {
	out := []string{}
	for _, t := range ts {
		out = append(out, t.ID)
	}
	return out
}
