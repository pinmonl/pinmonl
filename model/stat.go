package model

import (
	"github.com/Masterminds/semver"
	"github.com/pinmonl/pinmonl/model/field"
)

type Stat struct {
	ID          string        `json:"id"`
	PkgID       string        `json:"pkgId"`
	ParentID    string        `json:"parentId"`
	RecordedAt  field.Time    `json:"recordedAt"`
	Kind        StatKind      `json:"kind"`
	Value       string        `json:"value"`
	ValueType   StatValueType `json:"valueType"`
	Checksum    string        `json:"checksum"`
	Weight      int           `json:"weight"`
	IsLatest    bool          `json:"isLatest"`
	HasChildren bool          `json:"hasChildren"`

	Substats *StatList `json:"substats,omitempty"`
}

func (s Stat) MorphKey() string  { return s.ID }
func (s Stat) MorphName() string { return "stat" }

type StatKind string

const (
	AnyStat       = StatKind("")
	TagStat       = StatKind("tag")
	AliasStat     = StatKind("alias")
	StarStat      = StatKind("star")
	ForkStat      = StatKind("fork")
	OpenIssueStat = StatKind("open_issue")
	LangStat      = StatKind("lang")
	FileCountStat = StatKind("file_count")
	DownloadStat  = StatKind("download")
	PullStat      = StatKind("pull")
	WatcherStat   = StatKind("watcher")
	StatusStat    = StatKind("status")
	LicenseStat   = StatKind("license")
)

type StatValueType int

const (
	StringStat StatValueType = iota
	IntegerStat
)

type StatList []*Stat

func (sl StatList) Keys() []string {
	keys := make([]string, len(sl))
	for i := range sl {
		keys[i] = sl[i].ID
	}
	return keys
}

func (sl StatList) GetKind(k StatKind) StatList {
	list := make([]*Stat, 0)
	for _, s := range sl {
		if s.Kind == k {
			list = append(list, s)
		}
	}
	return list
}

func (sl StatList) GetLatest() StatList {
	list := make([]*Stat, 0)
	for _, s := range sl {
		if s.IsLatest {
			list = append(list, s)
		}
	}
	return list
}

func (sl StatList) GetPkgID(pkgID string) StatList {
	list := make([]*Stat, 0)
	for _, s := range sl {
		if s.PkgID == pkgID {
			list = append(list, s)
		}
	}
	return list
}

func (sl StatList) GetParentID(parentID string) StatList {
	list := make([]*Stat, 0)
	for _, s := range sl {
		if s.ParentID == parentID {
			list = append(list, s)
		}
	}
	return list
}

func (sl StatList) GetValue(value string) StatList {
	list := make([]*Stat, 0)
	for _, s := range sl {
		if s.Value == value {
			list = append(list, s)
		}
	}
	return list
}

type StatBySemver StatList

func (sl StatBySemver) Len() int { return len(sl) }

func (sl StatBySemver) Swap(i, j int) { sl[i], sl[j] = sl[j], sl[i] }

func (sl StatBySemver) Less(i, j int) bool {
	iv, err := semver.NewVersion(sl[i].Value)
	if err != nil {
		return false
	}
	ij, err := semver.NewVersion(sl[j].Value)
	if err != nil {
		return false
	}
	return iv.Compare(ij) < 0
}
