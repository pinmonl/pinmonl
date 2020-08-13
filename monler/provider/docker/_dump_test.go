package docker

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pinmonl/pinmonl/model"
)

func testRepoName() string {
	return "drone/drone"
	// return "rclone/rclone"
	// return "linuxserver/jellyfin"
	// return "library/php"
}

func testRepoFile() string {
	return strings.ReplaceAll(testRepoName(), "/", "_")
}

func testFilepath(filename string) string {
	path, _ := filepath.Abs("./" + filename)
	return path
}

func TestDumpRepo(t *testing.T) {
	client := &Client{client: &http.Client{}}

	// Set digest compare output
	f, _ := os.Create(testFilepath("dump-digestcompare"))
	tagDigestStdout = f
	defer f.Close()

	// Fetch
	var tags []*model.Stat
	t.Run("fetch", func(t *testing.T) {
		start := time.Now()
		tags, _ = fetchAllTags(client, testRepoName())
		t.Log("took", time.Since(start))
	})

	var bucket *tagBucket
	t.Run("bucket", func(t *testing.T) {
		start := time.Now()
		bucket, _ = newTagBucket(tags)
		t.Log("took", time.Since(start))
		t.Log("latest tag is :", bucket.latest.Value)
	})

	t.Run("dump", func(t *testing.T) {
		start := time.Now()
		testDumpStatList(tags, testFilepath("dump-tags"))
		testDumpStatList(bucket.semvers, testFilepath("dump-semvers"))
		testDumpStatList(bucket.orphans, testFilepath("dump-orphans"))
		testDumpStatList(bucket.channels, testFilepath("dump-channels"))
		testDumpTagChildren(bucket.children, testFilepath("dump-children"))
		t.Log("took", time.Since(start))
	})
}

func testDumpStatList(list model.StatList, filename string) {
	file, _ := os.Create(filename)
	defer file.Close()

	for _, stat := range list {
		images := make([]string, 0)
		manifests := (*stat.Substats).GetKind(model.ManifestStat)
		for _, image := range manifests {
			images = append(images, fmt.Sprintf("%s__%s", image.Value, image.Checksum))
		}
		fmt.Fprintf(file, "%s : %s\n", stat.Value, strings.Join(images, ", "))
	}
}

func testDumpTagChildren(childrenMap map[*model.Stat]model.StatList, filename string) {
	file, _ := os.Create(filename)
	defer file.Close()

	for tag, children := range childrenMap {
		aliases := make([]string, 0)
		for _, child := range children {
			aliases = append(aliases, child.Value)
		}

		fmt.Fprintf(file, "%s : %s\n", tag.Value, strings.Join(aliases, ", "))
	}
}
