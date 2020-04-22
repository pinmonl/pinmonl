package npm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	ru := "https://www.npmjs.com/package/vue"
	u, err := ParseURL(ru)
	assert.Nil(t, err)
	if u != nil {
		assert.Equal(t, "vue", u.URI)
		assert.Equal(t, ru, u.String())
	}

	ru = "https://npmjs.com/package/vue"
	u, err = ParseURL(ru)
	assert.Nil(t, err)
	if u != nil {
		assert.Equal(t, "vue", u.URI)
	}
}
