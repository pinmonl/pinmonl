package git

import (
	"bytes"
	"regexp"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pkg/payload"
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
	if npmURL, err := guessNpmURI(r.repo); err == nil {
		report, err := mlrepo.Open("npm", npmURL.String(), cred)
		if err == nil {
			derived = append(derived, report)
		}
	}
	if readme, err := getReadme(r.repo); err == nil {
		reports, err := ReportFromReadme(mlrepo, readme, cred)
		if err == nil {
		merge:
			for _, report := range reports {
				for _, dr := range derived {
					if report.Provider() == dr.Provider() && report.ProviderURI() == dr.ProviderURI() {
						continue merge
					}
				}
				derived = append(derived, report)
			}
		}
	}
	return derived, nil
}

// ReportFromReadme aggregates reports with provided monler.Repository.
func ReportFromReadme(mlrepo *monler.Repository, readme string, cred monler.Credential) ([]monler.Report, error) {
	var reps []monler.Report
	prds := []string{"docker"}
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
	// Sorts tags by date if tags are not using semver.
	if r.latestTag == nil && len(tags) > 0 {
		sort.Sort(sort.Reverse(monler.StatByDate(tags)))
		r.tags = tags
		r.latestTag = tags[0]
	}

	return nil
}

func parseTags(repo *git.Repository) ([]*monler.Stat, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	var tags []*monler.Stat
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tag, err := parseTag(repo, ref)
		if err != nil {
			return err
		}
		tags = append(tags, tag)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(monler.StatByVersion(tags)))
	return tags, nil
}

func parseTag(repo *git.Repository, ref *plumbing.Reference) (*monler.Stat, error) {
	// Parse for annotated tag.
	tag, err := repo.TagObject(ref.Hash())
	if err == nil {
		return &monler.Stat{
			Kind:       monler.KindTag,
			Value:      tag.Name,
			RecordedAt: field.Time(tag.Tagger.When),
		}, nil
	}

	// Parse for lightweight tag.
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	return &monler.Stat{
		Kind:       monler.KindTag,
		Value:      strings.TrimPrefix(string(ref.Name()), "refs/tags/"),
		RecordedAt: field.Time(commit.Committer.When),
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
	if err != nil {
		return "", err
	}
	return file.Contents()
}

func guessNpmURI(repo *git.Repository) (*monler.URL, error) {
	json, err := getContent(repo, "package.json")
	if err != nil {
		return nil, err
	}
	var pkg struct {
		Name string
	}
	err = payload.UnmarshalJSON(bytes.NewBufferString(json), &pkg)
	if err != nil {
		return nil, err
	}
	return monler.NewURLFromRaw("https://npmjs.com/package/"+pkg.Name, pkg.Name)
}
