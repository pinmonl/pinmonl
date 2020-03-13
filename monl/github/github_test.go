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
			url:   "http://invalid_domain/ahshum/empty",
			wants: false,
		},
		{
			url:   "http://github.com/ahshum/empty",
			wants: true,
		},
		{
			url:   "https://github.com/ahshum/not-existed",
			wants: false,
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
