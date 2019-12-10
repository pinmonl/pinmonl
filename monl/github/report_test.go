package github

import (
	"context"
	"net/http"
	"os"
	"testing"

	"golang.org/x/oauth2"
)

func createClient() *http.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return httpClient
}

func TestReportUri(t *testing.T) {
	testCases := []struct {
		name   string
		rawurl string
		wants  string
	}{
		{
			name:   "https",
			rawurl: "https://github.com/owner/name",
			wants:  "owner/name",
		},
		{
			name:   "http",
			rawurl: "http://github.com/owner/name",
			wants:  "owner/name",
		},
	}

	for _, tc := range testCases {
		r := &Report{rawurl: tc.rawurl}
		got := r.URI()
		if got != tc.wants {
			t.Errorf("case %q fails, gots %q", tc.name, got)
		}
	}
}

func TestGithubRepo(t *testing.T) {
	testCases := []struct {
		name   string
		rawurl string
		passFn func(*Report, error) bool
	}{
		{
			name:   "ahxshum/empty",
			rawurl: "http://github.com/ahxshum/empty",
			passFn: func(r *Report, err error) bool {
				return err == nil &&
					r != nil &&
					r.Length() == 0 &&
					r.Latest() == nil
			},
		},
		{
			name:   "ahxshum/not-exist",
			rawurl: "https://github.com/ahxshum/not-exist",
			passFn: func(r *Report, err error) bool {
				return err != nil &&
					r != nil
			},
		},
		{
			name:   "ahxshum/release",
			rawurl: "https://github.com/ahxshum/release",
			passFn: func(r *Report, err error) bool {
				if err != nil || r == nil {
					return false
				}

				prev := r.Previous()
				return r.Length() == 2 &&
					r.Latest() != nil &&
					r.Latest().Value() == "v2.0.0" &&
					prev != nil &&
					prev.Value() == "v1.0.0"
			},
		},
	}

	for _, tc := range testCases {
		r, _ := NewReport("test", tc.rawurl, createClient())
		err := r.Download()
		if !tc.passFn(r, err) {
			t.Errorf("case %q fails, err: %v", tc.name, err)
		}
	}
}
