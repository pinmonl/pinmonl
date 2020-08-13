package youtube

import (
	"context"
	"errors"

	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	*youtube.Service
	ctx context.Context
}

func NewClient(apiKeys []string) (*Client, error) {
	if len(apiKeys) == 0 {
		return nil, errors.New("please provide at least one api key")
	}

	ctx := context.TODO()
	svc, err := youtube.NewService(ctx, option.WithAPIKey(apiKeys[0]))
	if err != nil {
		return nil, err
	}

	return &Client{
		Service: svc,
		ctx:     ctx,
	}, nil
}

func (c *Client) ChannelsList(id string, part []string) *youtube.ChannelsListCall {
	call := c.Service.Channels.List(part)

	if pkguri.IsYoutubeValidChannelId(id) {
		call = call.Id(id)
	} else {
		call = call.ForUsername(id)
	}

	return call
}

func (c *Client) PlaylistItemsList(playlistId string, part []string) *youtube.PlaylistItemsListCall {
	call := c.Service.PlaylistItems.List(part).
		PlaylistId(playlistId).
		MaxResults(50)

	return call
}
