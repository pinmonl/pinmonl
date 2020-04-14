package npm

import (
	"net/http"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	client := &http.Client{}

	rep, err := NewReport(&ReportOpts{
		URI:    "vue",
		Client: client,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		if v, _ := semver.NewVersion(rep.LatestTag().Value); v != nil {
			assert.LessOrEqual(t, 2, int(v.Major()))
		}
	}

	rep, err = NewReport(&ReportOpts{
		URI:    "lodash",
		Client: client,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		if v, _ := semver.NewVersion(rep.LatestTag().Value); v != nil {
			assert.LessOrEqual(t, 4, int(v.Major()))
		}
	}
}
