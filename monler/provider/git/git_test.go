package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderPing(t *testing.T) {
	p, err := NewProvider()
	assert.Nil(t, err)
	assert.Nil(t, p.Ping("https://github.com/ahshum/empty"))
	assert.Nil(t, p.Ping("https://github.com/ahshum/release"))
	assert.NotNil(t, p.Ping("https://github.com/ahshum/not-exist"))
}

func TestRepoAnalyze(t *testing.T) {
	var (
		repo   *Repo
		report *Report
		err    error
	)

	// Empty repo.
	repo, err = newRepo("https://github.com/ahshum/empty")
	assert.Nil(t, err)
	if assert.NotNil(t, repo) {
		report, err = repo.analyze()
		assert.Nil(t, err)
		assert.NotNil(t, report)
		repo.Close()
	}

	// Annotated tags.
	repo, err = newRepo("https://github.com/ahshum/release")
	assert.Nil(t, err)
	if assert.NotNil(t, repo) {
		report, err = repo.analyze()
		assert.Nil(t, err)
		if assert.True(t, report.Next()) {
			tag, err := report.Tag()
			assert.Nil(t, err)
			assert.Equal(t, "v1.0.0", tag.Value)
		}
		if assert.True(t, report.Next()) {
			tag, err := report.Tag()
			assert.Nil(t, err)
			assert.Equal(t, "v2.0.0", tag.Value)
			assert.True(t, tag.IsLatest)
		}
		assert.False(t, report.Next())
		repo.Close()
	}

	// Lightweight tags.
	repo, err = newRepo("https://github.com/ahshum/release-lw")
	assert.Nil(t, err)
	if assert.NotNil(t, repo) {
		report, err = repo.analyze()
		assert.Nil(t, err)
		if assert.True(t, report.Next()) {
			tag, err := report.Tag()
			assert.Nil(t, err)
			assert.Equal(t, "v1.0.0", tag.Value)
		}
		if assert.True(t, report.Next()) {
			tag, err := report.Tag()
			assert.Nil(t, err)
			assert.Equal(t, "v2.0.0", tag.Value)
			assert.True(t, tag.IsLatest)
		}
		assert.False(t, report.Next())
		repo.Close()
	}
}
