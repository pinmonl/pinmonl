package git

import (
	"regexp"
	"sort"

	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// ReportOpts defines the options of creating report.
type ReportOpts struct {
	URL string
}

// Report reports stat of Git repository.
type Report struct {
	url       string
	repo      *git.Repository
	tags      []*monler.Stat
	latestTag *monler.Stat
	cursor    int
}

// NewReport creates report.
func NewReport(opts *ReportOpts) (*Report, error) {
	r := &Report{
		url:    opts.URL,
		cursor: -1,
	}
	return r, nil
}

// URL returns the url to web page.
func (r *Report) URL() string { return r.url }

// Provider returns the name of provider.
func (r *Report) Provider() string { return Name }

// ProviderURI returns the unique identifier in provider.
func (r *Report) ProviderURI() string { return r.url }

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
	var derived []monler.Report
	if readme, err := getReadme(r.repo); err == nil {
		if reps, err := ReportFromReadme(mlrepo, readme, cred); err == nil {
			derived = append(derived, reps...)
		}
	}
	return derived, nil
}

// ReportFromReadme aggregates reports with provided monler.Repository.
func ReportFromReadme(mlrepo *monler.Repository, readme string, cred monler.Credential) ([]monler.Report, error) {
	var reps []monler.Report
	prds := mlrepo.Providers()
	for _, url := range extractURLs(readme) {
		for _, prd := range prds {
			if prd == Name {
				continue
			}
			if err := mlrepo.Ping(prd, url, cred); err == nil {
				if rep, err := mlrepo.Open(prd, url, cred); err == nil {
					reps = append(reps, rep)
				}
			}
		}
	}
	return reps, nil
}

func (r *Report) download() error {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: r.URL(),
	})
	if ErrIsEmptyGitRepository(err) {
		return nil
	}
	if err != nil {
		return err
	}
	r.repo = repo

	tags, err := parseTags(repo)
	if err != nil {
		return err
	}
	r.tags = tags
	r.latestTag = monler.StatList(tags).LatestStable()

	return nil
}

func parseTags(repo *git.Repository) ([]*monler.Stat, error) {
	iter, err := repo.TagObjects()
	if err != nil {
		return nil, err
	}
	var tags []*monler.Stat
	err = iter.ForEach(func(tag *object.Tag) error {
		pt, err := parseTag(tag)
		if err != nil {
			return err
		}
		tags = append(tags, pt)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(monler.StatByVersion(tags)))
	return tags, nil
}

func parseTag(tag *object.Tag) (*monler.Stat, error) {
	return &monler.Stat{
		Kind:       monler.KindTag,
		Value:      tag.Name,
		RecordedAt: field.Time(tag.Tagger.When),
	}, nil
}

func getReadme(repo *git.Repository) (string, error) {
	readme, err := getContent(repo, "README.md")
	if err != nil {
		readme, err = getContent(repo, "readme.md")
	}
	return readme, err
}

func extractURLs(content string) []string {
	var urls []string
	keys := make(map[string]bool)
	re := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9]{1,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		u := monler.URLNormalize(m[0])
		if _, ok := keys[u]; !ok {
			urls = append(urls, u)
			keys[u] = true
		}
	}
	return urls
}

func getContent(repo *git.Repository, path string) (string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}
	file, err := commit.File(path)
	return file.Contents()
}
