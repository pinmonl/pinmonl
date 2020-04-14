package helm

import (
	"net/http"

	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler"
)

// ReportOpts defines the options of creating report.
type ReportOpts struct {
	URI    string
	Client *http.Client
}

// Report reports stats of helm repo.
type Report struct {
	uri       string
	client    *apiClient
	tags      []*monler.Stat
	latestTag *monler.Stat
	cursor    int
}

var _ monler.Report = &Report{}

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
func (r *Report) URL() string { return DefaultEndpoint + "/charts/" + r.uri }

// Provider returns the name of provider.
func (r *Report) Provider() string { return Name }

// ProviderURI returns the unique identifier in provider.
func (r *Report) ProviderURI() string { return r.uri }

// Stats returns the stats of repo, such as star, fork, major language.
func (r *Report) Stats() []*monler.Stat {
	return nil
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
	if r.cursor >= 0 && r.cursor < r.Len() {
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
	vers, err := r.client.listVersions(r.ProviderURI())
	if err != nil {
		return err
	}

	tags, err := parseTags(vers)
	if err != nil {
		return err
	}
	r.tags = tags
	r.latestTag = monler.StatList(tags).LatestStable()
	return nil
}

func parseTags(res []*ChartVersionResponse) ([]*monler.Stat, error) {
	var tags []*monler.Stat
	for _, vres := range res {
		at := field.Time(vres.Attributes.Created)
		tags = append(tags, &monler.Stat{
			Kind:       monler.KindTag,
			Value:      vres.Attributes.Version,
			Digest:     vres.Attributes.Digest,
			RecordedAt: at,
		})
	}
	return tags, nil
}
