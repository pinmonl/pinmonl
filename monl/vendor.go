package monl

import "context"

// Vendor performs low-level url check and creates report.
type Vendor interface {
	// Name of vendor
	Name() string

	// Quick url check for vendor
	Check(url string) bool

	// Load report from url
	Load(ctx context.Context, url string) (Report, error)
}
