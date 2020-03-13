package api

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

func TestGetRepo(t *testing.T) {
	testCases := []struct {
		name      string
		repoOwner string
		repoName  string
		passFn    func(*Repo, error) bool
	}{
		{
			name:      "ahshum/empty",
			repoOwner: "ahshum",
			repoName:  "empty",
			passFn: func(repo *Repo, err error) bool {
				return err == nil &&
					repo != nil &&
					len(repo.Releases.Nodes) == 0
			},
		},
	}

	c := NewClient(createClient())
	for _, tc := range testCases {
		r, err := c.GetRepo(context.Background(), tc.repoOwner, tc.repoName)
		if !tc.passFn(r, err) {
			t.Errorf("case %q fails", tc.name)
		}
	}
}

func TestListRepoReleases(t *testing.T) {
	//
}

func TestListRepoTags(t *testing.T) {
	//
}
