package model

import "github.com/pinmonl/pinmonl/model/field"

// Pkg defines the package vendor under monl.
type Pkg struct {
	ID           string       `json:"id"           db:"pkg_id"`
	URL          string       `json:"url"          db:"pkg_url"`
	Provider     string       `json:"provider"     db:"pkg_provider"`
	ProviderHost string       `json:"providerHost" db:"pkg_provider_host"`
	ProviderURI  string       `json:"providerUri"  db:"pkg_provider_uri"`
	Title        string       `json:"title"        db:"pkg_title"`
	Description  string       `json:"description"  db:"pkg_description"`
	Readme       string       `json:"readme"       db:"pkg_readme"`
	ImageID      string       `json:"imageId"      db:"pkg_image_id"`
	Labels       field.Labels `json:"labels"       db:"pkg_labels"`
	CreatedAt    field.Time   `json:"createdAt"    db:"pkg_created_at"`
	UpdatedAt    field.Time   `json:"updatedAt"    db:"pkg_updated_at"`
}

// PkgList is slice of Pkg.
type PkgList []Pkg

// Keys gets ID slice from Pkgs.
func (pl PkgList) Keys() []string {
	out := make([]string, len(pl))
	for i, p := range pl {
		out[i] = p.ID
	}
	return out
}
