package github

import (
	"context"
	"regexp"

	"github.com/pinmonl/pinmonl/monl"
	"golang.org/x/oauth2"
)

// Opts defines the options of Github vendor initiation.
type Opts struct {
	Token string
}

// Vendor handles Github url.
type Vendor struct {
	token string
}

// NewVendor creates Github vendor.
func NewVendor(opts Opts) monl.Vendor {
	return &Vendor{
		token: opts.Token,
	}
}

// Name returns the vendor name.
func (v *Vendor) Name() string { return "github" }

// Check passes if the url matches one of the patterns.
func (v *Vendor) Check(rawurl string) bool { return v.isValidURL(rawurl) }

// Load returns Github report.
func (v *Vendor) Load(ctx context.Context, rawurl string) (monl.Report, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: v.token},
	)
	httpClient := oauth2.NewClient(ctx, src)

	r, err := NewReport(v.Name(), rawurl, httpClient)
	if err != nil {
		return nil, err
	}

	err = r.Download()
	if err != nil {
		return nil, err
	}
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
