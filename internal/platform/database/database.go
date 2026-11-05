package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/config"
)

func ConnectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	pgURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := pgxpool.New(context.Background(), pgURL)
	if err != nil {
		return nil, err
	}

	// Check DB connection
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	slog.Info("connected to the database")
	return db, nil
}
