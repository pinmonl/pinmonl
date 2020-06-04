package model

import (
	"net/url"
	"strings"

	"github.com/pinmonl/pinmonl/model/field"
)

type Pkg struct {
	ID           string     `json:"id"`
	URL          string     `json:"url"`
	Provider     string     `json:"provider"`
	ProviderHost string     `json:"providerHost"`
	ProviderURI  string     `json:"providerUri"`
	CreatedAt    field.Time `json:"createdAt"`
	UpdatedAt    field.Time `json:"updatedAt"`
}

func (p Pkg) MorphKey() string  { return p.ID }
func (p Pkg) MorphName() string { return "pkg" }

func (p Pkg) MarshalURI() ([]byte, error) {
	uri := &url.URL{
		Scheme: p.Provider,
		Host:   p.ProviderHost,
		Path:   p.ProviderURI,
	}
	return []byte(uri.String()), nil
}

func (p *Pkg) UnmarshalURI(data []byte) error {
	u, err := url.Parse(string(data))
	if err != nil {
		return err
	}
	p.Provider = u.Scheme
	p.ProviderHost = u.Host
	p.ProviderURI = strings.Trim(u.Path, "/")
	return nil
}
