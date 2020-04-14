package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	assert.Nil(t, Ping("https://github.com/ahshum/empty", nil))
	assert.Nil(t, Ping("https://github.com/ahshum/release", nil))
	assert.NotNil(t, Ping("https://github.com/ahshum/not-existed", nil))
}
