package version

import "github.com/Masterminds/semver"

var (
	// VersionString is the plain version string.
	VersionString string = "0.3.0"
	// Version is parsed semver object from string.
	Version *semver.Version = semver.MustParse(VersionString)
)
