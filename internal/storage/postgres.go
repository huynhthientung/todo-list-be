package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // register pgx with database/sql

	"todo-list-be/internal/config"
)

// OpenPostgres opens a PostgreSQL-backed *sql.DB with basic pooling settings.
func OpenPostgres(ctx context.Context, cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseURL())
	if err != nil {
		return nil, err
	}

	// Conservative defaults for small services.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
