package github

import (
	"testing"

	"github.com/pinmonl/pinmonl/monler"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	tests := []struct {
		url  string
		want error
	}{
		{
			url:  "http://invalid_domain/ahshum/empty",
			want: monler.ErrNotSupport,
		},
		{
			url:  "http://github.com/ahshum/empty",
			want: nil,
		},
		{
			url:  "https://github.com/ahshum/not-existed",
			want: monler.ErrNotExist,
		},
	}

	for _, test := range tests {
		got, want := Ping(test.url, nil), test.want
		assert.Equal(t, want, got)
	}
}

func TestExtractURI(t *testing.T) {
	o, r := ExtractURI("owner/repo")
	assert.Equal(t, "owner", o)
	assert.Equal(t, "repo", r)

	o, r = ExtractURI("abc/def")
	assert.Equal(t, "abc", o)
	assert.Equal(t, "def", r)
}
