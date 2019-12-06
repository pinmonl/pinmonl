package monl

import "time"

type Stat interface {
	Date() time.Time
	Group() string
	Value() string
	Manifest() Manifest
}

type Manifest map[string]string

func NewManifestFromString(raw string) Manifest {
	// TODO
	return nil
}

func (m *Manifest) ToString() string {
	// TODO
	return ""
}

type StatCollection []Stat

func (s StatCollection) FindGroup(group string) StatCollection {
	out := make([]Stat, 0)
	for _, stat := range s {
		if stat.Group() == group {
			out = append(out, stat)
		}
	}
	return out
}

type SimpleStat struct {
	date     time.Time
	group    string
	value    string
	manifest Manifest
}

func NewStat(date time.Time, group, value string, manifest Manifest) *SimpleStat {
	return &SimpleStat{
		date:     date,
		group:    group,
		value:    value,
		manifest: manifest,
	}
}

func (s *SimpleStat) Date() time.Time { return s.date }

func (s *SimpleStat) Group() string { return s.group }

func (s *SimpleStat) Value() string { return s.value }

func (s *SimpleStat) Manifest() Manifest { return s.manifest }
