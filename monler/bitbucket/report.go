package bitbucket

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

// Report stores stats and tags from git.
type Report struct {
	uri       string
	client    *apiClient
	stats     monler.StatList
	gitReport *git.Report
}

// NewReport creates report.
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
	ws, repo := ExtractURI(r.uri)

	fRes, err := r.client.getForks(ws, repo)
	if err != nil {
		return err
	}
	fStat, err := parseForksStat(fRes)
	if err != nil {
		return err
	}
	r.stats = append(r.stats, fStat)

	wRes, err := r.client.getWatchers(ws, repo)
	if err != nil {
		return err
	}
	wStat, err := parseWatchersStat(wRes)
	if err != nil {
		return err
	}
	r.stats = append(r.stats, wStat)

	return nil
}

func parseWatchersStat(res *listResponse) (*monler.Stat, error) {
	return &monler.Stat{
		Kind:       monler.KindWatcher,
		Value:      strconv.Itoa(res.Size),
		RecordedAt: field.Time(time.Now()),
	}, nil
}

func parseForksStat(res *listResponse) (*monler.Stat, error) {
	return &monler.Stat{
		Kind:       monler.KindFork,
		Value:      strconv.Itoa(res.Size),
		RecordedAt: field.Time(time.Now()),
	}, nil
}
