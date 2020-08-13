package docker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pinmonl/pinmonl/model"
)

// Patterns of semver.
var (
	ChannelPattern      = regexp.MustCompile("^[a-zA-Z]+$")
	SemverPattern       = regexp.MustCompile("^v?(\\d+)(\\.\\d+)?(\\.\\d+)?")
	StrictSemverPattern = regexp.MustCompile("^v?(\\d+)(\\.\\d+)(\\.\\d+)")
)

var tagDigestStdout = ioutil.Discard

// tagBucket stores different kinds of tags.
type tagBucket struct {
	latest   *model.Stat
	semvers  model.StatList
	channels model.StatList
	orphans  model.StatList
	children map[*model.Stat]model.StatList
}

// newTagBucket creates tagBucket from tags.
func newTagBucket(tags []*model.Stat) (*tagBucket, error) {
	var (
		latest         *model.Stat
		latestManifest = make(map[string]int)
	)
	// Retrieve "latest" tag.
	if got := model.StatList(tags).GetValue("latest"); len(got) > 0 {
		latest = got[0]
	} else {
		return nil, errors.New("docker: do not have latest tag")
	}
	// Prepare for latest tag comparsion.
	for _, digest := range digestsFromTag(latest) {
		latestManifest[digest]++
	}

	// Define for calculation.
	var (
		latestTag      *model.Stat
		fallbackTag    *model.Stat
		latestVer      = &semver.Version{}
		fallbackVer    = &semver.Version{}
		digestToTags   = make(map[string]model.StatList)
		tagToDigests   = tagDigestSet{}
		tagChildren    = tagMapSet{}
		tagAllChildren = tagMapSet{}
		tagParents     = tagMapSet{}
		tagIdentical   = tagIdenticalSet{}
		tagOrphans     = tagSet{}
		tagChannels    = tagSet{}
		tagSemvers     = tagSet{}
	)

	// 1st loop to search for latest tag.
	for _, tag := range tags {
		if tver, err := semver.NewVersion(tag.Value); err == nil {
			isLatest := (tagCompare(tag, latest) == tagIsIdentical)
			if isLatest && tver.GreaterThan(latestVer) {
				latestTag = tag
				latestVer = tver
			}

			isStable := (tver.Metadata() == "" && tver.Prerelease() == "")
			if isStable && tver.GreaterThan(fallbackVer) {
				fallbackTag = tag
				fallbackVer = tver
			}
		}

		for _, digest := range digestsFromTag(tag) {
			digestToTags[digest] = append(digestToTags[digest], tag)
			tagToDigests.Add(tag, digest)
		}

		if ChannelPattern.MatchString(tag.Value) {
			tagChannels.Add(tag)
		}
	}

	// 2nd loop for digest comparing.
	for _, tag := range tags {
		digests := tagToDigests.GetSlice(tag)

		if len(digests) == 0 {
			tagOrphans.Add(tag)
			continue
		}

		intersectCount := 0

		for _, digest := range digests {
			for _, intersect := range digestToTags[digest] {
				if tag == intersect {
					continue
				}

				intersectCount++
				intersectDigests := tagToDigests.GetSlice(intersect)

				compared := digestsCompare(digests, intersectDigests)

				fmt.Fprintf(tagDigestStdout, "digest compare: %s(%d) vs %s(%d) == %d\n",
					tag.Value,
					len(digests),
					intersect.Value,
					len(intersectDigests),
					compared)

				switch compared {
				case tagIsIdentical:
					digestKey := tagDigestKey(digests)
					tagIdentical.Add(digestKey, tag)
					tagIdentical.Add(digestKey, intersect)
					tagAllChildren.Add(tag, intersect)
					tagAllChildren.Add(intersect, tag)
				case tagIsSuperset:
					tagAllChildren.Add(tag, intersect)
					tagChildren.Add(tag, intersect)
					tagParents.Add(intersect, tag)
				case tagIsSubset:
					tagAllChildren.Add(intersect, tag)
					tagChildren.Add(intersect, tag)
					tagParents.Add(tag, intersect)
				}
			}
		}

		if intersectCount == 0 {
			tagChildren[tag] = tagSet{}
		}
	}

	// 3rd loop for fill in semvers.

	// Copy orphans.
	for tag := range tagOrphans {
		tagSemvers.Set(tag, 1)
	}

	// Hide channels.
	for tag := range tagChannels {
		tagSemvers.Set(tag, 0)
	}

	// Choose the tag with shortest name and
	// strict semver from each identical group.
	for _, tagset := range tagIdentical {
		var (
			shortest *model.Stat
			slen     int
		)
		for tag := range tagset {
			if val, has := tagSemvers[tag]; has && val == 0 {
				continue
			}

			switch tlen := len(tag.Value); {
			// Is children
			case len(tagParents[tag]) > 0:
				tagSemvers.Set(tag, 0)
			// Is not strict semver
			case !StrictSemverPattern.MatchString(tag.Value):
				tagSemvers.Set(tag, 0)
			// Shortest is not set.
			case shortest == nil:
				tagSemvers.Set(tag, 1)
				shortest, slen = tag, tlen
			// Tag name is shorter.
			case tlen < slen:
				tagSemvers.Set(shortest, 0)
				tagSemvers.Set(tag, 1)
				shortest, slen = tag, tlen
			default:
				tagSemvers.Set(tag, 0)
			}
		}
	}

	// Copy parents.
	for tag, childset := range tagChildren {
		for child := range childset {
			if _, skip := tagSemvers[child]; skip {
				continue
			}
			tagSemvers.Set(child, 0)
		}

		if len(tagParents[tag]) > 0 {
			tagSemvers.Set(tag, 0)
		} else if _, skip := tagSemvers[tag]; skip {
			continue
		} else {
			tagSemvers.Set(tag, 1)
		}
	}

	// Latest tag fallback.
	if latestTag == nil {
		latestTag = fallbackTag
	}

	return &tagBucket{
		latest:   latestTag,
		semvers:  tagSemvers.Slice(),
		channels: tagChannels.Slice(),
		orphans:  tagOrphans.Slice(),
		children: tagAllChildren.SliceMap(),
	}, nil
}

// tagIntersection describes the intersection between tags.
type tagIntersection int

// intersection types.
const (
	tagIsSuperset tagIntersection = iota
	tagIsSubset
	tagIsIdentical
	tagIsUndefined
)

// tagCompare reports the intersection result between tags.
func tagCompare(src, other *model.Stat) tagIntersection {
	return digestsCompare(
		digestsFromTag(src),
		digestsFromTag(other),
	)
}

// digestsCompare reports the intersection result between digests.
func digestsCompare(src, other []string) tagIntersection {
	checks := make(map[string]int)
	for _, digest := range src {
		checks[digest] = 0
	}
	intersectCount := 0
	for _, digest := range other {
		if _, intersect := checks[digest]; intersect {
			intersectCount++
			checks[digest]++
		}
	}

	var (
		slen = len(src)
		olen = len(other)
	)
	if intersectCount == slen && intersectCount == olen {
		return tagIsIdentical
	}
	if intersectCount == slen && intersectCount < olen {
		return tagIsSubset
	}
	if intersectCount == olen && intersectCount < slen {
		return tagIsSuperset
	}
	return tagIsUndefined
}

// digestsFromTag extracts digests from the manifests of tag.
func digestsFromTag(tag *model.Stat) []string {
	if tag.Substats == nil || len((*tag.Substats).GetKind(model.ManifestStat)) == 0 {
		return nil
	}

	var (
		manifests = (*tag.Substats).GetKind(model.ManifestStat)
		checks    = make(map[string]int)
		out       = make([]string, 0)
	)
	for _, image := range manifests {
		digest := image.Checksum
		if digest == "" {
			continue
		}
		if _, has := checks[digest]; has {
			continue
		}

		out = append(out, digest)
		checks[digest]++
	}

	return out
}

func tagDigestKey(digests []string) string {
	sort.Sort(sort.StringSlice(digests))
	return strings.Join(digests, ",")
}

// TagBySemver sorts tags.
type TagBySemver model.StatList

// Len reports the length.
func (tl TagBySemver) Len() int { return len(tl) }

// Swap swaps elements.
func (tl TagBySemver) Swap(i, j int) { tl[i], tl[j] = tl[j], tl[i] }

// Less reports whether item i has a semver less than item j.
// Always returns true if item i is not a semver and
// false if item j is not a semver.
func (tl TagBySemver) Less(i, j int) bool {
	iv, err := newVersion(tl[i].Value)
	if err != nil {
		// If error occurs, sort to top.
		return true
	}
	ij, err := newVersion(tl[j].Value)
	if err != nil {
		// If error occurs, sort to top.
		return false
	}
	return iv.Compare(ij) < 0
}

// newVersion creates semver and tries fallback version string if fails.
func newVersion(s string) (ver *semver.Version, outerr error) {
	if v, err := semver.NewVersion(s); err == nil {
		ver = v
	} else if SemverPattern.MatchString(s) {
		s = SemverPattern.ReplaceAllString(s, "$0-")
		if v, err := semver.NewVersion(s); err == nil {
			ver = v
		} else {
			outerr = err
		}
		return
	} else {
		outerr = err
	}

	return
}

// tagSet stores unique tags.
type tagSet map[*model.Stat]int

// Add adds tag.
func (t tagSet) Add(tag *model.Stat) {
	t[tag]++
}

// Set sets tag with value.
func (t tagSet) Set(tag *model.Stat, value int) {
	t[tag] = value
}

// Del removes tag.
func (t tagSet) Del(tag *model.Stat) {
	delete(t, tag)
}

// Slice returns model.StatList.
func (t tagSet) Slice() model.StatList {
	out := model.StatList{}
	for tag, value := range t {
		if value > 0 {
			out = append(out, tag)
		}
	}
	return out
}

// tagMapSet stores unique child tags of tag.
type tagMapSet map[*model.Stat]tagSet

// Add adds child to tag.
func (t tagMapSet) Add(tag, child *model.Stat) {
	if _, inited := t[tag]; !inited {
		t[tag] = tagSet{}
	}

	t[tag].Add(child)
}

// Set sets child of tag.
func (t tagMapSet) Set(tag, child *model.Stat) {
	t[tag] = tagSet{child: 1}
}

// Get gets the first child of tag.
func (t tagMapSet) Get(tag *model.Stat) *model.Stat {
	for value := range t[tag] {
		return value
	}
	return nil
}

// Del removes child from tag.
func (t tagMapSet) Del(tag, value *model.Stat) {
	if _, has := t[tag]; !has {
		return
	}

	delete(t[tag], value)
}

// SliceMap returns tag-children paired map.
func (t tagMapSet) SliceMap() map[*model.Stat]model.StatList {
	out := make(map[*model.Stat]model.StatList)
	for tag, tagset := range t {
		out[tag] = tagset.Slice()
	}
	return out
}

// tagIdenticalSet stores unique tags of string key.
type tagIdenticalSet map[string]tagSet

// Add adds tag to key.
func (t tagIdenticalSet) Add(key string, tag *model.Stat) {
	if _, inited := t[key]; !inited {
		t[key] = tagSet{}
	}

	t[key].Add(tag)
}

// Set sets tag to key.
func (t tagIdenticalSet) Set(key string, tag *model.Stat) {
	t[key] = tagSet{tag: 1}
}

// Del removes tag from key.
func (t tagIdenticalSet) Del(key string, tag *model.Stat) {
	if _, has := t[key]; !has {
		return
	}

	delete(t[key], tag)
}

// tagDigestSet stores unique digest string of tag.
type tagDigestSet map[*model.Stat]map[string]int

// Add adds digest to tag.
func (t tagDigestSet) Add(tag *model.Stat, digest string) {
	if _, inited := t[tag]; !inited {
		t[tag] = make(map[string]int)
	}

	t[tag][digest] = 1
}

// GetSlice returns digests of tag.
func (t tagDigestSet) GetSlice(tag *model.Stat) []string {
	digestMap := t[tag]
	out := make([]string, len(digestMap))

	i := 0
	for digest := range digestMap {
		out[i] = digest
		i++
	}
	return out
}
