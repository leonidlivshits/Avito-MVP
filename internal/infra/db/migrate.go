package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool, migrationsDir string) error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is empty")
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("open database for migrations: %w", err)
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping database for migrations: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create migrate driver: %w", err)
	}

	sourceURL := "file://" + migrationsDir
	migr, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance (source=%s): %w", sourceURL, err)
	}

	if err := migr.Up(); err != nil {
		if err == migrate.ErrNoChange {
			srcErr, dbErr := migr.Close()
			if srcErr != nil || dbErr != nil {
				return fmt.Errorf("migrate close after no change: sourceErr=%v dbErr=%v", srcErr, dbErr)
			}
			return nil
		}
		srcErr, dbErr := migr.Close()
		if srcErr != nil || dbErr != nil {
			return fmt.Errorf("migrate up failed: %v; additionally close errors: sourceErr=%v dbErr=%v", err, srcErr, dbErr)
		}
		return fmt.Errorf("migrate up failed: %w", err)
	}

	srcErr, dbErr := migr.Close()
	if srcErr != nil || dbErr != nil {
		return fmt.Errorf("close migrate instance: sourceErr=%v dbErr=%v", srcErr, dbErr)
	}

	return nil
}
