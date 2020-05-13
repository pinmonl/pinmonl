package github

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/git"
	"github.com/pinmonl/pinmonl/monler/helm"
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
	gitReport *git.Report
	stats     monler.StatList
}

var _ monler.Report = &Report{}

// NewReport creates Github report.
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
func (r *Report) URL() string {
	return DefaultEndpoint + "/" + r.uri
}

// Provider returns the name of provider.
func (r *Report) Provider() string { return Name }

// ProviderURI returns the unique identifier in provider.
func (r *Report) ProviderURI() string {
	return r.uri
}

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
	reports, err := r.gitReport.Derived(mlrepo, cred)
	if err != nil {
		return nil, err
	}

	// Test whether helm provider is registered.
	if _, err := mlrepo.Get(helm.Name); err != nil {
		logx.Debugf("github monler: err %v", err)
		return reports, nil
	}
	charts, err := helm.Search(r.URL())
	if err != nil {
		return reports, nil
	}
	if len(charts) == 0 {
		return reports, nil
	}
	// Open report for each helm chart.
	for _, chart := range charts {
		report, err := mlrepo.Open(helm.Name, chart.ID, cred)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func (r *Report) download() error {
	owner, repo := ExtractURI(r.ProviderURI())
	res, err := r.client.getRepo(owner, repo)
	if err != nil {
		return err
	}
	stats, err := parseStats(res)
	if err != nil {
		return err
	}
	r.stats = stats
	return nil
}

func parseStats(res *packageResponse) (monler.StatList, error) {
	now := field.Time(time.Now())

	return monler.StatList{
		&monler.Stat{
			Kind:       monler.KindFork,
			Value:      strconv.Itoa(res.ForksCount),
			RecordedAt: now,
		},
		&monler.Stat{
			Kind:       monler.KindOpenIssue,
			Value:      strconv.Itoa(res.OpenIssuesCount),
			RecordedAt: now,
		},
		// &monler.Stat{
		// 	Kind:       monler.KindWatcher,
		// 	Value:      strconv.Itoa(res.WatchersCount),
		// 	RecordedAt: now,
		// },
		&monler.Stat{
			Kind:       monler.KindStar,
			Value:      strconv.Itoa(res.StargazersCount),
			RecordedAt: now,
		},
	}, nil
}
