package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	AdminToken   string
	UserToken    string
	LogLevel     string
	MaxDBConns   int
	MigrationsDir string
}


func Load() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:password@db:5432/avito_mvp?sslmode=disable"),
		AdminToken:   getEnv("ADMIN_TOKEN", ""),
		UserToken:    getEnv("USER_TOKEN", ""),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		MigrationsDir: getEnv("MIGRATIONS_DIR", "/migrations"),
	}

	if v := getEnv("MAX_DB_CONNS", "20"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.MaxDBConns = n
		} else {
			c.MaxDBConns = 20
		}
	} else {
		c.MaxDBConns = 20
	}
	return c, nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
