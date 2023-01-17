package db

import (
	"database/sql"
	"strings"

	_ "github.com/glebarez/go-sqlite"
)

type DB struct {
	DB *sql.DB
}

func Get(dbpath string) (*DB, error) {
	flags := []string{
		"_pragma=foreign_keys(1)",
	}

	db, err := sql.Open("sqlite", dbpath+"?"+strings.Join(flags, "&"))
	if err != nil {
		return nil, err
	}
	return &DB{
		DB: db,
	}, nil
}
