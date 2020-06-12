package git

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/monler/prvdutils"
	"github.com/pinmonl/pinmonl/pkgs/pkgrepo"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/sirupsen/logrus"
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
	return pkguri.GitProvider
}

func (p *Provider) Open(rawurl string) (provider.Repo, error) {
	return p.open(rawurl)
}

func (p *Provider) Parse(uri string) (provider.Repo, error) {
	pu, err := pkguri.Parse(uri)
	if err != nil {
		return nil, err
	}
	if pu.Provider != p.ProviderName() {
		return nil, ErrNotSupportedURI
	}
	return p.open(pu.URL().String())
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
	tags := make([]*model.Stat, 0)
	tagIter.ForEach(func(ref *plumbing.Reference) error {
		// Parse annotated tag.
		tag, err := r.repo.TagObject(ref.Hash())
		if err == nil {
			tags = append(tags, &model.Stat{
				RecordedAt: field.Time(tag.Tagger.When),
				Kind:       model.TagStat,
				Value:      tag.Name,
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
	pu, err := pkguri.ParseFromGit(r.gitURL)
	if err != nil {
		return nil, err
	}

	return newReport(pu, tags)
}

func (r *Repo) Derived() ([]provider.Report, error) {
	return nil, nil
}

func (r *Repo) Skipped() []string {
	return nil
}

func (r *Repo) Close() error {
	return nil
}

func (r *Repo) openFile(file string) (*os.File, error) {
	return os.Open(filepath.Join(r.tempDir, file))
}

func (r *Repo) openFilePath(file pkgrepo.FilePath) (*os.File, error) {
	return r.openFile(file.String())
}

func (r *Repo) GetNpmURI() (*pkguri.PkgURI, error) {
	f, err := r.openFilePath(pkgrepo.NpmPackage)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var pkg struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(f).Decode(&pkg)
	if err != nil {
		logrus.Debugf("git: GetNpmURI parse err(%v)", err)
		return nil, err
	}
	if pkg.Name == "" {
		return nil, errors.New("git: npm package name is not valid")
	}

	return &pkguri.PkgURI{
		Provider: pkguri.NpmProvider,
		URI:      pkg.Name,
	}, nil
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
