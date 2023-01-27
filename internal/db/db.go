package db

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/glebarez/go-sqlite"

	"github.com/HardDie/event_tracker/internal/logger"
)

type DB struct {
	DB *sql.DB
}

func Get(dbpath string) (*DB, error) {
	flags := []string{
		"_fk=true",      // "_pragma=foreign_keys(1)",
		"_timeout=5000", // "_pragma=busy_timeout(3000)",
		"_journal=WAL",  // "_pragma=journal_mode(WAL)",
	}

	db, err := sql.Open("sqlite3", dbpath+"?"+strings.Join(flags, "&"))
	if err != nil {
		return nil, err
	}

	// "_pragma=journal_size_limit(1073741824)", // 1Gb limit for -wal file
	_, err = db.Exec("PRAGMA journal_size_limit=1073741824")
	if err != nil {
		logger.Error.Println("error pragma journal_size_limit:", err.Error())
	}

	return &DB{
		DB: db,
	}, nil
}
