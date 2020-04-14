package npm

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler"
)

// ReportOpts defines options of creating report.
type ReportOpts struct {
	URI    string
	Client *http.Client
}

// Report shows stats and tags from API.
type Report struct {
	uri       string
	client    *apiClient
	stats     monler.StatList
	tags      monler.StatList
	latestTag *monler.Stat
	cursor    int
}

var _ monler.Report = &Report{}

// NewReport creates NPM report.
func NewReport(opts *ReportOpts) (*Report, error) {
	r := &Report{
		uri:    opts.URI,
		client: &apiClient{client: opts.Client},
		cursor: -1,
	}
	return r, nil
}

// URL returns the url to web page.
func (r *Report) URL() string { return DefaultEndpoint + "/package/" + r.uri }

// Provider returns the name of provider.
func (r *Report) Provider() string { return Name }

// ProviderURI returns the unique identifier in provider.
func (r *Report) ProviderURI() string { return r.uri }

// Stats returns the stats of repo, such as star, fork, major language.
func (r *Report) Stats() []*monler.Stat { return r.stats }

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
	if r.cursor < 0 {
		return nil
	}
	if r.cursor >= r.Len() {
		return nil
	}
	return r.tags[r.cursor]
}

// LatestTag returns the latest tag.
func (r *Report) LatestTag() *monler.Stat { return r.latestTag }

// Len returns total count of tags.
func (r *Report) Len() int { return len(r.tags) }

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
	pkgRes, err := r.client.getPackage(r.ProviderURI())
	if err != nil {
		return err
	}
	dlRes, err := r.client.getDownloadPoint(r.ProviderURI(), time.Now().Add(-time.Hour*24*30), time.Now())
	if err != nil {
		return err
	}

	tags, err := parseTags(pkgRes)
	if err != nil {
		return err
	}
	r.tags = tags
	r.latestTag = monler.StatList(tags).LatestStable()

	dStats, err := parseDownloadStats(dlRes)
	if err != nil {
		return err
	}
	r.stats = append(r.stats, dStats...)

	pStats, err := parsePackageStats(pkgRes)
	if err != nil {
		return err
	}
	r.stats = append(r.stats, pStats...)

	return nil
}

func parseTags(res *packageResponse) ([]*monler.Stat, error) {
	var stats []*monler.Stat
	for ver := range res.Versions {
		stats = append(stats, &monler.Stat{
			Kind:       monler.KindTag,
			Value:      ver,
			RecordedAt: field.Time(res.Time[ver]),
		})
	}
	sort.Sort(sort.Reverse(monler.StatByVersion(stats)))
	return stats, nil
}

func parsePackageStats(res *packageResponse) ([]*monler.Stat, error) {
	now := field.Time(time.Now())

	var stats []*monler.Stat
	if len(res.DistTags) > 0 {
		for alias, tag := range res.DistTags {
			stats = append(stats, &monler.Stat{
				Kind:       monler.KindTagAlias,
				Value:      alias,
				RecordedAt: now,
				Substats: []*monler.Stat{
					{Kind: monler.KindTag, Value: tag},
				},
			})
		}
	}
	latestTag := ""
	if lTag, ok := res.DistTags["latest"]; ok {
		latestTag = lTag
	}
	if ver, ok := res.Versions[latestTag]; ok {
		stats = append(stats,
			&monler.Stat{
				Kind:       monler.KindSize,
				Value:      strconv.Itoa(ver.Dist.UnpackedSize),
				RecordedAt: now,
			},
			&monler.Stat{
				Kind:       monler.KindFileCount,
				Value:      strconv.Itoa(ver.Dist.FileCount),
				RecordedAt: now,
			},
		)
	}
	return stats, nil
}

func parseDownloadStats(res *downloadResponse) ([]*monler.Stat, error) {
	sum := 0
	for _, dl := range res.Downloads {
		sum += dl.Downloads
	}
	return []*monler.Stat{
		{
			Kind:       monler.KindDownload,
			Value:      strconv.Itoa(sum),
			RecordedAt: field.Time(res.RequestedAt),
			Labels:     field.Labels{"range": "last_30d"},
		},
	}, nil
}
