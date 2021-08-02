package monl

import (
	"testing"

	"github.com/pinmonl/pinmonl/monlvars"
	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		rawurl string
		want   *URL
	}{
		{
			rawurl: "https://github.com/owner/repo",
			want: &URL{
				MonlerName: monlvars.Github,
				ObjectName: monlvars.GithubRepo,
				Prefix:     "owner",
				Name:       "repo",
			},
		},
		{
			rawurl: "https://github.com/owner",
			want: &URL{
				MonlerName: monlvars.Github,
				ObjectName: monlvars.GithubUser,
				Name:       "owner",
			},
		},
	}

	for _, test := range tests {
		got, err := ParseURL(test.rawurl)
		want := test.want

		if assert.Nil(t, err) {
			assert.Equal(t, want, got)
		}
	}
}
