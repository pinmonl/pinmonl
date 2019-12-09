package github

import (
	"testing"
)

func TestVendorCheck(t *testing.T) {
	tests := []struct {
		url   string
		wants bool
	}{
		{
			url:   "http://invalid_domain/ahxshum/empty",
			wants: false,
		},
		{
			url:   "http://github.com/ahxshum/empty",
			wants: true,
		},
		{
			url:   "https://github.com/ahxshum/not-existed",
			wants: true,
		},
	}

	v := &Vendor{}
	for _, test := range tests {
		if got := v.Check(test.url); got != test.wants {
			t.Errorf("%q should be %t", test.url, test.wants)
		}
	}
}

func TestVendorReport(t *testing.T) {
	//
}
