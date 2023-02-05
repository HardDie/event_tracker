package db

import (
	"context"
	"fmt"

	"github.com/HardDie/godb/v2"
	_ "github.com/lib/pq"

	"github.com/HardDie/event_tracker/internal/logger"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type DB struct {
	DB *godb.DBO
}

func Get(cfg DBConfig) (*DB, error) {
	conf := godb.PostgresConnectionConfig{
		ConnectionConfig: godb.ConnectionConfig{
			Host:                   cfg.Host,
			Port:                   cfg.Port,
			Name:                   cfg.Database,
			User:                   cfg.User,
			Password:               cfg.Password,
			MaxConnections:         50,
			ConnectionIdleLifetime: 15,
		},
		SSLMode: "disable",
	}

	dbo, err := godb.DBO{
		Options: godb.Options{
			//Debug:  true,
			Logger: logger.Debug,
		},
		Connection: &conf,
	}.Init()
	if err != nil {
		return nil, fmt.Errorf("error open connection to db: %w", err)
	}

	return &DB{
		DB: dbo,
	}, nil
}

func (db *DB) BeginTx(ctx context.Context) (*godb.SqlTx, error) {
	return db.DB.BeginContext(ctx)
}
func (db *DB) EndTx(tx *godb.SqlTx, err error) error {
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
