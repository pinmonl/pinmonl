package monl

import "errors"

var (
	// ErrVendorNotFound indicates that the vendor does not found by name or url
	ErrVendorNotFound = errors.New("monl vendor not found")
)

// Opts defines the options of monl initiation.
type Opts struct {
	//
}

// Monl holds list of vendors.
type Monl struct {
	vendors map[string]Vendor
}

// New creates an instance of Monl.
func New(opts Opts) *Monl {
	return &Monl{
		vendors: make(map[string]Vendor),
	}
}

// Register pushes vendor into the list.
func (m *Monl) Register(v Vendor) error {
	m.vendors[v.Name()] = v
	return nil
}

// Get finds vendor by name.
func (m *Monl) Get(name string) (Vendor, error) {
	if v, ok := m.vendors[name]; ok {
		return v, nil
	}
	return nil, ErrVendorNotFound
}

// GuessURL find vendor by calling vendor's Check().
func (m *Monl) GuessURL(url string) []Vendor {
	vs := make([]Vendor, 0)
	for _, v := range m.vendors {
		if v.Check(url) {
			vs = append(vs, v)
		}
	}
	return vs
}
