package git

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractURLs(t *testing.T) {
	getReadme := func(url string) string {
		res, _ := http.Get(url)
		b, _ := ioutil.ReadAll(res.Body)
		return string(b)
	}

	urls := extractURLs(getReadme("https://raw.githubusercontent.com/vuejs/vue/dev/README.md"))
	assert.NotNil(t, urls)
	assert.Less(t, 1, len(urls))
}

func TestReport(t *testing.T) {
	r, err := NewReport(&ReportOpts{URL: "https://github.com/ahshum/release"})
	assert.Nil(t, err)
	assert.Nil(t, r.Download())
	if r != nil {
		assert.Equal(t, 2, r.Len())
	}
}
