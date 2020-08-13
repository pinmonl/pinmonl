package npm

import (
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/monler/prvdutils"
	"github.com/pinmonl/pinmonl/pkgs/pkgdata"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/sirupsen/logrus"
)

type Provider struct{}

func NewProvider() (*Provider, error) {
	return &Provider{}, nil
}

func (p *Provider) ProviderName() string {
	return pkgdata.NpmProvider
}

func (p *Provider) Open(rawurl string) (provider.Repo, error) {
	pu, err := pkguri.ParseNpm(rawurl)
	if err != nil {
		return nil, err
	}

	return newRepo(pu)
}

func (p *Provider) Parse(uri string) (provider.Repo, error) {
	pu, err := pkguri.NewFromURI(uri)
	if err != nil {
		return nil, err
	}

	return newRepo(pu)
}

func (p *Provider) Ping(rawurl string) error {
	_, err := pkguri.ParseNpm(rawurl)
	if err != nil {
		return err
	}

	resp, err := http.Get(rawurl)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return provider.ErrNotFound
	}

	return nil
}

type Repo struct {
	pu                *pkguri.PkgURI
	lastPackage       *PackageResponse
	lastDownloadCount *DownloadCountResonse
}

func newRepo(pu *pkguri.PkgURI) (*Repo, error) {
	return &Repo{
		pu: pu,
	}, nil
}

func (r *Repo) Analyze() (provider.Report, error) {
	return r.analyze()
}

func (r *Repo) Derived() ([]string, error) {
	if r.lastPackage == nil {
		if _, err := r.analyze(); err != nil {
			return nil, err
		}
	}

	derived := make([]string, 0)

	if r.lastPackage.Homepage != "" {
		if u, err := url.Parse(r.lastPackage.Homepage); err == nil {
			u.Fragment = ""
			derived = append(derived, u.String())
		}
	}

	if r.lastPackage.Repository.URL != "" {
		if u, err := url.Parse(r.lastPackage.Repository.URL); err == nil {
			// Splits "git+https" scheme
			schemes := strings.Split(u.Scheme, "+")
			for _, scheme := range schemes {
				u2 := u
				u2.Scheme = scheme
				derived = append(derived, u2.String())
			}
		}
	}

	return derived, nil
}

func (r *Repo) analyze() (*Report, error) {
	client := &Client{client: &http.Client{}}

	pkg, err := client.Package(r.pu.URI)
	if err != nil {
		logrus.Debugln("npm:", err)
		return nil, err
	}

	dlcount, err := client.DownloadCount(r.pu.URI, &DownloadCountOpts{LastMonth: true})
	if err != nil {
		return nil, err
	}

	now := field.Now()
	latest := pkg.Versions[pkg.DistTags["latest"]]

	stats := []*model.Stat{
		&model.Stat{
			Kind:       model.DownloadCountStat,
			Value:      strconv.FormatUint(dlcount.Downloads, 10),
			IsLatest:   true,
			RecordedAt: now,
		},
		&model.Stat{
			Kind:       model.FileCountStat,
			Value:      strconv.FormatUint(latest.Dist.FileCount, 10),
			IsLatest:   true,
			RecordedAt: now,
		},
		&model.Stat{
			Kind:       model.SizeStat,
			Value:      strconv.FormatUint(latest.Dist.UnpackedSize, 10),
			IsLatest:   true,
			RecordedAt: now,
		},
	}

	tags := model.StatList{}
	for tag, pkgVer := range pkg.Versions {
		releasedAt, err := time.Parse(time.RFC3339, pkg.Time[tag])
		if err != nil {
			continue
		}

		tags = append(tags, &model.Stat{
			Kind:       model.TagStat,
			Value:      tag,
			Checksum:   pkgVer.Dist.Shasum,
			RecordedAt: field.Time(releasedAt),
			IsLatest:   (tag == pkg.DistTags["latest"]),
		})
	}
	for channel, tag := range pkg.DistTags {
		substats := model.StatList{
			&model.Stat{
				Kind:  model.AliasStat,
				Name:  channel,
				Value: tag,
			},
		}

		tags = append(tags, &model.Stat{
			Kind:     model.ChannelStat,
			Value:    channel,
			IsLatest: channel == "latest",
			Substats: &substats,
		})
	}

	sort.Sort(model.StatBySemver(tags))

	r.lastPackage = pkg
	r.lastDownloadCount = dlcount
	return newReport(r.pu, stats, tags)
}

func (r *Repo) Close() error {
	return nil
}

type Report struct {
	*prvdutils.StaticReport
}

func newReport(pu *pkguri.PkgURI, stats, tags []*model.Stat) (*Report, error) {
	report := prvdutils.NewStaticReport(pu, stats, tags)
	return &Report{report}, nil
}

var _ provider.Provider = &Provider{}
var _ provider.Repo = &Repo{}
var _ provider.Report = &Report{}
