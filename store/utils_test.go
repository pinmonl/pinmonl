package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUID(t *testing.T) {
	assert.NotEmpty(t, newUID())

	id1 := newUID()
	id2 := newUID()

	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id2, newUID())
}

func TestTimestamp(t *testing.T) {
	assert.False(t, timestamp().Time().IsZero())
}

func TestBindQueryIDs(t *testing.T) {
	prefix := "test"
	ids := []string{"id1", "id2", "id3"}
	ks, args := bindQueryIDs(prefix, ids)
	for i, k := range ks {
		assert.Equal(t, ":", k[:1])
		assert.Equal(t, prefix, k[1:1+len(prefix)])
		got, ok := args[k[1:]]
		assert.True(t, ok)
		want := ids[i]
		assert.Equal(t, want, got)
	}
}
