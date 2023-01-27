package db

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/glebarez/go-sqlite"
)

type DB struct {
	DB *sql.DB
}

func Get(dbpath string) (*DB, error) {
	flags := []string{
		"_pragma=foreign_keys(1)",
		"_pragma=busy_timeout(3000)",
		"_pragma=journal_mode(WAL)",

		//"_fk=true",
		//"_timeout=5000",
		//"_journal=WAL",

		"_pragma=journal_size_limit(1073741824)", // 1Gb limit for -wal file
	}

	db, err := sql.Open("sqlite3", dbpath+"?"+strings.Join(flags, "&"))
	if err != nil {
		return nil, err
	}
	//_, err = db.Exec("PRAGMA journal_size_limit=1073741824")
	//if err != nil {
	//	return nil, err
	//}
	return &DB{
		DB: db,
	}, nil
}
