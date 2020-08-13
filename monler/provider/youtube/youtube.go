package youtube

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/monler/prvdutils"
	"github.com/pinmonl/pinmonl/pkgs/pkgdata"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"google.golang.org/api/youtube/v3"
)

type Provider struct {
	client *Client
}

func NewProvider(apiKeys []string) (*Provider, error) {
	client, err := NewClient(apiKeys)
	if err != nil {
		return nil, err
	}

	return &Provider{
		client: client,
	}, nil
}

func (*Provider) ProviderName() string {
	return pkgdata.YoutubeProvider
}

func (p *Provider) Open(rawurl string) (provider.Repo, error) {
	pu, err := pkguri.ParseYoutube(rawurl)
	if err != nil {
		return nil, err
	}

	return newRepo(p.client, pu.URI)
}

func (p *Provider) Parse(uri string) (provider.Repo, error) {
	pu, err := pkguri.NewFromURI(uri)
	if err != nil {
		return nil, err
	}

	return newRepo(p.client, pu.URI)
}

func (*Provider) Ping(rawurl string) error {
	if _, err := pkguri.ParseYoutube(rawurl); err != nil {
		return err
	}

	res, err := http.Get(rawurl)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return provider.ErrNotSupport
	}
	return nil
}

type Repo struct {
	client     *Client
	channelId  string
	lastReport *Report
}

func newRepo(client *Client, channelId string) (*Repo, error) {
	return &Repo{
		client:    client,
		channelId: channelId,
	}, nil
}

func (r *Repo) Analyze() (provider.Report, error) {
	return r.analyze()
}

func (r *Repo) analyze() (*Report, error) {
	report, err := newReport(r.client, r.channelId)
	r.lastReport = report
	return report, err
}

func (r *Repo) Derived() ([]string, error) {
	if r.lastReport == nil {
		if _, err := r.analyze(); err != nil {
			return nil, err
		}
	}

	derived := make([]string, 0)

	if !pkguri.IsYoutubeValidChannelId(r.channelId) {
		if pu, err := r.lastReport.URI(); err == nil {
			derived = append(derived, pkguri.ToURL(pu))
		}
	}

	return derived, nil
}

func (r *Repo) Close() error {
	return nil
}

type Report struct {
	*prvdutils.PagesReport
}

func newReport(client *Client, channelId string) (*Report, error) {
	channelResponse, err := client.ChannelsList(channelId, []string{"contentDetails", "statistics"}).Do()
	if err != nil {
		return nil, err
	}
	if len(channelResponse.Items) == 0 {
		return nil, provider.ErrNotFound
	}

	channel := channelResponse.Items[0]
	uploadId := channel.ContentDetails.RelatedPlaylists.Uploads

	pu := &pkguri.PkgURI{
		Provider: pkgdata.YoutubeProvider,
		URI:      channel.Id,
		Proto:    pkguri.DefaultProto,
	}
	statFn := reportStatFn(channel)
	videoFn := reportVideoFn(client, uploadId)

	report := prvdutils.NewPagesReport(pu, statFn, videoFn)
	return &Report{PagesReport: report}, nil
}

func reportStatFn(channel *youtube.Channel) prvdutils.PageFunc {
	return func(_ int64) ([]*model.Stat, int64, bool, error) {
		now := field.Now()
		stats := []*model.Stat{
			&model.Stat{
				Kind:       model.SubscriberCountStat,
				Value:      strconv.FormatUint(channel.Statistics.SubscriberCount, 10),
				RecordedAt: now,
				IsLatest:   true,
			},
			&model.Stat{
				Kind:       model.VideoCountStat,
				Value:      strconv.FormatUint(channel.Statistics.VideoCount, 10),
				RecordedAt: now,
				IsLatest:   true,
			},
			&model.Stat{
				Kind:       model.ViewCountStat,
				Value:      strconv.FormatUint(channel.Statistics.ViewCount, 10),
				RecordedAt: now,
				IsLatest:   true,
			},
		}

		slen := int64(len(stats))
		return stats, slen, false, nil
	}
}

func reportVideoFn(client *Client, playlistId string) prvdutils.PageFunc {
	var (
		total         int64
		perPage       int64
		nextPageToken string
	)

	return func(_ int64) ([]*model.Stat, int64, bool, error) {
		call := client.PlaylistItemsList(playlistId, []string{"snippet"})

		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		itemResponse, err := call.Do()
		if err != nil {
			return nil, 0, false, err
		}

		total = itemResponse.PageInfo.TotalResults
		perPage = itemResponse.PageInfo.ResultsPerPage
		nextPageToken = itemResponse.NextPageToken

		videos := model.StatList{}
		for _, item := range itemResponse.Items {
			v, err := parsePlaylistItem(item)
			if err != nil {
				return nil, 0, false, err
			}
			videos = append(videos, v)
		}

		// TODO: enable next page flag.
		return videos, total, false, nil
	}
}

func parsePlaylistItem(item *youtube.PlaylistItem) (*model.Stat, error) {
	recordedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
	if err != nil {
		return nil, err
	}

	return &model.Stat{
		Kind:       model.VideoStat,
		RecordedAt: field.Time(recordedAt),
		Value:      item.Snippet.ResourceId.VideoId,
		IsLatest:   item.Snippet.Position == 0,
	}, nil
}

var _ provider.Provider = &Provider{}
var _ provider.Repo = &Repo{}
var _ provider.Report = &Report{}
