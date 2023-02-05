package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/HardDie/event_tracker/internal/db"
	"github.com/HardDie/event_tracker/internal/logger"
)

type Config struct {
	DB             db.DBConfig
	Port           string
	PwdMaxAttempts int
	PwdBlockTime   int
	RequestTimeout int
}

func Get() *Config {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error.Printf("failed to load env vars: %s", err)
		}
	}

	return &Config{
		DB: db.DBConfig{
			Host:     getEnv("DB_HOST", "db"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "event_tracker"),
			Password: getEnv("DB_PWD", "event_tracker"),
			Database: getEnv("DB_DB", "db"),
		},
		Port:           getEnv("PORT", ":8080"),
		PwdMaxAttempts: getEnvAsInt("PWD_MAX_ATTEMPTS", 5),
		PwdBlockTime:   getEnvAsInt("PWD_BLOCK_TIME", 24),
		RequestTimeout: getEnvAsInt("REQUEST_TIMEOUT", 3),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func getEnvAsInt(key string, defaultValue int) int {
	value := getEnv(key, "")
	if v, e := strconv.Atoi(value); e == nil {
		return v
	}
	return defaultValue
}
