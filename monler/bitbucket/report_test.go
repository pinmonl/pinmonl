package bitbucket

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	c := &http.Client{}

	rep, err := NewReport(&ReportOpts{
		URI:    "ahshum/release",
		Client: c,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		assert.Equal(t, 2, rep.Len())
	}

	rep, err = NewReport(&ReportOpts{
		URI:    "ahshum/not-exist",
		Client: c,
	})
	assert.Nil(t, err)
	assert.NotNil(t, rep.Download())

	rep, err = NewReport(&ReportOpts{
		URI:    "ahshum/release-lw",
		Client: c,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		assert.Equal(t, 0, rep.Len())
	}
}
