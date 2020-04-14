package bitbucket

import (
	"net/http"

	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pkg/payload"
)

type apiClient struct {
	client *http.Client
}

func (a *apiClient) getForks(workspace, repo string) (*listResponse, error) {
	dest := DefaultAPIEndpoint + "/2.0/repositories/" + workspace + "/" + repo + "/forks"
	req, err := http.NewRequest("GET", dest+"?pagelen=0", nil)
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
	var out listResponse
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (a *apiClient) getWatchers(workspace, repo string) (*listResponse, error) {
	dest := DefaultAPIEndpoint + "/2.0/repositories/" + workspace + "/" + repo + "/watchers"
	req, err := http.NewRequest("GET", dest+"?pagelen=0", nil)
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
	var out listResponse
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type listResponse struct {
	Size     int    `json:"size"`
	Pagelen  int    `json:"pagelen"`
	Page     int    `json:"page"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}
