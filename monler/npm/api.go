package npm

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

func (a *apiClient) getPackage(uri string) (*packageResponse, error) {
	req, err := http.NewRequest("GET", DefaultRegistryEndpoint+"/"+url.PathEscape(uri), nil)
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

func (a *apiClient) getDownloadPoint(uri string, start, end time.Time) (*downloadResponse, error) {
	format := "2006-01-02"
	dest := DefaultAPIEndpoint + "/downloads/range/"
	dest += start.Format(format) + ":" + end.Format(format)
	dest += "/" + url.PathEscape(uri)

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
	var out downloadResponse
	if err := payload.UnmarshalJSON(res.Body, &out); err != nil {
		return nil, err
	}
	out.RequestedAt = end
	return &out, nil
}

type packageResponse struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	DistTags    map[string]string         `json:"dist-tags"`
	Versions    map[string]packageVersion `json:"versions"`
	Time        map[string]time.Time      `json:"time"`
	License     string                    `json:"license"`
	Homepage    string                    `json:"homepage"`
}

type packageVersion struct {
	Version string `json:"version"`
	Dist    struct {
		Shasum       string `json:"shasum"`
		Tarball      string `json:"tarball"`
		FileCount    int    `json:"fileCount"`
		UnpackedSize int    `json:"unpackedSize"`
	} `json:"dist"`
}

type downloadResponse struct {
	Start       string    `json:"start"`
	End         string    `json:"end"`
	RequestedAt time.Time `json:"-"`
	Downloads   []struct {
		Downloads int    `json:"downloads"`
		Day       string `json:"day"`
	} `json:"downloads"`
}
