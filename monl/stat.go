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
	Labels() field.Labels
	Substats() []Stat
}

// NewStat creates simple stat which does not need handling.
func NewStat(date time.Time, group, value string, labels field.Labels, substats []Stat) Stat {
	return &basicStat{
		date:     date,
		group:    group,
		value:    value,
		labels:   labels,
		substats: substats,
	}
}

// basicStat stores stat which does not require extra handling.
type basicStat struct {
	date     time.Time
	group    string
	value    string
	labels   field.Labels
	substats []Stat
}

// Date returns the date of stat.
func (bs *basicStat) Date() time.Time { return bs.date }

// Group returns the group of stat.
func (bs *basicStat) Group() string { return bs.group }

// Value returns the value of stat.
func (bs *basicStat) Value() string { return bs.value }

// Labels returns the labels of stat.
func (bs *basicStat) Labels() field.Labels { return bs.labels }

// Substats returns list of substat.
func (bs *basicStat) Substats() []Stat { return bs.substats }

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
