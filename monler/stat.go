package monler

import (
	"github.com/Masterminds/semver"
	"github.com/pinmonl/pinmonl/model/field"
)

// Stat stores stat data.
type Stat struct {
	RecordedAt field.Time
	Kind       StatKind
	Value      string
	Digest     string
	Labels     field.Labels
	Substats   StatList
}

// StatList is slice of stat.
type StatList []*Stat

// LatestStable finds the latest stable version.
func (sl StatList) LatestStable() *Stat {
	var (
		lVer   *semver.Version
		latest *Stat
	)
	for _, s := range sl {
		sVer, err := semver.NewVersion(s.Value)
		if err != nil {
			continue
		}
		if sVer.Prerelease() != "" {
			continue
		}
		if lVer == nil || sVer.Compare(lVer) > 0 {
			lVer = sVer
			latest = s
		}
	}
	return latest
}

// StatKind represents kind of stat.
type StatKind string

// Stat kinds.
const (
	KindTag         StatKind = StatKind("tag")
	KindFork                 = StatKind("fork")
	KindStar                 = StatKind("star")
	KindWatcher              = StatKind("watcher")
	KindOpenIssue            = StatKind("open_issue")
	KindLang                 = StatKind("lang")
	KindDownload             = StatKind("download")
	KindPull                 = StatKind("pull")
	KindTagAlias             = StatKind("tag_alias")
	KindSize                 = StatKind("size")
	KindFileCount            = StatKind("file_count")
	KindLastUpdated          = StatKind("last_updated")
	KindImage                = StatKind("image")
)

// String returns plain string.
func (s StatKind) String() string { return string(s) }

// StatByVersion sorts value as version.
type StatByVersion StatList

// Len returns the length of slice.
func (s StatByVersion) Len() int { return len(s) }

// Less compares value of item i and j.
func (s StatByVersion) Less(i, j int) bool {
	v1, e1 := semver.NewVersion(s[i].Value)
	v2, e2 := semver.NewVersion(s[j].Value)
	if e1 != nil || e2 != nil {
		return false
	}
	return v1.Compare(v2) < 0
}

// Swap swaps item i and j.
func (s StatByVersion) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// StatByDate sorts item by RecordedAt.
type StatByDate StatList

// Len returns the length of slice.
func (s StatByDate) Len() int { return len(s) }

// Less compares value of item i and j.
func (s StatByDate) Less(i, j int) bool {
	ti, tj := s[i].RecordedAt.Time(), s[j].RecordedAt.Time()
	return ti.Before(tj)
}

// Swap swaps item i and j.
func (s StatByDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
