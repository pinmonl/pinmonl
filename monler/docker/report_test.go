package docker

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	c := &http.Client{}

	rep, err := NewReport(&ReportOpts{
		URI:    "drone/drone",
		Client: c,
	})
	assert.Nil(t, err)
	if rep != nil {
		assert.Nil(t, rep.Download())
		assert.NotNil(t, rep.LatestTag())
		assert.LessOrEqual(t, 1, rep.Len())
	}
}

func TestParseTagStats(t *testing.T) {
	res := []tagResponse{
		{
			Name:        "1.0.0",
			FullSize:    100,
			LastUpdated: time.Now(),
			Images: []tagImageResponse{
				{Digest: "digest:1"},
			},
		},
	}

	stats, err := parseTagStats(res)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stats))
	if len(stats) > 0 {
		assert.Equal(t, "1.0.0", stats[0].Value)
	}
}

func TestVersionRegex(t *testing.T) {
	tests := []struct {
		version string
		match   bool
		want    string
	}{
		{
			version: "14.2.3.99202004021627-6962-29e9594ubuntu18.04.1-ls69",
			match:   true,
			want:    "14.2.3",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.match, versionRegex.MatchString(test.version))
		assert.Equal(t, test.want, versionRegex.FindString(test.version))
	}
}
