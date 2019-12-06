package monl

import "errors"

var (
	ErrVendorNotFound = errors.New("monl vendor not found")
)

type Monl struct {
	vendors map[string]Vendor
}

func New() *Monl {
	return &Monl{
		vendors: make(map[string]Vendor),
	}
}

func (m *Monl) Add(v Vendor) error {
	m.vendors[v.Name()] = v
	return nil
}

func (m *Monl) Get(name string) (Vendor, error) {
	if v, ok := m.vendors[name]; ok {
		return v, nil
	}
	return nil, ErrVendorNotFound
}

func (m *Monl) GuessUrl(url string) (Vendor, error) {
	for _, v := range m.vendors {
		if v.Check(url) {
			return v, nil
		}
	}
	return nil, ErrVendorNotFound
}
