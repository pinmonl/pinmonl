package docker

import (
	"fmt"
	"net/http"
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
	return pkgdata.DockerProvider
}

func (p *Provider) Open(rawurl string) (provider.Repo, error) {
	pu, err := pkguri.ParseDocker(rawurl)
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
	_, err := pkguri.ParseDocker(rawurl)
	if err != nil {
		return err
	}

	resp, err := http.Get(rawurl)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return provider.ErrNotSupport
	}
	return nil
}

type Repo struct {
	pu     *pkguri.PkgURI
	client *Client
}

func newRepo(pu *pkguri.PkgURI) (*Repo, error) {
	return &Repo{
		pu:     pu,
		client: &Client{client: &http.Client{}},
	}, nil
}

func (r *Repo) Analyze() (provider.Report, error) {
	return r.analyze()
}

func (r *Repo) Derived() ([]string, error) {
	buildResponse, err := r.client.Builds(r.pu.URI)
	if err != nil {
		return nil, err
	}
	userResponse, err := r.client.User(r.pu.Namespace())
	if err != nil {
		return nil, err
	}

	derived := make([]string, 0)

	if userResponse.ProfileUrl != "" {
		derived = append(derived, userResponse.ProfileUrl)
	}

	if len(buildResponse.Objects) > 0 {
		object := buildResponse.Objects[0]
		if strings.ToLower(object.Provider) == "github" {
			repoUrl := fmt.Sprintf("https://github.com/%s", object.Image)
			derived = append(derived, repoUrl)
		}
	}

	return derived, nil
}

func (r *Repo) analyze() (*Report, error) {
	return newReport(r.client, r.pu)
}

func (r *Repo) Close() error {
	return nil
}

type Report struct {
	*prvdutils.StaticReport
}

func newReport(client *Client, pu *pkguri.PkgURI) (*Report, error) {
	repositoryResponse, err := client.Repository(pu.URI)
	if err != nil {
		return nil, err
	}

	now := field.Now()
	stats := []*model.Stat{
		&model.Stat{
			Kind:       model.PullCountStat,
			Value:      strconv.FormatInt(repositoryResponse.PullCount, 10),
			RecordedAt: now,
			IsLatest:   true,
		},
		&model.Stat{
			Kind:       model.StarCountStat,
			Value:      strconv.FormatInt(repositoryResponse.StarCount, 10),
			RecordedAt: now,
			IsLatest:   true,
		},
	}

	rawtags, err := fetchAllTags(client, pu.URI)
	if err != nil {
		return nil, err
	}
	bucket, err := newTagBucket(rawtags)
	if err != nil {
		return nil, err
	}

	tags := make([]*model.Stat, 0)
	for _, tag := range bucket.semvers {
		if tag == bucket.latest {
			tag.IsLatest = true
		}

		substats := *tag.Substats
		for _, child := range bucket.children[tag] {
			substats = append(substats, &model.Stat{
				Kind:  model.AliasStat,
				Name:  tag.Value,
				Value: child.Value,
			})
		}
		tag.Substats = &substats

		delete(bucket.children, tag)
		tags = append(tags, tag)
	}
	for _, tag := range bucket.channels {
		tag.IsLatest = tag.Value == "latest"
		tag.Kind = model.ChannelStat

		substats := *tag.Substats
		for _, child := range bucket.children[tag] {
			substats = append(substats, &model.Stat{
				Kind:  model.AliasStat,
				Name:  tag.Value,
				Value: child.Value,
			})
		}
		tag.Substats = &substats

		delete(bucket.children, tag)
		stats = append(stats, tag)
	}
	for tag, children := range bucket.children {
		tag.Kind = model.AliasStat
		substats := *tag.Substats
		for _, child := range children {
			substats = append(substats, &model.Stat{
				Kind:  model.AliasStat,
				Name:  tag.Value,
				Value: child.Value,
			})
		}
		tag.Substats = &substats

		tags = append(tags, tag)
	}

	sort.Sort(TagBySemver(tags))

	report := prvdutils.NewStaticReport(pu, nil, stats, tags)
	return &Report{report}, nil
}

func fetchAllTags(client *Client, repoName string) ([]*model.Stat, error) {
	var (
		tags    model.StatList
		page    = int64(1)
		hasNext = true
	)

	for hasNext {
		resp, err := client.Tags(repoName, page)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Results {
			tag, err := parseTag(item)
			if err != nil {
				return nil, err
			}

			tags = append(tags, tag)
		}

		hasNext = resp.Next != ""
		page++
	}

	return tags, nil
}

func parseTag(tag *RepositoryTag) (*model.Stat, error) {
	at, err := time.Parse(time.RFC3339, tag.LastUpdated)
	if err != nil {
		at = time.Time{}
		logrus.Debugf("docker: parse tag last_updated err(%s)", err)
	}

	recordedAt := field.Time(at)
	images := model.StatList{}

	for _, image := range tag.Images {
		substats := model.StatList{
			&model.Stat{
				Name:       "architecture",
				Value:      image.Architecture,
				RecordedAt: recordedAt,
			},
			&model.Stat{
				Name:       "features",
				Value:      image.Features,
				RecordedAt: recordedAt,
			},
			&model.Stat{
				Name:       "variant",
				Value:      image.Variant,
				RecordedAt: recordedAt,
			},
			&model.Stat{
				Name:       "os",
				Value:      image.Os,
				RecordedAt: recordedAt,
			},
			&model.Stat{
				Name:       "os_features",
				Value:      image.OsFeatures,
				RecordedAt: recordedAt,
			},
			&model.Stat{
				Name:       "os_version",
				Value:      image.OsVersion,
				RecordedAt: recordedAt,
			},
			&model.Stat{
				Kind:       model.SizeStat,
				Value:      strconv.FormatInt(image.Size, 10),
				RecordedAt: recordedAt,
			},
		}

		images = append(images, &model.Stat{
			Kind:       model.ManifestStat,
			RecordedAt: recordedAt,
			Value:      image.Os + "/" + image.Architecture,
			Checksum:   image.Digest,
			Substats:   &substats,
		})
	}

	return &model.Stat{
		Kind:       model.TagStat,
		RecordedAt: recordedAt,
		Value:      tag.Name,
		Substats:   &images,
	}, nil
}

var _ provider.Provider = &Provider{}
var _ provider.Repo = &Repo{}
var _ provider.Report = &Report{}
