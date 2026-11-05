package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/config"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/database"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/health/handler"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/health/repository"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/health/service"
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

	// Database init
	db, err := database.ConnectDB(cfg)
	if err != nil {
		slog.Error("failed connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// Main router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// System check init
	systemCheckRepository := repository.NewSystemCheckRepository(db)
	systemCheckService := service.NewSystemCheckService(systemCheckRepository)
	systemCheckHandler := handler.NewSystemCheckHandler(systemCheckService)
	handler.SystemCheckRouter(r, systemCheckHandler)

	// Server init
	server := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// ⚠️ GRACEFUL SHUTDOWN WILL BE IMPLEMENTED ⚠️

	slog.Info("starting server on :" + cfg.ServerPort)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server error", "err", err)
		os.Exit(1)
	}
}
