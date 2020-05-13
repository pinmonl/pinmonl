package github

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestReport(t *testing.T) {
	c := &http.Client{}
	if t := os.Getenv("GITHUB_TOKEN"); t != "" {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: t},
		)
		c = oauth2.NewClient(context.TODO(), src)
	}

	var (
		rep *Report
		err error
	)
	rep, err = NewReport(&ReportOpts{
		URI:    "ahshum/empty",
		Client: c,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		assert.Equal(t, 0, rep.Len())
	}

	rep, err = NewReport(&ReportOpts{
		URI:    "ahshum/not-exist",
		Client: c,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.NotNil(t, rep.Download())
	}

	rep, err = NewReport(&ReportOpts{
		URI:    "ahshum/release",
		Client: c,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		assert.Equal(t, 2, rep.Len())
		assert.Equal(t, true, rep.Next())
		assert.Equal(t, "v2.0.0", rep.Tag().Value)
		assert.Equal(t, true, rep.Next())
		assert.Equal(t, "v1.0.0", rep.Tag().Value)
		assert.Equal(t, 3, len(rep.Stats()))
	}
}
