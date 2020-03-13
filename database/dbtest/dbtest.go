package dbtest

import (
	"os"

	// Import database drivers.
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/migrations"
)

// Open opens a database connection for testing.
func Open() (*database.DB, error) {
	driver := "sqlite3"
	dsn := ":memory:?cache=shared"
	if os.Getenv("DB_DRIVER") != "" {
		driver = os.Getenv("DB_DRIVER")
		dsn = os.Getenv("DB_DSN")
	}

	db, err := database.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	mgsrc := database.PackrMigrationSource{
		Box: migrations.PackrBox(),
		Dir: db.DriverName(),
	}
	dbm := database.NewMigrationPlan(db.DB, mgsrc)
	err = dbm.Install()
	if err != nil {
		return nil, err
	}
	err = dbm.Up()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Reset clears the database.
func Reset(db *database.DB) {
	db.Exec("DELETE FROM images")
	db.Exec("DELETE FROM jobs")
	db.Exec("DELETE FROM monls")
	db.Exec("DELETE FROM pinls")
	db.Exec("DELETE FROM pinmonls")
	db.Exec("DELETE FROM pkgs")
	db.Exec("DELETE FROM share_tags")
	db.Exec("DELETE FROM shares")
	db.Exec("DELETE FROM stats")
	db.Exec("DELETE FROM taggables")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM users")
}

// Close closes the database connection.
func Close(db *database.DB) error {
	return db.Close()
}
