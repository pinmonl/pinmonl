package git

import (
	"testing"

	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/stretchr/testify/assert"
)

func TestProviderPing(t *testing.T) {
	p, err := NewProvider()
	assert.Nil(t, err)
	assert.Equal(t, ErrNoPing, p.Ping("anything"))
}

func TestNewReport(t *testing.T) {
	var (
		r   *Report
		err error
	)

	// Empty repo.
	r, err = newReport(&pkguri.PkgURI{
		URI: "https://github.com/ahshum/empty",
	})
	assert.Nil(t, err)
	assert.NotNil(t, r)
	r.Close()

	// Annotated tags.
	r, err = newReport(&pkguri.PkgURI{
		URI: "https://github.com/ahshum/release",
	})
	assert.Nil(t, err)
	if assert.NotNil(t, r) {
		if assert.True(t, r.Next()) {
			tag, err := r.Tag()
			assert.Nil(t, err)
			assert.Equal(t, "v1.0.0", tag.Value)
		}
		if assert.True(t, r.Next()) {
			tag, err := r.Tag()
			assert.Nil(t, err)
			assert.Equal(t, "v2.0.0", tag.Value)
			assert.True(t, tag.IsLatest)
		}
		assert.False(t, r.Next())
		r.Close()
	}

	// Lightweight tags.
	r, err = newReport(&pkguri.PkgURI{
		URI: "https://github.com/ahshum/release-lw",
	})
	assert.Nil(t, err)
	if assert.NotNil(t, r) {
		if assert.True(t, r.Next()) {
			tag, err := r.Tag()
			assert.Nil(t, err)
			assert.Equal(t, "v1.0.0", tag.Value)
		}
		if assert.True(t, r.Next()) {
			tag, err := r.Tag()
			assert.Nil(t, err)
			assert.Equal(t, "v2.0.0", tag.Value)
			assert.True(t, tag.IsLatest)
		}
		assert.False(t, r.Next())
		r.Close()
	}
}
