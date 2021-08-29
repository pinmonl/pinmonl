package monlurl

import (
	"net/url"
	"testing"

	"github.com/pinmonl/pinmonl/monlvars"
	"github.com/stretchr/testify/assert"
)

func TestSplitExplicit(t *testing.T) {
	tests := []struct {
		rawurl     string
		wantMonler string
		wantUrl    string
	}{
		{
			rawurl:     "test,https://example.com",
			wantMonler: "test",
			wantUrl:    "https://example.com",
		},
		{
			rawurl:     "https://example.com",
			wantMonler: "",
			wantUrl:    "https://example.com",
		},
	}

	for _, test := range tests {
		gotMonler, gotUrl := SplitExplicit(test.rawurl)
		wantMonler, wantUrl := test.wantMonler, test.wantUrl
		assert.Equal(t, wantMonler, gotMonler)
		assert.Equal(t, wantUrl, gotUrl)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		rawurl string
		want   *Url
		err    error
	}{
		{
			rawurl: "https://example.com",
			want:   nil,
			err:    ErrNoMatch,
		},
		{
			rawurl: "git,https://git.example.com",
			want: &Url{
				MonlerName: monlvars.Git,
				ObjectName: monlvars.ObjectGitRepo,
				Name:       "/",
				Scheme:     "https",
				Host:       "git.example.com",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		got, err := Parse(test.rawurl)
		if assert.Equal(t, test.err, err) && err == nil {
			want := test.want
			assert.Equal(t, want, got)
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		rawurl string
		want   *url.URL
	}{
		{
			rawurl: "https://example.com/foo/bar/",
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/foo/bar",
			},
		},
		{
			rawurl: "https://some.example.com",
			want: &url.URL{
				Scheme: "https",
				Host:   "some.example.com",
				Path:   "/",
			},
		},
	}

	for _, test := range tests {
		u, _ := url.Parse(test.rawurl)
		got := Normalize(u)
		want := test.want
		assert.Equal(t, want, got)
	}
}
