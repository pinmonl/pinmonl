package gitlab

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/git"
)

// ReportOpts defines the options of creating report.
type ReportOpts struct {
	URI    string
	Client *http.Client
}

// Report stores stats from API and tags from git.
type Report struct {
	client    *apiClient
	uri       string
	stats     monler.StatList
	gitReport *git.Report
}

var _ monler.Report = &Report{}

// NewReport creates Gitlab report.
func NewReport(opts *ReportOpts) (*Report, error) {
	r := &Report{
		uri:    opts.URI,
		client: &apiClient{client: opts.Client},
	}

	if gitRep, err := git.NewReport(&git.ReportOpts{
		URL: r.URL(),
	}); err == nil {
		r.gitReport = gitRep
	} else {
		return nil, err
	}

	return r, nil
}

// URL returns the url to web page.
func (r *Report) URL() string { return DefaultEndpoint + "/" + r.uri }

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
	return r.gitReport.Next()
}

// Prev checks whether has previous tag and moves the cursor.
func (r *Report) Prev() bool {
	return r.gitReport.Prev()
}

// Tag returns tag at the current cursor.
func (r *Report) Tag() *monler.Stat {
	return r.gitReport.Tag()
}

// LatestTag returns the latest tag.
func (r *Report) LatestTag() *monler.Stat {
	return r.gitReport.LatestTag()
}

// Len returns total count of tags.
func (r *Report) Len() int {
	return r.gitReport.Len()
}

// Download downloads stats and tags.
func (r *Report) Download() error {
	if err := r.gitReport.Download(); err != nil {
		return err
	}
	return r.download()
}

// Close closes the report.
func (r *Report) Close() error {
	return r.gitReport.Close()
}

// Derived returns related reports.
func (r *Report) Derived(mlrepo *monler.Repository, cred monler.Credential) ([]monler.Report, error) {
	return r.gitReport.Derived(mlrepo, cred)
}

func (r *Report) download() error {
	res, err := r.client.getProject(r.ProviderURI())
	if err != nil {
		return err
	}
	if stats, err := parseStats(res); err == nil {
		r.stats = stats
	} else {
		return err
	}
	return nil
}

func parseStats(res *projectResponse) ([]*monler.Stat, error) {
	now := field.Time(time.Now())

	return []*monler.Stat{
		{
			Kind:       monler.KindStar,
			Value:      strconv.Itoa(res.StarCount),
			RecordedAt: now,
		},
		{
			Kind:       monler.KindFork,
			Value:      strconv.Itoa(res.ForksCount),
			RecordedAt: now,
		},
	}, nil
}
