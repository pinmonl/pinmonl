package monl

// Vendor performs low-level url check and creates report.
type Vendor interface {
	// Name of vendor
	Name() string

	// Quick url check for vendor
	Check(url string) bool

	// Load report from url
	Load(url string) (Report, error)
}
