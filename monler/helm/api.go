package helm

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

func newAPIClient(cred monler.Credential) (*apiClient, error) {
	return &apiClient{client: &http.Client{}}, nil
}

func (a *apiClient) searchChart(repo, query string) ([]*ChartResponse, error) {
	dest := DefaultAPIEndpoint + "/charts"
	if repo != "" {
		dest += "/" + repo
	}
	dest += "/search"
	req, err := http.NewRequest("GET", dest, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = url.Values{"q": []string{query}}.Encode()
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var out struct {
		Data []*ChartResponse
	}
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

func (a *apiClient) getChart(uri string) (*ChartResponse, error) {
	dest := DefaultAPIEndpoint + "/charts/" + uri
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
	var out struct {
		Data *ChartResponse
	}
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

func (a *apiClient) listVersions(uri string) ([]*ChartVersionResponse, error) {
	dest := DefaultAPIEndpoint + "/charts/" + uri + "/versions"
	req, err := http.NewRequest("GET", dest, nil)
	if err != nil {
		return nil, err
	}
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var out struct {
		Data []*ChartVersionResponse
	}
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// ChartResponse defines the structure of chart.
type ChartResponse struct {
	ID         string
	Type       string
	Attributes ChartAttributes
}

// ChartAttributes defines the structure of chart attributes.
type ChartAttributes struct {
	Name string
	Repo struct {
		Name string
		URL  string
	}
	Description string
	Home        string
	Keywords    []string
	Maintainers []ChartMaintainer
	Sources     []string
	Icon        string
}

// ChartMaintainer defines the structure of chart maintainer.
type ChartMaintainer struct {
	Name  string
	Email string
}

// ChartVersionResponse defines the structure of chart version.
type ChartVersionResponse struct {
	ID         string
	Type       string
	Attributes ChartVersionAttributes
}

// ChartVersionAttributes defines the structure of chart version attributes.
type ChartVersionAttributes struct {
	Version    string
	AppVersion string
	Created    time.Time
	Digest     string
	URLs       []string
	Readme     string
	Values     string
}
