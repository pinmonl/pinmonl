package monlurl

import (
	"net/url"
	"testing"

	"github.com/pinmonl/pinmonl/monlvars"
	"github.com/stretchr/testify/assert"
)

func TestParseGitUrl(t *testing.T) {
	tests := []struct {
		rawurl string
		want   *Url
	}{
		{
			rawurl: "https://git.example.com/namespace/repo",
			want: &Url{
				MonlerName: monlvars.Git,
				ObjectName: monlvars.ObjectGitRepo,
				Name:       "/namespace/repo",
				Host:       "git.example.com",
				Scheme:     "https",
			},
		},
	}

	for _, test := range tests {
		u, _ := url.Parse(test.rawurl)
		got, err := ParseGitUrl(u)
		if assert.Nil(t, err) {
			want := test.want
			assert.Equal(t, want, got)
		}
	}
}

func TestBuildGitUrl(t *testing.T) {
	tests := []struct {
		url  *Url
		want string
	}{
		{
			url: &Url{
				Name:   "git.example.com/repo",
				Scheme: "https",
			},
			want: "https://git.example.com/repo",
		},
	}

	for _, test := range tests {
		got := BuildGitUrl(test.url)
		want := test.want
		assert.Equal(t, want, got)
	}
}
