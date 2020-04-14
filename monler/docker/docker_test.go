package docker

import (
	"testing"

	"github.com/pinmonl/pinmonl/monler"
	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{
			url:  "https://hub.docker.com/_/golang",
			want: "library/golang",
		},
		{
			url:  "https://hub.docker.com/r/drone/drone",
			want: "drone/drone",
		},
	}

	for _, test := range tests {
		u, err := ParseURL(test.url)
		assert.Nil(t, err)
		got, want := u.URI, test.want
		assert.Equal(t, want, got)
	}
}

func TestParseURLError(t *testing.T) {
	tests := []struct {
		url  string
		want error
	}{
		{
			url:  "https://google.com",
			want: monler.ErrNotSupport,
		},
	}

	for _, test := range tests {
		u, err := ParseURL(test.url)
		assert.Nil(t, u)
		got, want := err, test.want
		assert.Equal(t, want, got)
	}
}
