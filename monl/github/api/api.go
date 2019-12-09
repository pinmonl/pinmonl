package api

import (
	"context"
	"net/http"

	"github.com/shurcooL/githubv4"
)

// Client defines to handle Github v4 API
type Client struct {
	client *githubv4.Client
}

// NewClient creates API client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{client: githubv4.NewClient(httpClient)}
}

// GetRepo returns basic information of the repository
func (c *Client) GetRepo(ctx context.Context, owner, name string) (*Repo, error) {
	var q struct {
		Repo *Repo `graphql:"repository(owner: $owner, name: $name)"`
	}
	qv := c.makeVars(owner, name)
	err := c.client.Query(ctx, &q, qv)
	if err != nil {
		return nil, err
	}
	return q.Repo, nil
}

// ListRepoReleases returns list of releases
func (c *Client) ListRepoReleases(ctx context.Context, owner, name string, relpo *PageOption) (*RepoReleases, error) {
	var q struct {
		Repo *RepoReleases `graphql:"repository(owner: $owner, name: $name)"`
	}
	qv := c.makeVars(owner, name)
	if relpo != nil {
		for k, v := range relpo.Scalar() {
			qv[k] = v
		}
	}
	err := c.client.Query(ctx, &q, qv)
	if err != nil {
		return nil, err
	}
	return q.Repo, nil
}

// ListRepoTags returns list of tags
func (c *Client) ListRepoTags(ctx context.Context, owner, name string, tagpo *PageOption) (*RepoTags, error) {
	var q struct {
		Repo *RepoTags `graphql:"repository(owner: $owner, name: $name)"`
	}
	qv := c.makeVars(owner, name)
	if tagpo != nil {
		for k, v := range tagpo.Scalar() {
			qv[k] = v
		}
	}
	err := c.client.Query(ctx, &q, qv)
	if err != nil {
		return nil, err
	}
	return q.Repo, nil
}

func (c *Client) makeVars(owner, name string) map[string]interface{} {
	return map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}
}
