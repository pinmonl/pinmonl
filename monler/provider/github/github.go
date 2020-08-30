package github

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/monler/provider/git"
	"github.com/pinmonl/pinmonl/pkgs/pkgdata"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/sirupsen/logrus"
)

// Github settings.
var (
	DefaultHost = pkgdata.GithubHost
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
	return pkgdata.GithubProvider
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

	logrus.Debugf("github: report created %s", pu)

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

	logrus.Debugf("github: report analyzed %s", r.pu)

	client := &Client{tokens: r.tokens}
	report, err := newReport(r.pu, client, gitReport)
	r.lastReport = report
	return report, err
}

func (r *Repo) Derived() ([]string, error) {
	if r.lastReport == nil {
		if _, err := r.analyze(); err != nil {
			return nil, err
		}
	}

	derived, err := r.gitRepo.Derived()
	if err != nil {
		return nil, err
	}

	if pu, err := r.lastReport.URI(); err == nil {
		derived = append(derived, pkguri.ToURL(pu))
	}

	return derived, nil
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
	repoInfo  *RepositoryResponse
	gitReport provider.Report
}

func newReport(pu *pkguri.PkgURI, client *Client, gitReport provider.Report) (*Report, error) {
	resp, err := client.GetRepository(pu.Namespace(), pu.RepoName())
	if err != nil {
		return nil, err
	}

	now := field.Now()
	stats := []*model.Stat{
		&model.Stat{
			Kind:       model.ForkCountStat,
			Value:      strconv.FormatInt(resp.ForkCount, 10),
			RecordedAt: now,
			IsLatest:   true,
		},
	}

	if resp.Stargazers != nil {
		stats = append(stats, &model.Stat{
			Kind:       model.StarCountStat,
			Value:      strconv.FormatInt(resp.Stargazers.TotalCount, 10),
			RecordedAt: now,
			IsLatest:   true,
		})
	}
	if resp.Watchers != nil {
		stats = append(stats, &model.Stat{
			Kind:       model.WatcherCountStat,
			Value:      strconv.FormatInt(resp.Watchers.TotalCount, 10),
			RecordedAt: now,
			IsLatest:   true,
		})
	}
	if resp.Issues != nil {
		stats = append(stats, &model.Stat{
			Kind:       model.OpenIssueCountStat,
			Value:      strconv.FormatInt(resp.Issues.TotalCount, 10),
			RecordedAt: now,
			IsLatest:   true,
		})
	}
	if resp.PrimaryLanguage != nil {
		stats = append(stats, &model.Stat{
			Kind:       model.LangStat,
			Value:      strings.ToLower(resp.PrimaryLanguage.Name),
			RecordedAt: now,
			IsLatest:   true,
		})
	}
	if resp.LicenseInfo != nil {
		stats = append(stats, &model.Stat{
			Kind:       model.LicenseStat,
			Value:      strings.ToLower(resp.LicenseInfo.Key),
			RecordedAt: now,
			IsLatest:   true,
		})
	}

	if len(resp.FundingLinks) > 0 {
		fundingStats := make(model.StatList, len(resp.FundingLinks))
		for i := range resp.FundingLinks {
			f := resp.FundingLinks[i]
			fundingStats[i] = &model.Stat{
				Name:       f.Platform,
				Value:      f.URL,
				RecordedAt: now,
			}
		}
		stats = append(stats, &model.Stat{
			Kind:       model.FundingStat,
			Value:      strconv.Itoa(len(fundingStats)),
			RecordedAt: now,
			IsLatest:   true,
			Substats:   &fundingStats,
		})
	}

	if gitStats, err := gitReport.Stats(); err == nil {
		stats = append(stats, gitStats...)
	} else {
		return nil, err
	}

	return &Report{
		PkgURI:    pu,
		stats:     stats,
		repoInfo:  resp,
		gitReport: gitReport,
	}, nil
}

func (r *Report) URI() (*pkguri.PkgURI, error) {
	return r.PkgURI, nil
}

func (r *Report) Pkg() (*model.Pkg, error) {
	return nil, nil
}

func (r *Report) Stats() ([]*model.Stat, error) {
	return r.stats, nil
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
