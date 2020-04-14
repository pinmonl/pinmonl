package github

import (
	"net/http"
	"time"

	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pkg/payload"
)

type apiClient struct {
	client *http.Client
}

func (a *apiClient) getRepo(owner, name string) (*packageResponse, error) {
	dest := DefaultAPIEndpoint + "/repos/" + owner + "/" + name
	req, err := http.NewRequest("GET", dest, nil)
	if err != nil {
		return nil, err
	}
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, monler.ErrNotExist
	}
	defer res.Body.Close()
	var out packageResponse
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type packageResponse struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login string `json:"login"`
		Type  string `json:"type"`
		URL   string `json:"url"`
	} `json:"owner"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PushedAt         time.Time `json:"pushed_at"`
	GitURL           string    `json:"git_url"`
	SSHURL           string    `json:"ssh_url"`
	CloneURL         string    `json:"clone_url"`
	Size             int       `json:"size"`
	StargazersCount  int       `json:"stargazers_count"`
	WatchersCount    int       `json:"watchers_count"`
	ForksCount       int       `json:"forks_count"`
	HasIssues        bool      `json:"has_issues"`
	OpenIssuesCount  int       `json:"open_issues_count"`
	SubscribersCount int       `json:"subscribers_count"`
}
