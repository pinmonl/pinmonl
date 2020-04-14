package docker

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pkg/payload"
)

type apiClient struct {
	client *http.Client
}

func (a *apiClient) getRepository(uri string) (*repositoryResponse, error) {
	req, err := http.NewRequest("GET", DefaultAPIEndpoint+"/v2/repositories/"+uri, nil)
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
	var out repositoryResponse
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *apiClient) listAllTags(uri string) ([]tagResponse, error) {
	size := 25
	count := 0
	var out []tagResponse
	for page := 1; (page-1)*size <= count; page++ {
		res, err := a.listTags(uri, page, size)
		if err != nil {
			return nil, err
		}
		out = append(out, res.Results...)
		count = res.Count
	}
	return out, nil
}

func (a *apiClient) listTags(uri string, page, pageSize int) (*tagsResponse, error) {
	req, err := http.NewRequest("GET", DefaultAPIEndpoint+"/v2/repositories/"+uri+"/tags", nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = fmt.Sprintf("page=%d&page_size=%d", page, pageSize)
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, monler.ErrNotExist
	}
	defer res.Body.Close()
	var out tagsResponse
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type repositoryResponse struct {
	User            string    `json:"user"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	RepositoryType  string    `json:"repository_type"`
	Status          int       `json:"status"`
	Description     string    `json:"description"`
	IsPrivate       bool      `json:"is_private"`
	IsAutomated     bool      `json:"is_automated"`
	CanEdit         bool      `json:"can_edit"`
	StarCount       int       `json:"star_count"`
	PullCount       int       `json:"pull_count"`
	LastUpdated     time.Time `json:"last_updated"`
	IsMigrated      bool      `json:"is_migrated"`
	HasStarred      bool      `json:"has_starred"`
	FullDescription string    `json:"full_description"`
	// "affiliation"
	Permissions struct {
		Read  bool `json:"read"`
		Write bool `json:"write"`
		Admin bool `json:"admin"`
	} `json:"permissions"`
}

type listResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type tagsResponse struct {
	listResponse
	Results []tagResponse `json:"results"`
}

type tagResponse struct {
	Name        string             `json:"name"`
	FullSize    int                `json:"full_size"`
	LastUpdated time.Time          `json:"last_updated"`
	Images      []tagImageResponse `json:"images"`
}

type tagImageResponse struct {
	Architecture string `json:"architecture"`
	Features     string `json:"features"`
	Variant      string `json:"variant"`
	Digest       string `json:"digest"`
	Os           string `json:"os"`
	OsFeatures   string `json:"os_features"`
	OsVersion    string `json:"os_version"`
	Size         int    `json:"size"`
}
