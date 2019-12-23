package database

import (
	"bytes"
	"path"
	"strings"
)

// MigrationSource defines the interface for source accessing.
type MigrationSource interface {
	List() MigrationList
}

// PackrBox defines the interface for using packr.
type PackrBox interface {
	List() []string
	Find(name string) ([]byte, error)
}

// PackrMigrationSource loads files from packr.Box.
type PackrMigrationSource struct {
	Box PackrBox
	Dir string
}

// List returns list of migration.
func (s PackrMigrationSource) List() MigrationList {
	prefix := ""
	if dir := path.Clean(s.Dir); dir != "." {
		prefix = dir + "/"
	}

	ms := make([]Migration, 0)

	for _, f := range s.Box.List() {
		if !strings.HasPrefix(f, prefix) {
			continue
		}
		if path.Ext(f) == ".sql" {
			c, _ := s.Box.Find(f)
			n := strings.TrimSuffix(path.Base(f), ".sql")
			m, err := parseMigration(n, bytes.NewBuffer(c))
			if err == nil {
				ms = append(ms, m)
			}
		}
	}

	return ms
}
