package npm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	ru := "https://www.npmjs.org/package/vue"
	u, err := ParseURL(ru)
	assert.Nil(t, err)
	assert.Equal(t, "vue", u.URI)
	assert.Equal(t, ru, u.String())
}
