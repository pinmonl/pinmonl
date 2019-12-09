package monl

import "time"

// Stat defines information of the repo at the time
type Stat interface {
	Date() time.Time
	Group() string
	Value() string
	Manifest() Manifest
}

// Manifest holds extra info for stat
// e.g. OS, Architecture
type Manifest map[string]string

// NewManifestFromString converts string into mnnifest map
func NewManifestFromString(raw string) Manifest {
	// TODO
	return nil
}

// String converts manifest map back into string
func (m *Manifest) String() string {
	// TODO
	return ""
}

// StatCollection holds array of stat
type StatCollection []Stat

// FindGroup is a shorthand to filter collection by group
func (s StatCollection) FindGroup(group string) StatCollection {
	out := make([]Stat, 0)
	for _, stat := range s {
		if stat.Group() == group {
			out = append(out, stat)
		}
	}
	return out
}

// SimpleStat stores stat which does not require extra handling
type SimpleStat struct {
	date     time.Time
	group    string
	value    string
	manifest Manifest
}

// NewStat creates simple stat which does not need handling
func NewStat(date time.Time, group, value string, manifest Manifest) *SimpleStat {
	return &SimpleStat{
		date:     date,
		group:    group,
		value:    value,
		manifest: manifest,
	}
}

// Date returns the date of stat
func (s *SimpleStat) Date() time.Time { return s.date }

// Group returns the group of stat
func (s *SimpleStat) Group() string { return s.group }

// Value returns the value of stat
func (s *SimpleStat) Value() string { return s.value }

// Manifest returns the manifest of stat
func (s *SimpleStat) Manifest() Manifest { return s.manifest }
