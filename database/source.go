package database

import (
	"bytes"
	"path"
	"strings"
)

// MigrationSource defines the interface for source accessing
type MigrationSource interface {
	List() MigrationList
}

// PackrBox defines the interface for using packr
type PackrBox interface {
	List() []string
	Find(name string) ([]byte, error)
}

// PackrMigrationSource loads files from packr.Box
type PackrMigrationSource struct {
	Box PackrBox
	Dir string
}

// List returns list of migration
func (s PackrMigrationSource) List() MigrationList {
	prefix := ""
	if dir := path.Clean(s.Dir); dir != "." {
		prefix = dir + "/"
	}

	migrations := make([]Migration, 0)

	for _, file := range s.Box.List() {
		if !strings.HasPrefix(file, prefix) {
			continue
		}
		if path.Ext(file) == ".sql" {
			content, _ := s.Box.Find(file)
			name := strings.TrimSuffix(path.Base(file), ".sql")
			if m, err := parseMigration(name, bytes.NewBuffer(content)); err == nil {
				migrations = append(migrations, m)
			}
		}
	}

	return migrations
}
