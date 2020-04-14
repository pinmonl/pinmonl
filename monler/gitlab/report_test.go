package gitlab

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	c := &http.Client{}
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
		assert.Nil(t, rep.LatestTag())
	}

	rep, err = NewReport(&ReportOpts{
		URI:    "ahshum/release",
		Client: c,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		assert.Equal(t, true, rep.Next())
		assert.Equal(t, "v2.0.0", rep.Tag().Value)
		assert.Equal(t, true, rep.Next())
		assert.Equal(t, "v1.0.0", rep.Tag().Value)
	}
}
