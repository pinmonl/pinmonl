package database

import (
	"sort"
	"testing"

	"github.com/gobuffalo/packr/v2"
	_ "github.com/mattn/go-sqlite3"
)

func TestMigrationInstall(t *testing.T) {
	db, _ := Open("sqlite3", "file::memory:?cache=shared")
	box := packr.New("./testdata/source/install", "./testdata/source/install")
	src := PackrMigrationSource{Box: box}
	mp := NewMigrationPlan(db.DB, src)

	if mp.HasMigrationTable() {
		t.Error("migration table should not exist.")
	}

	mp.Install()
	if !mp.HasMigrationTable() {
		t.Error("migration table should have installed.")
	}
	if r := mp.Records(); len(r) > 0 {
		t.Error("migration should be empty.")
	}
}

func TestMigrationRun(t *testing.T) {
	db, _ := Open("sqlite3", "file::memory:?cache=shared")
	box := packr.New("./testdata/source/install", "./testdata/source/install")
	src := PackrMigrationSource{Box: box}
	mp := NewMigrationPlan(db.DB, src)

	if err := mp.Install(); err != nil {
		t.Error("migration install fails")
	}

	records := []string{"0001_one", "0002_two", "0003_three"}
	tables := []string{"one", "two", "three"}

	testCases := []struct {
		name    string
		run     func(*MigrationPlan)
		tables  []string
		records []string
	}{
		{
			name:    "Up",
			run:     func(m *MigrationPlan) { m.Up() },
			tables:  tables,
			records: records,
		},
		{
			name:    "Up 1",
			run:     func(m *MigrationPlan) { m.UpTo(1) },
			tables:  tables[:1],
			records: records[:1],
		},
		{
			name:    "Up 2",
			run:     func(m *MigrationPlan) { m.UpTo(2) },
			tables:  tables[:2],
			records: records[:2],
		},
		{
			name:    "Down",
			run:     func(m *MigrationPlan) { m.Up(); m.Down() },
			tables:  nil,
			records: nil,
		},
		{
			name:    "Down 1",
			run:     func(m *MigrationPlan) { m.Up(); m.DownTo(1) },
			tables:  tables[:2],
			records: records[:2],
		},
		{
			name:    "Down 2",
			run:     func(m *MigrationPlan) { m.Up(); m.DownTo(2) },
			tables:  tables[:1],
			records: records[:1],
		},
	}

	for _, tc := range testCases {
		mp.Down()
		tc.run(mp)
		for i, r := range mp.Records() {
			if r.Name != tc.records[i] {
				t.Errorf(
					"case %q record does not match, wants %s, got %s",
					tc.name,
					tc.records[i],
					r.Name,
				)
			}
		}
		for _, table := range tc.tables {
			if !mp.hasTable(table) {
				t.Errorf("case %q table %q does not exist", tc.name, table)
			}
		}
	}
}

func TestMigrationListSort(t *testing.T) {
	testCases := []struct {
		name  string
		list  MigrationList
		wants []string
	}{
		{
			name: "fixed width",
			list: []Migration{
				Migration{Name: "0001_one"},
				Migration{Name: "0003_three"},
				Migration{Name: "0002_two"},
				Migration{Name: "0004_four"},
			},
			wants: []string{"0001_one", "0002_two", "0003_three", "0004_four"},
		},
		{
			name: "number",
			list: []Migration{
				Migration{Name: "3_three"},
				Migration{Name: "1_one"},
				Migration{Name: "4_four"},
				Migration{Name: "2_two"},
			},
			wants: []string{"1_one", "2_two", "3_three", "4_four"},
		},
		{
			name: "timestamp",
			list: []Migration{
				Migration{Name: "20191216010213_two"},
				Migration{Name: "20191216010203_one"},
				Migration{Name: "20191216010223_three"},
				Migration{Name: "20191216010233_four"},
			},
			wants: []string{"20191216010203_one", "20191216010213_two", "20191216010223_three", "20191216010233_four"},
		},
		{
			name: "same prefix",
			list: []Migration{
				Migration{Name: "0001_c"},
				Migration{Name: "0001_b"},
				Migration{Name: "0001_a"},
				Migration{Name: "0001_d"},
			},
			wants: []string{"0001_a", "0001_b", "0001_c", "0001_d"},
		},
	}

	for _, tc := range testCases {
		list := tc.list
		sort.Sort(list)
		for i, m := range list {
			if m.Name != tc.wants[i] {
				t.Errorf("case %q sorts incorrectly at %d", tc.name, i)
			}
		}
	}
}

func TestVersionString(t *testing.T) {
	testCases := []struct {
		name    string
		version string
		wants   string
	}{
		{
			name:    "empty",
			version: "",
			wants:   "",
		},
		{
			name:    "incorrect format",
			version: "abc_123123",
			wants:   "",
		},
		{
			name:    "fixed width",
			version: "0001_name",
			wants:   "0001",
		},
		{
			name:    "number",
			version: "1_name",
			wants:   "1",
		},
		{
			name:    "timestamp",
			version: "20191216010203_name",
			wants:   "20191216010203",
		},
	}

	for _, tc := range testCases {
		got := versionFrom(tc.version)
		if got != tc.wants {
			t.Errorf("case %q fails", tc.name)
		}
	}
}
