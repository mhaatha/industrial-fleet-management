package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/config"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/database"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/health"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Logger init
	config.LogInit()

	// Config init
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("error occures while calling LoadConfig", "err", err)
		os.Exit(1)
	}

	// Database init
	db, err := database.ConnectDB(cfg)
	if err != nil {
		slog.Error("failed connect to database", "err", err)
		os.Exit(1)
	}

	// Redis init
	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Main router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// System check init
	systemCheckRepository := health.NewSystemCheckRepository(db, rdb)
	systemCheckService := health.NewSystemCheckService(systemCheckRepository)
	systemCheckHandler := health.NewSystemCheckHandler(systemCheckService)
	health.SystemCheckRouter(r, systemCheckHandler)

	// Server init
	server := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Server goroutine
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting server on :" + cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Block until it receives SIGINT or SIGTERM
	select {
	case err := <-serverErr:
		slog.Error("HTTP server error", "err", err)
		os.Exit(1)
	case <-ctx.Done():
		slog.Info("signal received, begin graceful shutdown")
	}

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("incomplete shutdown", "err", err)
		server.Close()
	}

	db.Close()
	if err := rdb.Close(); err != nil {
		slog.Error("failed to close redis conn", "err", err)
	}
	slog.Info("graceful shutdown success")
}
