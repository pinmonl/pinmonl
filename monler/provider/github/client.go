package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pinmonl/pinmonl/pkgs/pkgrepo"
)

type Client struct {
	client *http.Client
	tokens *TokenStore
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	client := c.client
	if client == nil {
		ti, err := c.tokens.Get()
		if err != nil {
			return nil, err
		}

		ti.Lock()
		defer ti.Unlock()

		client = &http.Client{
			Transport: &Transport{tokenInfo: ti},
		}
	}
	return client.Do(req)
}

func (c *Client) GetRepository(owner, repo string) (*ApiRepoResponse, error) {
	query := `{
  repository(owner: "` + owner + `", name: "` + repo + `") {
    stargazers {
      totalCount
    }
    updatedAt
    watchers {
      totalCount
    }
    homepageUrl
    issues(filterBy: {states: OPEN}) {
      totalCount
    }
    pullRequests(states: OPEN) {
      totalCount
    }
    isArchived
    isDisabled
    forkCount
    isMirror
    primaryLanguage {
      name
      color
    }
    licenseInfo {
      name
      key
    }
  }
}`

	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(struct {
		Query string `json:"query"`
	}{query})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", body)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github: api response got %d", resp.StatusCode)
	}

	var info struct {
		Data struct {
			Repo ApiRepoResponse `json:"repository"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info.Data.Repo, nil
}

type Transport struct {
	base      http.RoundTripper
	tokenInfo *TokenInfo
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	baset := t.base
	if baset == nil {
		baset = http.DefaultTransport
	}
	if t.tokenInfo != nil && t.tokenInfo.token != "" {
		req.Header.Add("Authorization", "Bearer "+t.tokenInfo.token)
	}

	resp, err := baset.RoundTrip(req)
	if resp != nil {
		t.tokenInfo.UpdateFromHeader(resp.Header)
	}
	return resp, err
}

type ApiRepoResponse struct {
	ForkCount       int                `json:"forkCount"`
	HomepageUrl     string             `json:"homepageUrl"`
	IsArchived      bool               `json:"isArchived"`
	IsDisabled      bool               `json:"isDisabled"`
	IsMirror        bool               `json:"isMirror"`
	UpdatedAt       time.Time          `json:"updatedAt"`
	Stargazers      ApiCountResponse   `json:"stargazers"`
	Watchers        ApiCountResponse   `json:"watchers"`
	Issues          ApiCountResponse   `json:"issues"`
	PullRequests    ApiCountResponse   `json:"pullRequests"`
	PrimaryLanguage ApiLangResponse    `json:"primaryLanguage"`
	LicenseInfo     ApiLicenseResponse `json:"licenseInfo"`
}

type ApiCountResponse struct {
	TotalCount int `json:"totalCount"`
}

type ApiLangResponse struct {
	Name  RepoLanguage `json:"name"`
	Color string       `json:"color"`
}

type ApiLicenseResponse struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type RepoLanguage string

func (rl RepoLanguage) String() string {
	return strings.ToLower(string(rl))
}

func (rl RepoLanguage) Language() pkgrepo.Language {
	switch rl.String() {
	case "javascript":
		return pkgrepo.Javascript
	default:
		return pkgrepo.Undefined
	}
}
