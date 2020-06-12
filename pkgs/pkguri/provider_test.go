package pkguri

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToNpm(t *testing.T) {
	tests := []struct {
		pu     *PkgURI
		expect string
	}{
		{
			pu: &PkgURI{
				URI: "my-pkg",
			},
			expect: "https://www.npmjs.com/package/my-pkg",
		},
		{
			pu: &PkgURI{
				Host: "another.npmjs.com",
				URI:  "my-pkg",
			},
			expect: "https://another.npmjs.com/package/my-pkg",
		},
	}

	for _, test := range tests {
		got := ToNpm(test.pu)
		assert.Equal(t, test.expect, got)
	}
}

func TestParseFromNpm(t *testing.T) {
	tests := []struct {
		rawurl string
		expect *PkgURI
		err    error
	}{
		{
			rawurl: "https://www.npmjs.com/package/my-pkg",
			expect: &PkgURI{
				Provider: NpmProvider,
				Host:     NpmHost,
				URI:      "my-pkg",
			},
			err: nil,
		},
		{
			rawurl: "https://www.npmjs.com/package/@org/my-pkg",
			expect: &PkgURI{
				Provider: NpmProvider,
				Host:     NpmHost,
				URI:      "@org/my-pkg",
			},
			err: nil,
		},
		{
			rawurl: "https://another.npmjs.com/package/my-pkg",
			expect: nil,
			err:    ErrHost,
		},
		{
			rawurl: "https://www.npmjs.com/package-somewhere/my-pkg",
			expect: nil,
			err:    ErrPath,
		},
	}

	for _, test := range tests {
		got, err := ParseFromNpm(test.rawurl)
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expect, got)
	}
}

func TestParseFromNpmRegistry(t *testing.T) {
	tests := []struct {
		rawurl string
		expect *PkgURI
		err    error
	}{
		{
			rawurl: "https://registry.npmjs.org/my-pkg",
			expect: &PkgURI{
				Provider: NpmProvider,
				Host:     NpmHost,
				URI:      "my-pkg",
			},
			err: nil,
		},
		{
			rawurl: "https://registry.npmjs.org/@org/my-pkg",
			expect: &PkgURI{
				Provider: NpmProvider,
				Host:     NpmHost,
				URI:      "@org/my-pkg",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		got, err := ParseFromNpmRegistry(test.rawurl)
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expect, got)
	}
}

func TestToGithub(t *testing.T) {
	tests := []struct {
		pu     *PkgURI
		expect string
	}{
		{
			pu: &PkgURI{
				Provider: GithubProvider,
				Host:     GithubHost,
				URI:      "owner/repo",
			},
			expect: "https://github.com/owner/repo",
		},
		{
			pu: &PkgURI{
				Provider: GithubProvider,
				Host:     "",
				URI:      "owner/repo",
			},
			expect: "https://github.com/owner/repo",
		},
	}

	for _, test := range tests {
		got := ToGithub(test.pu)
		assert.Equal(t, test.expect, got)
	}
}

func TestParseFromGithub(t *testing.T) {
	tests := []struct {
		rawurl string
		expect *PkgURI
		err    error
	}{
		{
			rawurl: "https://github.com/owner/repo",
			expect: &PkgURI{
				Provider: GithubProvider,
				Host:     GithubHost,
				URI:      "owner/repo",
			},
			err: nil,
		},
		{
			rawurl: "https://github.com/owner",
			expect: nil,
			err:    ErrNoURI,
		},
		{
			rawurl: "https://github.com/owner/repo/tree/dev/folder",
			expect: &PkgURI{
				Provider: GithubProvider,
				Host:     GithubHost,
				URI:      "owner/repo",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		got, err := ParseFromGithub(test.rawurl)
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.expect, got)
	}
}
