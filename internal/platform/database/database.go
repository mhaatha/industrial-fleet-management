package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/config"
)

func ConnectDB(cfg *config.Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	config.ConnConfig.User = cfg.DBUser
	config.ConnConfig.Password = cfg.DBPassword
	config.ConnConfig.Host = cfg.DBHost
	config.ConnConfig.Database = cfg.DBName
	config.ConnConfig.Port = uint16(cfg.DBPort)

	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnIdleTime = 30 * time.Minute
	config.MaxConnLifetime = 1 * time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Check DB connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	slog.Info("connected to the database")
	return pool, nil
}
