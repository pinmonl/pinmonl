package monl

// Report handles the vendor API response and transform to repository information
// releases should be in descending order by date.
type Report interface {
	// RawURL returns the original url, should exclude the part of hash and search
	RawURL() string

	// URI returns the unique string which can identify its location
	URI() string

	// Vendor returns the name of vendor
	Vendor() string

	// Popularity returns a list contains the numbers which can represent its reputation
	Popularity() StatCollection

	// Latest returns the latest release always
	Latest() Stat

	// Stat returns the stat pointed by the current cursor.
	Stat() Stat

	// Next moves cursor to next stat.
	Next() bool

	// Previous moves cursor to previous stat.
	Previous() bool

	// Close closes the report.
	Close() error

	// Length returns the total count of release
	Length() int

	// Derived returns key-value pairs with `vendor name` as key and `url` as value
	Derived() map[string]string
}
