package pkguri

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		rawurl string
		expect *PkgURI
	}{
		{
			rawurl: "pvd://provider.com/owner/repo",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner/repo",
			},
		},
		{
			rawurl: "gitlab://gitlab.com/group/subgroup/repo",
			expect: &PkgURI{
				Provider: "gitlab",
				Host:     "gitlab.com",
				URI:      "group/subgroup/repo",
			},
		},
		{
			rawurl: "pvd://provider.com/owner",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner",
			},
		},
		{
			rawurl: "pvd:///owner",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "",
				URI:      "owner",
			},
		},
	}

	for _, test := range tests {
		pu, err := Parse(test.rawurl)
		got, want := pu, test.expect
		assert.Nil(t, err)
		assert.Equal(t, want.Provider, got.Provider)
		assert.Equal(t, want.Host, got.Host)
		assert.Equal(t, want.URI, got.URI)
	}

	_, err := Parse("pvd://provider.com")
	assert.Equal(t, ErrNoURI, err)
}
