package database

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// MigrationPlan controls the database version based on the database driver.
type MigrationPlan struct {
	db  *sqlx.DB
	src MigrationSource

	tableName string
}

// NewMigrationPlan creates a migration plan.
func NewMigrationPlan(db *sqlx.DB, src MigrationSource) *MigrationPlan {
	return &MigrationPlan{
		db:  db,
		src: src,
	}
}

// TableName returns the defined table name for migration records.
func (mp *MigrationPlan) TableName() string {
	if mp.tableName == "" {
		return "migrations"
	}
	return mp.tableName
}

// SetTableName sets the table name for migration records.
func (mp *MigrationPlan) SetTableName(name string) {
	mp.tableName = name
}

// Install creates the table of migration records.
func (mp *MigrationPlan) Install() error {
	if mp.HasMigrationTable() {
		return nil
	}

	var err error
	switch mp.db.DriverName() {
	case "sqlite3":
		_, err = mp.db.Exec(fmt.Sprintf(`
			CREATE TABLE %s (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR(255) UNIQUE
			)`, mp.TableName()))
	case "postgres":
		_, err = mp.db.Exec(fmt.Sprintf(`
			CREATE TABLE %s (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) UNIQUE
			)`, mp.TableName()))
	}
	return err
}

// HasMigrationTable reports whether the migration table exists or not.
func (mp *MigrationPlan) HasMigrationTable() bool {
	return mp.hasTable(mp.TableName())
}

func (mp *MigrationPlan) hasTable(name string) bool {
	_, err := mp.db.Queryx(fmt.Sprintf(`SELECT 1 FROM %s LIMIT 1`, name))
	return err == nil
}

// Up runs UpTo with limit disabled.
func (mp *MigrationPlan) Up() error {
	return mp.UpTo(-1)
}

// UpTo runs up statements by limit.
func (mp *MigrationPlan) UpTo(limit int) error {
	return mp.runMigration(dirUp, limit)
}

// Down runs DownTo with limit disabled.
func (mp *MigrationPlan) Down() error {
	return mp.DownTo(-1)
}

// DownTo runs down statements by limit.
func (mp *MigrationPlan) DownTo(limit int) error {
	return mp.runMigration(dirDown, limit)
}

func (mp *MigrationPlan) runMigration(dir direction, limit int) error {
	var (
		ms    = mp.src.List()
		start = 0
		rstmt string
		rs    = mp.Records()
	)

	switch dir {
	case dirUp:
		rstmt = fmt.Sprintf("INSERT INTO %s (name) VALUES (:name)", mp.TableName())
		sort.Sort(ms)
		if len(rs) > 0 {
			start = start + 1
		}
	case dirDown:
		if len(rs) == 0 {
			return nil
		}
		rstmt = fmt.Sprintf("DELETE FROM %s WHERE name = :name", mp.TableName())
		sort.Sort(sort.Reverse(ms))
	}
	if len(rs) > 0 {
		latest := rs[len(rs)-1]
		start = start + ms.IndexOf(latest.Name)
	}

	ms2 := ms[start:]
	if limit > 0 {
		end := limit
		ms2 = ms2[:end]
	}

	var err error
	tx, _ := mp.db.Beginx()
migrate:
	for _, m := range ms2 {
		args := map[string]interface{}{"name": m.Name}
		_, err = tx.NamedExec(rstmt, args)
		if err != nil {
			break migrate
		}
		for _, stmt := range m.Stmts(dir) {
			_, err = tx.Exec(stmt)
			if err != nil {
				break migrate
			}
		}
	}

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Records gets all the migration records from database.
func (mp *MigrationPlan) Records() []MigrationRecord {
	rows, err := mp.db.Queryx(fmt.Sprintf(`
		SELECT id, name FROM %s
		ORDER BY id ASC, name ASC`, mp.TableName()))
	if err != nil {
		return nil
	}

	mrs := make([]MigrationRecord, 0)
	for rows.Next() {
		var mr MigrationRecord
		err = rows.StructScan(&mr)
		if err != nil {
			return nil
		}
		mrs = append(mrs, mr)
	}
	return mrs
}

var migrationNameRegex = regexp.MustCompile(`^(\d+).*$`)

func nameMatches(name string) []string {
	return migrationNameRegex.FindStringSubmatch(name)
}

func versionFrom(name string) string {
	nm := nameMatches(name)
	if len(nm) <= 1 {
		return ""
	}
	return nm[1]
}

func versionInt(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}

// Migration stores the up and down sql statements.
type Migration struct {
	Name string
	Up   []string
	Down []string
}

// Version returns the version part from the name.
func (m Migration) Version() string {
	return versionFrom(m.Name)
}

// VersionInt returns the version in integer.
func (m Migration) VersionInt() int {
	return versionInt(m.Version())
}

// Stmts returns the migration statements by direction.
func (m Migration) Stmts(dir direction) []string {
	switch dir {
	case dirUp:
		return m.Up
	case dirDown:
		return m.Down
	default:
		return nil
	}
}

// MigrationRecord defines the structure of storing a ran migration.
type MigrationRecord struct {
	ID   int
	Name string
}

// Version returns the version part from the name.
func (m MigrationRecord) Version() string {
	return versionFrom(m.Name)
}

// VersionInt returns the version in integer.
func (m MigrationRecord) VersionInt() int {
	return versionInt(m.Version())
}

// MigrationList is the array of Migration.
type MigrationList []Migration

// Len is the number of elements in the collection.
func (ml MigrationList) Len() int { return len(ml) }

// Swap swaps the elements with indexes i and j.
func (ml MigrationList) Swap(i, j int) { ml[i], ml[j] = ml[j], ml[i] }

// Less reports whether the element with
// index i should sort before the element with index j.
func (ml MigrationList) Less(i, j int) bool {
	if ml[i].VersionInt() < ml[j].VersionInt() {
		return true
	}
	return strings.Compare(ml[i].Name, ml[j].Name) == -1
}

// IndexOf finds the index of the element by name.
func (ml MigrationList) IndexOf(name string) int {
	for i, m := range ml {
		if m.Name == name {
			return i
		}
	}
	return -1
}
