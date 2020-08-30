package model

import (
	"github.com/Masterminds/semver/v3"
	"github.com/pinmonl/pinmonl/model/field"
)

type Stat struct {
	ID          string        `json:"id"`
	PkgID       string        `json:"pkgId"`
	ParentID    string        `json:"parentId"`
	RecordedAt  field.Time    `json:"recordedAt"`
	Kind        StatKind      `json:"kind"`
	Name        string        `json:"name"`
	Value       string        `json:"value"`
	ValueType   StatValueType `json:"valueType"`
	Checksum    string        `json:"checksum"`
	Weight      int           `json:"weight"`
	IsLatest    bool          `json:"isLatest"`
	HasChildren bool          `json:"hasChildren"`

	Substats   *StatList `json:"substats,omitempty"`
	SubstatIDs *[]string `json:"substatIds,omitempty"`
}

func (s Stat) MorphKey() string  { return s.ID }
func (s Stat) MorphName() string { return "stat" }

type StatKind string

const (
	AnyStat             = StatKind("")
	AliasStat           = StatKind("alias")
	ChannelStat         = StatKind("channel")
	DownloadCountStat   = StatKind("download_count")
	FileCountStat       = StatKind("file_count")
	ForkCountStat       = StatKind("fork_count")
	FundingStat         = StatKind("funding")
	LangStat            = StatKind("lang")
	LicenseStat         = StatKind("license")
	ManifestStat        = StatKind("manifest")
	OpenIssueCountStat  = StatKind("open_issue_count")
	PullCountStat       = StatKind("pull_count")
	SizeStat            = StatKind("size")
	StarCountStat       = StatKind("star_count")
	StatusStat          = StatKind("status")
	SubscriberCountStat = StatKind("subscriber_count")
	TagStat             = StatKind("tag")
	VideoCountStat      = StatKind("video_count")
	VideoStat           = StatKind("video")
	ViewCountStat       = StatKind("view_count")
	WatcherCountStat    = StatKind("watcher_count")
)

var ReleaseStatKinds = []StatKind{
	TagStat,
	VideoStat,
}

func IsReleaseStatKind(kind StatKind) bool {
	for _, k := range ReleaseStatKinds {
		if k == kind {
			return true
		}
	}
	return false
}

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

func (sl StatList) GetHasChildren() StatList {
	list := make([]*Stat, 0)
	for _, s := range sl {
		if s.HasChildren {
			list = append(list, s)
		}
	}
	return list
}

func (sl StatList) MustSemver() StatList {
	list := make([]*Stat, 0)
	for _, s := range sl {
		if _, err := semver.NewVersion(s.Value); err == nil {
			list = append(list, s)
		}
	}
	return list
}

func (sl StatList) Contains(val *Stat) bool {
	for _, s := range sl {
		if s == val {
			return true
		}
	}
	return false
}

type StatBySemver StatList

func (sl StatBySemver) Len() int { return len(sl) }

func (sl StatBySemver) Swap(i, j int) { sl[i], sl[j] = sl[j], sl[i] }

func (sl StatBySemver) Less(i, j int) bool {
	iv, erri := semver.NewVersion(sl[i].Value)
	jv, errj := semver.NewVersion(sl[j].Value)
	if erri != nil && errj != nil {
		it := sl[i].RecordedAt.Time()
		jt := sl[j].RecordedAt.Time()
		return it.Before(jt)
	}
	if erri != nil {
		// If error occurs, sort to top.
		return true
	}
	if errj != nil {
		// If error occurs, sort to top.
		return false
	}
	return iv.Compare(jv) < 0
}

type StatByRecordedAt StatList

func (sl StatByRecordedAt) Len() int { return len(sl) }

func (sl StatByRecordedAt) Swap(i, j int) { sl[i], sl[j] = sl[j], sl[i] }

func (sl StatByRecordedAt) Less(i, j int) bool {
	ir := sl[i].RecordedAt.Time()
	jr := sl[j].RecordedAt.Time()
	return ir.Before(jr)
}
