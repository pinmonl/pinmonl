package monl

import (
	"time"

	"github.com/pinmonl/pinmonl/model/field"
)

// Stat defines information of the repo at the time.
type Stat interface {
	Date() time.Time
	Group() string
	Value() string
	Manifest() field.Labels
}

// NewStat creates simple stat which does not need handling.
func NewStat(date time.Time, group, value string, manifest field.Labels) Stat {
	return &basicStat{
		date:     date,
		group:    group,
		value:    value,
		manifest: manifest,
	}
}

// basicStat stores stat which does not require extra handling.
type basicStat struct {
	date     time.Time
	group    string
	value    string
	manifest field.Labels
}

// Date returns the date of stat.
func (bs *basicStat) Date() time.Time { return bs.date }

// Group returns the group of stat.
func (bs *basicStat) Group() string { return bs.group }

// Value returns the value of stat.
func (bs *basicStat) Value() string { return bs.value }

// Manifest returns the manifest of stat.
func (bs *basicStat) Manifest() field.Labels { return bs.manifest }

// StatCollection holds array of stat.
type StatCollection []Stat

// FindGroup is a shorthand to filter collection by group.
func (sc StatCollection) FindGroup(group string) StatCollection {
	out := make([]Stat, 0)
	for _, s := range sc {
		if s.Group() == group {
			out = append(out, s)
		}
	}
	return out
}
