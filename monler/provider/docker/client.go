package docker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	client *http.Client
}

func (c *Client) get(dest string, out interface{}) error {
	dest = "https://hub.docker.com/" + dest

	resp, err := c.client.Get(dest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) Repository(repoName string) (*RepositoryResponse, error) {
	dest := "v2/repositories/" + repoName

	var out RepositoryResponse
	if err := c.get(dest, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Tags(repoName string, page int64) (*TagListResponse, error) {
	qs := url.Values{}
	qs.Set("page_size", "100")
	qs.Set("page", strconv.FormatInt(page, 10))
	dest := fmt.Sprintf("v2/repositories/%s/tags?%s", repoName, qs.Encode())

	var out TagListResponse
	if err := c.get(dest, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Builds(repoName string) (*BuildListResponse, error) {
	qs := url.Values{}
	qs.Set("image", repoName)
	dest := fmt.Sprintf("api/build/v1/source?%s", qs.Encode())

	var out BuildListResponse
	if err := c.get(dest, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) User(username string) (*UserResponse, error) {
	dest := fmt.Sprintf("v2/users/%s", username)

	var out UserResponse
	if err := c.get(dest, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type RepositoryResponse struct {
	User            string          `json:"user"`
	Name            string          `json:"name"`
	Namespace       string          `json:"namespace"`
	RepositoryType  string          `json:"repository_type"`
	Status          int64           `json:"status"`
	Description     string          `json:"description"`
	IsPrivate       bool            `json:"is_private"`
	IsAutomated     bool            `json:"is_automated"`
	CanEdit         bool            `json:"can_edit"`
	StarCount       int64           `json:"star_count"`
	PullCount       int64           `json:"pull_count"`
	LastUpdated     string          `json:"last_updated"`
	IsMigrated      bool            `json:"is_migrated"`
	HasStarred      bool            `json:"has_starred"`
	FullDescription string          `json:"full_description"`
	Affiliation     string          `json:"affiliation"`
	Permissions     map[string]bool `json:"permissions"`
}

type TagListResponse struct {
	Count    int64            `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []*RepositoryTag `json:"results"`
}

type RepositoryTag struct {
	Creator             int64                 `json:"creator"`
	Id                  int64                 `json:"id"`
	ImageId             int64                 `json:"image_id"`
	Images              []*RepositoryTagImage `json:"images"`
	LastUpdated         string                `json:"last_updated"`
	LastUpdater         int64                 `json:"last_updater"`
	LastUpdaterUsername string                `json:"last_updater_username"`
	Name                string                `json:"name"`
	Repository          int64                 `json:"repository"`
	FullSize            int64                 `json:"full_size"`
	V2                  bool                  `json:"v2"`
}

type RepositoryTagImage struct {
	Architecture string `json:"architecture"`
	Features     string `json:"features"`
	Variant      string `json:"variant"`
	Digest       string `json:"digest"`
	Os           string `json:"os"`
	OsFeatures   string `json:"os_features"`
	OsVersion    string `json:"os_version"`
	Size         int64  `json:"size"`
}

type BuildListResponse struct {
	Meta    *BuildListMeta `json:"meta"`
	Objects []*BuildObject `json:"objects"`
}

type BuildObject struct {
	Autotests     string   `json:"autotests"`
	BuildInFarm   bool     `json:"build_in_farm"`
	BuildSettings []string `json:"build_settings"`
	Channel       string   `json:"channel"`
	Image         string   `json:"image"`
	Owner         string   `json:"owner"`
	Provider      string   `json:"provider"`
	RepoLinks     bool     `json:"repo_links"`
	Repository    string   `json:"repository"`
	ResourceUri   string   `json:"resource_uri"`
	State         string   `json:"state"`
	Uuid          string   `json:"uuid"`
}

type BuildListMeta struct {
	Limit      int64  `json:"limit"`
	Next       string `json:"next"`
	Offset     int64  `json:"offset"`
	Previous   string `json:"previous"`
	TotalCount int64  `json:"total_count"`
}

type UserResponse struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Orgname       string `json:"orgname"`
	FullName      string `json:"full_name"`
	Location      string `json:"location"`
	Company       string `json:"company"`
	ProfileUrl    string `json:"profile_url"`
	DateJoined    string `json:"date_joined"`
	GravatarUrl   string `json:"gravatar_url"`
	GravatarEmail string `json:"gravatar_email"`
	Type          string `json:"type"`
}
