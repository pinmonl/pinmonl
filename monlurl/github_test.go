package monlurl

import (
	"net/url"
	"testing"

	"github.com/pinmonl/pinmonl/monlvars"
	"github.com/stretchr/testify/assert"
)

func TestParseGithubUrl(t *testing.T) {
	tests := []struct {
		rawurl string
		want   *Url
	}{
		{
			rawurl: "https://github.com/owner/repo",
			want: &Url{
				MonlerName: monlvars.Github,
				ObjectName: monlvars.ObjectGithubRepo,
				Host:       monlvars.GithubHost,
				Prefix:     "owner",
				Name:       "repo",
			},
		},
		{
			rawurl: "https://github.com/owner",
			want: &Url{
				MonlerName: monlvars.Github,
				ObjectName: monlvars.ObjectGithubUser,
				Host:       monlvars.GithubHost,
				Name:       "owner",
			},
		},
	}

	for _, test := range tests {
		u, _ := url.Parse(test.rawurl)
		got, err := ParseGithubUrl(u)
		want := test.want

		if assert.Nil(t, err) {
			assert.Equal(t, want, got)
		}
	}
}

func TestBuildGithubUrl(t *testing.T) {
	tests := []struct {
		url  *Url
		want string
	}{
		{
			url: &Url{
				ObjectName: monlvars.ObjectGithubRepo,
				Prefix:     "owner",
				Name:       "repo",
			},
			want: "https://github.com/owner/repo",
		},
		{
			url: &Url{
				ObjectName: monlvars.ObjectGithubUser,
				Name:       "user",
			},
			want: "https://github.com/user",
		},
	}

	for _, test := range tests {
		got := BuildGithubUrl(test.url)
		want := test.want
		assert.Equal(t, want, got)
	}
}
