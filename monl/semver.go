package monl

import (
	"github.com/Masterminds/semver"
)

type BySemver []Stat

func (b BySemver) Len() int      { return len(b) }
func (b BySemver) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b BySemver) Less(i, j int) bool {
	vi, err := semver.NewVersion(b[i].Value())
	if err != nil {
		return false
	}
	vj, err := semver.NewVersion(b[j].Value())
	if err != nil {
		return true
	}
	return vi.Compare(vj) > 0
}

type SemverTags struct {
	Stable  []Stat
	Next    []Stat
	Unknown []Stat
}

func NewSemverTags(stats []Stat, latestStable Stat) (SemverTags, error) {
	var vl *semver.Version
	if latestStable != nil {
		vl, _ = semver.NewVersion(latestStable.Value())
	}

	out := SemverTags{}
	for _, stat := range stats {
		v, err := semver.NewVersion(stat.Value())
		switch {
		case err != nil:
			out.Unknown = append(out.Unknown, stat)
		case vl != nil && v.Compare(vl) > 0:
			out.Next = append(out.Next, stat)
		default:
			out.Stable = append(out.Stable, stat)
		}
	}
	return out, nil
}
