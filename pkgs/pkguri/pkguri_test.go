package pkguri

import (
	"net/url"
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

func TestParseURL(t *testing.T) {
	tests := []struct {
		rawurl string
		expect *url.URL
		err    error
	}{
		{
			rawurl: "https://github.com/owner/repo",
			expect: &url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "/owner/repo",
			},
			err: nil,
		},
		{
			rawurl: "/owner/repo",
			expect: nil,
			err:    ErrHost,
		},
		{
			rawurl: "github.com/owner/repo",
			expect: &url.URL{
				Scheme: DefaultProto,
				Host:   "github.com",
				Path:   "/owner/repo",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		got, err := ParseURL(test.rawurl)
		assert.Equal(t, test.err, err)
		if test.expect != nil {
			assert.Equal(t, test.expect.Scheme, got.Scheme)
			assert.Equal(t, test.expect.Host, got.Host)
			assert.Equal(t, test.expect.Path, got.Path)
		}
	}
}
