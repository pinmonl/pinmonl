package gitlab

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pkg/payload"
)

type apiClient struct {
	client *http.Client
}

func (a *apiClient) getProject(id string) (*projectResponse, error) {
	dest := DefaultAPIEndpoint + "/v4/projects/" + url.PathEscape(id)
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
	var out projectResponse
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type projectResponse struct {
	Name              string    `json:"name"`
	NameWithNamespace string    `json:"name_with_namespace"`
	Path              string    `json:"path"`
	PathWithNamespace string    `json:"path_with_namespace"`
	StarCount         int       `json:"star_count"`
	ForksCount        int       `json:"forks_count"`
	CreatedAt         time.Time `json:"created_at"`
	LastActivityAt    time.Time `json:"last_activity_at"`
	SSHURLToRepo      string    `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string    `json:"http_url_to_repo"`
	Namespace         struct {
		Name     string `json:"name"`
		Path     string `json:"path"`
		Kind     string `json:"kind"`
		FullPath string `json:"full_path"`
	} `json:"namespace"`
}
