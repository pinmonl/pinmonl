package github

import (
	"regexp"
)

// Vendor handles Github url
type Vendor struct {
	token string
}

// NewVendor creates Github vendor
func NewVendor(token string) (*Vendor, error) {
	v := &Vendor{
		token: token,
	}
	return v, nil
}

// Name returns the vendor name
func (v *Vendor) Name() string { return "github" }

// Check passes if the url matches one of the patterns
func (v *Vendor) Check(rawurl string) bool { return v.isValidURL(rawurl) }

// Load returns Github report
func (v *Vendor) Load(rawurl string) (*Report, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: v.token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	r, err := NewReport(v.Name(), rawurl, httpClient)
	if err != nil {
		return nil, err
	}

	r.Download()
	return r, nil
}

func (v *Vendor) isValidURL(rawurl string) bool {
	patterns := []string{
		`^https?://github\.com/([^/]+)/([^/]+)`,
	}
	for _, pattern := range patterns {
		if regexp.MustCompile(pattern).MatchString(rawurl) {
			return true
		}
	}
	return false
}
