package monl

// Report handles the vendor API response and transform to repository information
// releases should be in descending order by date
type Report interface {
	RawURL() string
	URI() string
	Vendor() string
	Popularity() StatCollection
	Latest() Stat
	Next() Stat
	Previous() Stat
	Length() int
}
