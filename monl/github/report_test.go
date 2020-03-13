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
			name:   "ahshum/empty",
			rawurl: "http://github.com/ahshum/empty",
			passFn: func(r *Report, err error) bool {
				return err == nil &&
					r != nil &&
					r.Length() == 0 &&
					r.Latest() == nil
			},
		},
		{
			name:   "ahshum/not-exist",
			rawurl: "https://github.com/ahshum/not-exist",
			passFn: func(r *Report, err error) bool {
				return err != nil &&
					r != nil
			},
		},
		{
			name:   "ahshum/release",
			rawurl: "https://github.com/ahshum/release",
			passFn: func(r *Report, err error) bool {
				if err != nil || r == nil {
					return false
				}

				return r.Length() == 2 &&
					r.Next() &&
					r.Stat().Value() == "v2.0.0" &&
					r.Next() &&
					r.Stat().Value() == "v1.0.0"
			},
		},
	}

	ctx := context.TODO()
	for _, tc := range testCases {
		r, _ := NewReport("test", tc.rawurl, createClient())
		err := r.Download(ctx)
		if !tc.passFn(r, err) {
			t.Errorf("case %q fails, err: %v", tc.name, err)
		}
	}
}
