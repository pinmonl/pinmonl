package main

import (
	"fmt"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/pkger"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	pkger.Include("/migrations")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
