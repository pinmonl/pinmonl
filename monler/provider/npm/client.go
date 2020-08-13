package npm

import (
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func (c *Client) get(dest string, out interface{}) error {
	resp, err := c.client.Get(dest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) Package(packageName string) (*PackageResponse, error) {
	dest := "https://registry.npmjs.org/" + packageName

	var out PackageResponse
	if err := c.get(dest, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DownloadCount(packageName string, opts *DownloadCountOpts) (*DownloadCountResonse, error) {
	dest := "https://api.npmjs.org/downloads/point/" + opts.Period() + "/" + packageName

	var out DownloadCountResonse
	if err := c.get(dest, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type DownloadCountOpts struct {
	LastDay   bool
	LastWeek  bool
	LastMonth bool
	StartDate time.Time
	EndDate   time.Time
}

func (o *DownloadCountOpts) Period() string {
	if o.LastDay {
		return "last-day"
	}
	if o.LastWeek {
		return "last-week"
	}
	if o.LastMonth {
		return "last-month"
	}

	if !o.StartDate.IsZero() {
		layout := "2006-01-02"
		period := o.StartDate.Format(layout)
		if !o.EndDate.IsZero() {
			period += ":" + o.EndDate.Format(layout)
		}
		return period
	}

	return ""
}

type PackageResponse struct {
	Name       string                     `json:"name"`
	DistTags   map[string]string          `json:"dist-tags"`
	Versions   map[string]*PackageVersion `json:"versions"`
	Time       map[string]string          `json:"time"`
	Repository *PackageRepository         `json:"repository"`
	Homepage   string                     `json:"homepage"`
	Bugs       *PackageBug                `json:"bugs"`
	License    string                     `json:"license"`
}

type PackageRepository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type PackageBug struct {
	URL string `json:"url"`
}

type PackageVersion struct {
	Name     string              `json:"name"`
	Version  string              `json:"version"`
	Dist     *PackageVersionDist `json:"dist"`
	GitHead  string              `json:"gitHead"`
	Homepage string              `json:"homepage"`
	License  string              `json:"license"`
}

type PackageVersionDist struct {
	Integrity    string `json:"integrity"`
	Shasum       string `json:"shasum"`
	Tarball      string `json:"tarball"`
	FileCount    uint64 `json:"fileCount"`
	UnpackedSize uint64 `json:"unpackedSize"`
}

type DownloadCountResonse struct {
	Downloads uint64 `json:"downloads"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Package   string `json:"package"`
}
