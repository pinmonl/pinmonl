package model

import "github.com/pinmonl/pinmonl/model/field"

// Pkg defines the package vendor under monl.
type Pkg struct {
	ID          string       `json:"id"          db:"id"`
	URL         string       `json:"url"         db:"url"`
	MonlID      string       `json:"monlId"      db:"monl_id"`
	Vendor      string       `json:"vendor"      db:"vendor"`
	VendorURI   string       `json:"vendorUri"   db:"vendor_uri"`
	Title       string       `json:"title"       db:"title"`
	Description string       `json:"description" db:"description"`
	Readme      string       `json:"readme"      db:"readme"`
	ImageID     string       `json:"imageId"     db:"image_id"`
	Labels      field.Labels `json:"labels"      db:"labels"`
	CreatedAt   field.Time   `json:"createdAt"   db:"created_at"`
	UpdatedAt   field.Time   `json:"updatedAt"   db:"updated_at"`
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
