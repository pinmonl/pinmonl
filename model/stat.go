package model

import (
	"github.com/pinmonl/pinmonl/model/field"
)

// Stat stores the value at the given date and time.
type Stat struct {
	ID         string       `json:"id"         db:"id"`
	MonlID     string       `json:"monlId"     db:"monl_id"`
	RecordedAt field.Time   `json:"recordedAt" db:"recorded_at"`
	Kind       string       `json:"kind"       db:"kind"`
	Value      string       `json:"value"      db:"value"`
	IsLatest   bool         `json:"isLatest"   db:"is_latest"`
	Manifest   field.Labels `json:"manifest"   db:"manifest"`
}

// StatList is slice of Stat.
type StatList []Stat

// FindKind filters stats by kind.
func (sl StatList) FindKind(kind string) StatList {
	var out []Stat
	for _, s := range sl {
		if kind == s.Kind {
			out = append(out, s)
		}
	}
	return out
}

// FindValue filters stats by value.
func (sl StatList) FindValue(value string) StatList {
	var out []Stat
	for _, s := range sl {
		if value == s.Value {
			out = append(out, s)
		}
	}
	return out
}
