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
				Proto:    DefaultProto,
			},
		},
		{
			rawurl: "gitlab://gitlab.com/group/subgroup/repo",
			expect: &PkgURI{
				Provider: "gitlab",
				Host:     "gitlab.com",
				URI:      "group/subgroup/repo",
				Proto:    DefaultProto,
			},
		},
		{
			rawurl: "pvd://provider.com/owner",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner",
				Proto:    DefaultProto,
			},
		},
		{
			rawurl: "pvd:///owner",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "",
				URI:      "owner",
				Proto:    DefaultProto,
			},
		},
		{
			rawurl: "pvd://provider.com/owner/repo?proto=git",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner/repo",
				Proto:    "git",
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
		assert.Equal(t, want.Proto, got.Proto)
		assert.Equal(t, test.rawurl, pu.String())
	}

	_, err := Parse("pvd://provider.com")
	assert.Equal(t, ErrNoURI, err)
}
