package github

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pinmonl/pinmonl/monl"
	"github.com/pinmonl/pinmonl/monl/github/api"
)

// Report defines the Github repo info from the rawurl.
type Report struct {
	client     *api.Client
	rawurl     string
	vendorName string

	repoData        *api.Repo
	repoRelPageInfo *api.PageInfo
	repoRels        monl.StatCollection

	latestRel monl.Stat
	totalRel  int
	cursor    int
}

// NewReport creates a report for handling the repo info.
func NewReport(vendorName, rawurl string, httpClient *http.Client) (*Report, error) {
	return &Report{
		client:     api.NewClient(httpClient),
		rawurl:     rawurl,
		vendorName: vendorName,
		repoRels:   make([]monl.Stat, 0),
		totalRel:   0,
		cursor:     -1,
	}, nil
}

// RawURL returns the raw url of the repo.
func (r *Report) RawURL() string { return r.rawurl }

// URI returns the unique name of the repo.
func (r *Report) URI() string { return strings.Join(r.uriPaths(), "/") }

// Vendor returns the vendor name.
func (r *Report) Vendor() string { return r.vendorName }

// Latest returns the latest release stat.
func (r *Report) Latest() monl.Stat { return r.latestRel }

// Length returns the total count of release.
func (r *Report) Length() int { return r.totalRel }

// Popularity returns a list of stats which shows its pupularity.
func (r *Report) Popularity() monl.StatCollection {
	if r.repoData == nil {
		return nil
	}

	var (
		now  = time.Now()
		data = r.repoData

		forkStat      = monl.NewStat(now, "fork", strconv.Itoa(data.ForkCount), nil)
		watcherStat   = monl.NewStat(now, "watcher", strconv.Itoa(data.Watchers.TotalCount), nil)
		starStat      = monl.NewStat(now, "star", strconv.Itoa(data.Stargazers.TotalCount), nil)
		openIssueStat = monl.NewStat(now, "open_issue", strconv.Itoa(data.OpenIssues.TotalCount), nil)
	)
	return []monl.Stat{
		starStat,
		forkStat,
		watcherStat,
		openIssueStat,
	}
}

// Next moves the cursor if not reach the end.
// If the cursor is out of range, it will proceed to previous page
// or return nil if reached the end.
func (r *Report) Next() bool {
	// Return if cursor is within the range of release cache
	if r.cursor+1 < len(r.repoRels) {
		r.cursor = r.cursor + 1
		return true
	}

	pi := r.repoRelPageInfo
	// Check before going to the next page
	if !pi.HasNextPage {
		return false
	}

	// Download the releases
	infos := r.uriPaths()
	po := &api.PageOption{First: 50, After: pi.EndCursor}
	repo, err := r.client.ListRepoReleases(context.Background(), infos[0], infos[1], po)
	if err != nil {
		return false
	}
	r.repoRelPageInfo = repo.Releases.PageInfo
	for _, rel := range repo.Releases.Nodes {
		relStat := r.parseRelease(rel)
		r.repoRels = append(r.repoRels, relStat)
	}

	r.cursor = r.cursor + 1
	return true
}

// Previous moves the cursor if has not reach the end.
func (r *Report) Previous() bool {
	if len(r.repoRels) > 0 && r.cursor-1 >= 0 {
		r.cursor = r.cursor - 1
		return true
	}
	return false
}

// Stat returns the stat at current cursor.
func (r *Report) Stat() monl.Stat {
	return r.repoRels[r.cursor]
}

// Download gets the repo info.
func (r *Report) Download() error {
	infos := r.uriPaths()
	repo, err := r.client.GetRepo(context.Background(), infos[0], infos[1])
	if err != nil {
		return err
	}

	r.repoData = repo
	r.repoRelPageInfo = repo.Releases.PageInfo
	r.totalRel = repo.Releases.TotalCount
	if nodes := repo.Releases.Nodes; len(nodes) > 0 {
		relStat := r.parseRelease(nodes[0])
		r.latestRel = relStat
		r.repoRels = append(r.repoRels, relStat)
	}
	return nil
}

// Derived returns the derived urls in other vendor.
func (r *Report) Derived() map[string]string {
	return nil
}

// Close closes the used resources.
func (r *Report) Close() error {
	return nil
}

func (r *Report) uriPaths() []string {
	url, err := url.Parse(r.RawURL())
	if err != nil {
		return []string{"", ""}
	}
	paths := strings.Split(strings.TrimPrefix(url.Path, "/"), "/")
	if len(paths) < 2 {
		return []string{"", ""}
	}
	return paths[:2]
}

func (*Report) parseRelease(rel api.RepoReleaseNode) monl.Stat {
	return monl.NewStat(rel.CreatedAt, "tag", rel.TagName, nil)
}
