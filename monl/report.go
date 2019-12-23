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

	// Next returns the next release
	Next() Stat

	// Previous returns the previous release
	Previous() Stat

	// Length returns the total count of release
	Length() int

	// Derived returns key-value pairs with `vendor name` as key and `url` as value
	Derived() map[string]string
}
