package helm

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	c := &http.Client{}

	rep, err := NewReport(&ReportOpts{
		URI:    "jetstack/cert-manager",
		Client: c,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		assert.NotNil(t, rep.LatestTag())
		assert.LessOrEqual(t, 1, rep.Len())
	}
}
