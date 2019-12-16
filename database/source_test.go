package database

import (
	"testing"

	"github.com/gobuffalo/packr/v2"
)

func TestPackrSource(t *testing.T) {
	testCases := []struct {
		name   string
		boxDir string
		dir    string
		passFn func(MigrationList) bool
	}{
		{
			name:   "standard",
			boxDir: "./testdata/source/standard",
			dir:    "",
			passFn: func(list MigrationList) bool {
				return len(list) == 2 &&
					list[0].Name == "0001_one" &&
					list[1].Name == "0002_two"
			},
		},
		{
			name:   "standard dir",
			boxDir: "./testdata/source",
			dir:    "standard",
			passFn: func(list MigrationList) bool {
				return len(list) == 2 &&
					list[0].Name == "0001_one" &&
					list[1].Name == "0002_two"
			},
		},
		{
			name:   "discardline",
			boxDir: "./testdata/source/discardline",
			dir:    "",
			passFn: func(list MigrationList) bool {
				return len(list) == 2 &&
					list[0].Name == "0001_one" &&
					list[1].Name == "0002_two"
			},
		},
		{
			name:   "discardline dir",
			boxDir: "./testdata/source",
			dir:    "discardline",
			passFn: func(list MigrationList) bool {
				return len(list) == 2 &&
					list[0].Name == "0001_one" &&
					list[1].Name == "0002_two"
			},
		},
	}

	for _, tc := range testCases {
		src := PackrMigrationSource{
			Box: packr.New(tc.boxDir, tc.boxDir),
			Dir: tc.dir,
		}
		passes := tc.passFn(src.List())
		if !passes {
			t.Errorf("case %q fails", tc.name)
		}
	}
}
