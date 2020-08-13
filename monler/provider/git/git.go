package git

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/monler/prvdutils"
	"github.com/pinmonl/pinmonl/pkgs/pkgdata"
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

func (p *Provider) ProviderName() string {
	return pkgdata.GitProvider
}

func (p *Provider) Open(rawurl string) (provider.Repo, error) {
	return p.open(rawurl)
}

func (p *Provider) Parse(uri string) (provider.Repo, error) {
	pu, err := pkguri.NewFromURI(uri)
	if err != nil {
		return nil, err
	}
	if pu.Provider != p.ProviderName() {
		return nil, ErrNotSupportedURI
	}
	return p.open(pkguri.ToURL(pu))
}

func (p *Provider) open(rawurl string) (*Repo, error) {
	return newRepo(rawurl)
}

func (p *Provider) Ping(rawurl string) error {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{rawurl},
	})
	_, err := remote.List(&git.ListOptions{})
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		return err
	}
	return nil
}

type Repo struct {
	gitURL  string
	tempDir string
	repo    *git.Repository
}

func NewRepo(gitURL string) (*Repo, error) {
	return newRepo(gitURL)
}

func newRepo(gitURL string) (*Repo, error) {
	// Git clone to temp directory.
	dir, err := getCloneDir(gitURL)
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:        gitURL,
		NoCheckout: true,
		Progress:   CloneProgress,
	})
	if err == git.ErrRepositoryAlreadyExists {
		repo, err = git.PlainOpen(dir)
		if err == nil {
			repo.Fetch(&git.FetchOptions{
				Progress: CloneProgress,
			})
		}
	}
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		if !IsDev {
			os.RemoveAll(dir)
		}
		return nil, err
	}

	return &Repo{
		gitURL:  gitURL,
		tempDir: dir,
		repo:    repo,
	}, nil
}

func (r *Repo) Analyze() (provider.Report, error) {
	return r.analyze()
}

func (r *Repo) analyze() (*Report, error) {
	// Fill in tags.
	tagIter, err := r.repo.Tags()
	if err != nil {
		return nil, err
	}
	tags := model.StatList{}
	tagIter.ForEach(func(ref *plumbing.Reference) error {
		// Parse annotated tag.
		tag, err := r.repo.TagObject(ref.Hash())
		if err == nil {
			tags = append(tags, &model.Stat{
				RecordedAt: field.Time(tag.Tagger.When),
				Kind:       model.TagStat,
				Value:      tag.Name,
				Checksum:   tag.Hash.String(),
			})
			return nil
		}
		// Parse lightweight tag.
		commit, err := r.repo.CommitObject(ref.Hash())
		if err == nil {
			tags = append(tags, &model.Stat{
				RecordedAt: field.Time(commit.Committer.When),
				Kind:       model.TagStat,
				Value:      strings.TrimPrefix(string(ref.Name()), "refs/tags/"),
				Checksum:   commit.Hash.String(),
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

	// Construct pkguri.
	pu, err := pkguri.ParseGit(r.gitURL)
	if err != nil {
		return nil, err
	}

	return newReport(pu, tags)
}

func (r *Repo) Derived() ([]string, error) {
	derived := make([]string, 0)

	if npmUrls, err := r.GuessNpm(); err == nil {
		derived = append(derived, npmUrls...)
	}

	return derived, nil
}

func (r *Repo) Close() error {
	if !IsDev {
		os.RemoveAll(r.tempDir)
	}
	return nil
}

func (r *Repo) GuessNpm() ([]string, error) {
	packageFile, err := r.file("package.json")
	if err != nil {
		return nil, err
	}

	fr, err := packageFile.Reader()
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	var packageContent struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(fr).Decode(&packageContent); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	if packageContent.Name != "" {
		pu := &pkguri.PkgURI{
			Provider: pkgdata.NpmProvider,
			URI:      packageContent.Name,
			Proto:    pkguri.DefaultProto,
		}
		urls = append(urls, pkguri.ToURL(pu))
	}

	return urls, nil
}

func (r *Repo) file(paths ...string) (file *object.File, err error) {
	ref, err := r.repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		f, e := commit.File(path)
		if e == nil {
			file = f
			return
		}

		err = e
	}
	return
}

// Report inherits StaticReport.
type Report struct {
	*prvdutils.StaticReport
}

func newReport(pu *pkguri.PkgURI, tags []*model.Stat) (*Report, error) {
	return &Report{
		StaticReport: prvdutils.NewStaticReport(pu, nil, tags),
	}, nil
}

var _ provider.Provider = &Provider{}
var _ provider.Repo = &Repo{}
var _ provider.Report = &Report{}
