package db

import (
	"context"
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
		"_timeout=5000", // "_pragma=busy_timeout(5000)",
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

func (db *DB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, nil)
}
func (db *DB) EndTx(tx *sql.Tx, err error) error {
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			logger.Error.Println("error rollback tx:", err.Error())
			return err
		}
		return nil
	}

	err = tx.Commit()
	if err != nil {
		logger.Error.Println("error commit tx:", err.Error())
		return err
	}
	return nil
}
