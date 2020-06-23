package pkguri

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		rawurl string
		expect *PkgURI
		err    error
	}{
		{
			rawurl: "pvd://provider.com/owner/repo",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner/repo",
				Proto:    DefaultProto,
			},
			err: nil,
		},
		{
			rawurl: "gitlab://gitlab.com/group/subgroup/repo",
			expect: &PkgURI{
				Provider: "gitlab",
				Host:     "gitlab.com",
				URI:      "group/subgroup/repo",
				Proto:    DefaultProto,
			},
			err: nil,
		},
		{
			rawurl: "pvd://provider.com/owner",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner",
				Proto:    DefaultProto,
			},
			err: nil,
		},
		{
			rawurl: "pvd:///owner",
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "",
				URI:      "owner",
				Proto:    DefaultProto,
			},
			err: nil,
		},
		{
			rawurl: "pvd://provider.com/owner/repo?proto=" + url.QueryEscape("git+ssh"),
			expect: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner/repo",
				Proto:    "git+ssh",
			},
			err: nil,
		},
		{
			rawurl: "pvd://provider.com",
			expect: nil,
			err:    ErrNoURI,
		},
	}

	for _, test := range tests {
		pu, err := unmarshal(test.rawurl)
		got, want := pu, test.expect
		if assert.Equal(t, test.err, err) && err == nil {
			assert.Equal(t, want.Provider, got.Provider)
			assert.Equal(t, want.Host, got.Host)
			assert.Equal(t, want.URI, got.URI)
			assert.Equal(t, want.Proto, got.Proto)
			assert.Equal(t, test.rawurl, pu.String())
		}
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		pu     *PkgURI
		expect string
		err    error
	}{
		{
			pu: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner/repo",
				Proto:    "https",
			},
			expect: "pvd://provider.com/owner/repo",
			err:    nil,
		},
		{
			pu: &PkgURI{
				Provider: "pvd",
				Host:     "provider.com",
				URI:      "owner/repo",
				Proto:    "git+ssh",
			},
			expect: "pvd://provider.com/owner/repo?proto=" + url.QueryEscape("git+ssh"),
			err:    nil,
		},
		{
			pu: &PkgURI{
				Provider: "github",
				Host:     "",
				URI:      "owner/repo",
			},
			expect: "github:///owner/repo",
			err:    nil,
		},
	}

	for _, test := range tests {
		pkguri, err := marshal(test.pu)
		if assert.Equal(t, test.err, err) && err == nil {
			assert.Equal(t, test.expect, pkguri)
		}
	}
}
