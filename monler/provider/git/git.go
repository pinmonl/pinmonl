package git

import (
	"errors"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
)

// Errors.
var (
	ErrNotSupportedURI = errors.New("git: uri is not supported")
	ErrNoPing          = errors.New("git: ping is not supported")
)

type Provider struct{}

func NewProvider() (*Provider, error) {
	return &Provider{}, nil
}

func (p *Provider) Name() string {
	return provider.Git
}

func (p *Provider) Open(rawurl string) (provider.Repo, error) {
	return p.open(rawurl)
}

func (p *Provider) Parse(uri string) (provider.Repo, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if u.Scheme != p.Name() {
		return nil, ErrNotSupportedURI
	}
	proto := "https"
	return p.open(proto + "://" + u.Host + "/" + strings.Trim(u.Path, "/"))
}

func (p *Provider) open(rawurl string) (provider.Repo, error) {
	return newRepo(rawurl)
}

func (p *Provider) Ping(_ string) error {
	return ErrNoPing
}

type Repo struct {
	*sync.Mutex
	gitURL string
}

func newRepo(gitURL string) (*Repo, error) {
	return &Repo{
		Mutex:  &sync.Mutex{},
		gitURL: gitURL,
	}, nil
}

func (r *Repo) Analyze() (provider.Report, error) {
	pu := &pkguri.PkgURI{
		Provider: provider.Git,
		URI:      r.gitURL,
	}
	return newReport(pu)
}

func (r *Repo) Derived() ([]provider.Report, error) {
	return nil, nil
}

func (r *Repo) Close() error {
	return nil
}

type Report struct {
	*pkguri.PkgURI
	*sync.Mutex
	cursor  int
	repo    *git.Repository
	tempDir string
	tags    model.StatList
}

func newReport(uri *pkguri.PkgURI) (*Report, error) {
	// Prepare temp directory.
	dir, err := getTempDir()
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:        uri.URI,
		NoCheckout: true,
		Progress:   CloneProgress,
	})
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		os.RemoveAll(dir)
		return nil, err
	}

	// Fill in tags.
	tagIter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	tags := make([]*model.Stat, 0)
	tagIter.ForEach(func(ref *plumbing.Reference) error {
		// Parse annotated tag.
		tag, err := repo.TagObject(ref.Hash())
		if err == nil {
			tags = append(tags, &model.Stat{
				RecordedAt: field.Time(tag.Tagger.When),
				Kind:       model.TagStat,
				Value:      tag.Name,
			})
			return nil
		}
		// Parse lightweight tag.
		commit, err := repo.CommitObject(ref.Hash())
		if err == nil {
			tags = append(tags, &model.Stat{
				RecordedAt: field.Time(commit.Committer.When),
				Kind:       model.TagStat,
				Value:      strings.TrimPrefix(string(ref.Name()), "refs/tags/"),
			})
			return nil
		}
		return nil
	})
	// Sort by semver.
	if len(tags) > 0 {
		sort.Sort(model.StatBySemver(tags))
		tags[len(tags)-1].IsLatest = true
	}

	// Create report.
	r := &Report{
		Mutex:   &sync.Mutex{},
		cursor:  -1,
		repo:    repo,
		tempDir: dir,
		tags:    tags,
	}
	return r, nil
}

func (r *Report) URI() (*pkguri.PkgURI, error) {
	return r.PkgURI, nil
}

func (r *Report) Stats() ([]*model.Stat, error) {
	return nil, nil
}

func (r *Report) Next() bool {
	r.Lock()
	defer r.Unlock()
	if r.cursor+1 < len(r.tags) {
		r.cursor++
		return true
	}
	return false
}

func (r *Report) Tag() (*model.Stat, error) {
	r.Lock()
	defer r.Unlock()
	if r.cursor < 0 {
		return nil, errors.New("git: out of range")
	}
	if r.cursor >= len(r.tags) {
		return nil, io.EOF
	}
	return r.tags[r.cursor], nil
}

func (r *Report) Close() error {
	if r.tempDir != "" {
		if err := os.RemoveAll(r.tempDir); err != nil {
			return err
		}
	}
	return nil
}

var _ provider.Provider = &Provider{}
var _ provider.Repo = &Repo{}
var _ provider.Report = &Report{}
