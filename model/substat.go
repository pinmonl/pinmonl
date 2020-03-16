package model

import "github.com/pinmonl/pinmonl/model/field"

// Substat stores extra information for Stat.
type Substat struct {
	ID     string       `json:"id"     db:"substat_id"`
	StatID string       `json:"statId" db:"substat_stat_id"`
	Kind   string       `json:"kind"   db:"substat_kind"`
	Labels field.Labels `json:"labels" db:"substat_labels"`
}

// SubstatList is slice of Substat.
type SubstatList []Substat

// FindStat filters Variants by Stat.
func (sl SubstatList) FindStat(t Stat) SubstatList {
	var out []Substat
	for _, s := range sl {
		if s.StatID == t.ID {
			out = append(out, s)
		}
	}
	return out
}
