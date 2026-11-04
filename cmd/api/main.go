package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/config"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/health"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Logger init
	config.LogInit()

	// Config init
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("godotenv fails to load .env file", "err", err)
		os.Exit(1)
	}

	// Main router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	ctx := context.Background()

	pgURL := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", cfg.DBUser, cfg.DBPassword, cfg.DBName)
	pool, err := pgxpool.New(ctx, pgURL)
	if err != nil {
		log.Fatalf("Gagal terhubung ke PostgreSQL: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	systemCheckHandler := health.NewSystemCheckHandler(pool, rdb)
	health.SystemCheckRouter(r, systemCheckHandler)

	http.ListenAndServe(":3000", r)
}
