package docker

import (
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler"
)

// ReportOpts defines the options for creating report.
type ReportOpts struct {
	URI    string
	Client *http.Client
}

// Report reports stats of Docker repo.
type Report struct {
	uri       string
	client    *apiClient
	stats     []*monler.Stat
	tags      []*monler.Stat
	latestTag *monler.Stat
	tagLen    int
	cursor    int
}

// NewReport creates report.
func NewReport(opts *ReportOpts) (*Report, error) {
	r := &Report{
		uri:    opts.URI,
		client: &apiClient{client: opts.Client},
		cursor: -1,
	}
	return r, nil
}

// URL returns the url to web page.
func (r *Report) URL() string {
	if strings.HasPrefix(r.uri, OfficialNamespace) {
		return DefaultEndpoint + "/" + OfficialURLPrefix + "/" + strings.TrimPrefix(r.uri, OfficialNamespace+"/")
	}
	return DefaultEndpoint + "/" + ThirdPartyURLPrefix + "/" + r.uri
}

// Provider returns the name of provider.
func (r *Report) Provider() string { return Name }

// ProviderURI returns the unique identifier in provider.
func (r *Report) ProviderURI() string { return r.uri }

// Stats returns the stats of repo, such as star, fork, major language.
func (r *Report) Stats() []*monler.Stat {
	return r.stats
}

// Next checks whether has next tag and moves the cursor.
func (r *Report) Next() bool {
	if r.cursor+1 < r.Len() {
		r.cursor++
		return true
	}
	return false
}

// Prev checks whether has previous tag and moves the cursor.
func (r *Report) Prev() bool {
	if r.cursor-1 >= 0 {
		r.cursor--
		return true
	}
	return false
}

// Tag returns tag at the current cursor.
func (r *Report) Tag() *monler.Stat {
	if r.cursor >= 0 || r.cursor < r.Len() {
		return r.tags[r.cursor]
	}
	return nil
}

// LatestTag returns the latest tag.
func (r *Report) LatestTag() *monler.Stat {
	return r.latestTag
}

// Len returns total count of tags.
func (r *Report) Len() int {
	return len(r.tags)
}

// Download downloads stats and tags.
func (r *Report) Download() error {
	return r.download()
}

// Close closes the report.
func (r *Report) Close() error {
	return nil
}

// Derived returns related reports.
func (r *Report) Derived(mlrepo *monler.Repository, cred monler.Credential) ([]monler.Report, error) {
	return nil, nil
}

func (r *Report) download() error {
	rRes, err := r.client.getRepository(r.ProviderURI())
	if err != nil {
		return err
	}
	rStats, err := parseRepositoryStats(rRes)
	if err != nil {
		return err
	}
	r.stats = rStats

	tRes, err := r.client.listAllTags(r.ProviderURI())
	if err != nil {
		return err
	}
	tags, err := parseTagStats(tRes)
	if err != nil {
		return err
	}
	r.tags = tags
	r.latestTag = findLatestTag(tags)

	return nil
}

func parseRepositoryStats(res *repositoryResponse) ([]*monler.Stat, error) {
	now := field.Time(time.Now())

	return []*monler.Stat{
		{
			Kind:       monler.KindStar,
			Value:      strconv.Itoa(res.StarCount),
			RecordedAt: now,
		},
		{
			Kind:       monler.KindPull,
			Value:      strconv.Itoa(res.PullCount),
			RecordedAt: now,
		},
	}, nil
}

func parseTagStats(res []tagResponse) ([]*monler.Stat, error) {
	var (
		stats   = make([]*monler.Stat, 0)
		aliases = make(map[string][]string)
	)
	// Parse response to stat.
	for _, tr := range res {
		s, err := parseTagStat(tr)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	// Scan for alias.
	for i, a := range stats {
		for _, b := range stats[i+1:] {
			// Skip if no intersection.
			if !tagIntersected(a, b) {
				continue
			}

			var parent, child string
			if tagIsSubset(a, b) {
				parent, child = b.Value, a.Value
			} else {
				parent, child = a.Value, b.Value
			}

			aliases[parent] = append(aliases[parent], child)
		}
	}
	// Appends alias to stat.
	for _, stat := range stats {
		if children, ok := aliases[stat.Value]; ok {
			for _, child := range children {
				stat.Substats = append(stat.Substats, &monler.Stat{
					Kind:  monler.KindTagAlias,
					Value: child,
				})
			}
		}
	}
	return stats, nil
}

func parseTagStat(res tagResponse) (*monler.Stat, error) {
	s := &monler.Stat{
		Kind:       monler.KindTag,
		Value:      res.Name,
		RecordedAt: field.Time(res.LastUpdated),
	}
	for _, im := range res.Images {
		s.Substats = append(s.Substats, &monler.Stat{
			Kind:   monler.KindImage,
			Value:  im.Digest,
			Digest: im.Digest,
			Labels: field.Labels{
				"architecture": im.Architecture,
				"features":     im.Features,
				"variant":      im.Variant,
				"os":           im.Os,
				"os_features":  im.OsFeatures,
				"os_version":   im.OsVersion,
				"size":         strconv.Itoa(im.Size),
			},
		})
	}
	return s, nil
}

func tagIntersected(a, b *monler.Stat) bool {
	checks := make(map[string]int)
	for _, s := range a.Substats {
		if s.Kind == monler.KindImage {
			checks[s.Digest]++
		}
	}
	for _, s := range b.Substats {
		if s.Kind != monler.KindImage {
			continue
		}
		if c, ok := checks[s.Digest]; ok && c > 0 {
			return true
		}
	}
	return false
}

func tagIsSubset(a, b *monler.Stat) bool {
	checks := make(map[string]int)
	for _, s := range b.Substats {
		if s.Kind == monler.KindImage {
			checks[s.Digest]++
		}
	}
	for _, s := range a.Substats {
		if s.Kind != monler.KindImage {
			continue
		}
		if _, ok := checks[s.Digest]; !ok {
			return false
		}
	}
	return len(a.Value) > len(b.Value)
}

var versionRegex = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?`)

func findLatestTag(stats []*monler.Stat) *monler.Stat {
	sort.Sort(sort.Reverse(monler.StatByDate(stats)))
	var (
		lVer   *semver.Version
		latest *monler.Stat
	)
	for _, s := range stats {
		switch {
		case latest == nil && s.Value == "latest":
			latest = s
		case versionRegex.MatchString(s.Value):
			m := versionRegex.FindString(s.Value)
			sVer, err := semver.NewVersion(m)
			if err != nil {
				continue
			}
			if lVer == nil {
				lVer = sVer
				latest = s
			} else if sVer.Compare(lVer) > 0 && len(sVer.Original()) > len(lVer.Original()) {
				lVer = sVer
				latest = s
			}
		}
	}
	return latest
}
