package github

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/monler/provider/git"
	"github.com/pinmonl/pinmonl/monler/provider/git/gitderive"
	"github.com/pinmonl/pinmonl/pkgs/pkgrepo"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/sirupsen/logrus"
)

// Github settings.
var (
	DefaultHost = pkguri.GithubHost
)

// Errors.
var (
	ErrNotSupport = errors.New("github: repo not support")
)

type Provider struct {
	tokens *TokenStore
}

func NewProvider() (*Provider, error) {
	return NewProviderWithTokens(nil)
}

func NewProviderWithTokens(tokens *TokenStore) (*Provider, error) {
	p := Provider{tokens: tokens}
	if p.tokens == nil {
		p.tokens = globalTokens
	}
	return &p, nil
}

func (p *Provider) ProviderName() string {
	return pkguri.GithubProvider
}

func (p *Provider) Open(rawurl string) (provider.Repo, error) {
	pu, err := pkguri.ParseGithub(rawurl)
	if err != nil {
		return nil, err
	}
	return p.open(pu)
}

func (p *Provider) Parse(uri string) (provider.Repo, error) {
	pu, err := pkguri.NewFromURI(uri)
	if err != nil {
		return nil, err
	}
	if pu.Provider != pkguri.GithubProvider {
		return nil, err
	}
	return p.open(pu)
}

func (p *Provider) open(pu *pkguri.PkgURI) (*Repo, error) {
	return newRepo(pu, p.tokens)
}

func (p *Provider) Ping(rawurl string) error {
	pu, err := pkguri.ParseGithub(rawurl)
	if err != nil {
		return err
	}

	resp, err := http.Get(pkguri.ToURL(pu))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return ErrNotSupport
	}
	return nil
}

type Repo struct {
	pu         *pkguri.PkgURI
	tokens     *TokenStore
	gitRepo    *git.Repo
	lastReport *Report
}

func newRepo(pu *pkguri.PkgURI, tokens *TokenStore) (*Repo, error) {
	gitURL := pkguri.ToURL(pu)
	gitRepo, err := git.NewRepo(gitURL)
	if err != nil {
		return nil, err
	}

	return &Repo{
		pu:      pu,
		tokens:  tokens,
		gitRepo: gitRepo,
	}, nil
}

func (r *Repo) Analyze() (provider.Report, error) {
	return r.analyze()
}

func (r *Repo) analyze() (*Report, error) {
	gitReport, err := r.gitRepo.Analyze()
	if err != nil {
		return nil, err
	}

	client := &Client{tokens: r.tokens}
	report, err := newReport(r.pu, client, gitReport)
	r.lastReport = report
	return report, err
}

func (r *Repo) Derived() ([]provider.Report, error) {
	if r.lastReport == nil {
		report, err := r.analyze()
		if err != nil {
			return nil, err
		}
		r.lastReport = report
	}

	derived := gitderive.New(r.gitRepo)
	primLang := r.lastReport.PrimaryLanguage()
	logrus.Debugf("github: repo's primary language is %q", primLang)

	if primLang == pkgrepo.Javascript {
		derived.AddNpm()
	}

	return derived.Reports(), nil
}

func (r *Repo) Skipped() []string {
	return []string{pkguri.GitProvider}
}

func (r *Repo) Close() error {
	if r.gitRepo != nil {
		if err := r.gitRepo.Close(); err != nil {
			return err
		}
	}
	return nil
}

type Report struct {
	*pkguri.PkgURI
	stats     []*model.Stat
	repoInfo  *ApiRepoResponse
	gitReport provider.Report
}

func newReport(pu *pkguri.PkgURI, client *Client, gitReport provider.Report) (*Report, error) {
	res, err := client.GetRepository(pu.Namespace(), pu.RepoName())
	if err != nil {
		return nil, err
	}

	now := field.Now()
	stats := make([]*model.Stat, 0)
	stats = append(stats,
		&model.Stat{
			Kind:       model.ForkStat,
			Value:      strconv.Itoa(res.ForkCount),
			RecordedAt: now,
			IsLatest:   true,
		},
		&model.Stat{
			Kind:       model.StarStat,
			Value:      strconv.Itoa(res.Stargazers.TotalCount),
			RecordedAt: now,
			IsLatest:   true,
		},
		&model.Stat{
			Kind:       model.WatcherStat,
			Value:      strconv.Itoa(res.Watchers.TotalCount),
			RecordedAt: now,
			IsLatest:   true,
		},
		&model.Stat{
			Kind:       model.LangStat,
			Value:      res.PrimaryLanguage.Name.String(),
			RecordedAt: now,
			IsLatest:   true,
		},
		&model.Stat{
			Kind:       model.OpenIssueStat,
			Value:      strconv.Itoa(res.Issues.TotalCount),
			RecordedAt: now,
			IsLatest:   true,
		},
		&model.Stat{
			Kind:       model.LicenseStat,
			Value:      res.LicenseInfo.Key,
			RecordedAt: now,
			IsLatest:   true,
		},
	)

	return &Report{
		PkgURI:    pu,
		stats:     stats,
		repoInfo:  res,
		gitReport: gitReport,
	}, nil
}

func (r *Report) URI() (*pkguri.PkgURI, error) {
	return r.PkgURI, nil
}

func (r *Report) Stats() ([]*model.Stat, error) {
	return r.stats, nil
}

func (r *Report) PrimaryLanguage() pkgrepo.Language {
	return r.repoInfo.PrimaryLanguage.Name.Language()
}

func (r *Report) Next() bool {
	return r.gitReport.Next()
}

func (r *Report) Tag() (*model.Stat, error) {
	return r.gitReport.Tag()
}

func (r *Report) Close() error {
	if err := r.gitReport.Close(); err != nil {
		return err
	}
	return nil
}

var _ provider.Provider = &Provider{}
var _ provider.Repo = &Repo{}
var _ provider.Report = &Report{}
