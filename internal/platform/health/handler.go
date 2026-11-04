package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// SystemCheckResponse defines a structured JSON format for status checking
type SystemCheckResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

type SystemCheckHandler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewSystemCheckHandler(db *pgxpool.Pool, redis *redis.Client) *SystemCheckHandler {
	return &SystemCheckHandler{db: db, redis: redis}
}

func (h *SystemCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(SystemCheckResponse{
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *SystemCheckHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	checks := make(map[string]string)
	isReady := true

	// PostgreSQL connection check
	if err := h.db.Ping(ctx); err != nil {
		slog.Error("postgres health check failed", "err", err)

		checks["postgres"] = "DOWN"
		isReady = false
	} else {
		checks["postgres"] = "UP"
	}

	// Redis connection check
	if err := h.redis.Ping(ctx).Err(); err != nil {
		slog.Error("redis health check failed", "err", err)

		checks["redis"] = "DOWN"
		isReady = false
	} else {
		checks["redis"] = "UP"
	}

	w.Header().Set("Content-Type", "application/json")

	if !isReady {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(SystemCheckResponse{
			Status:    "DOWN",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Checks:    checks,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SystemCheckResponse{
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	})
}
