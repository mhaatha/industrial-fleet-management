package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/mhaatha/industrial-fleet-management/internal/platform/health/repository"
	systemCheckWeb "github.com/mhaatha/industrial-fleet-management/internal/platform/health/web"
	"github.com/redis/go-redis/v9"
)

// SystemCheckService is a contract interface for system checks.
type SystemCheckService interface {
	ReadinessCheck(ctx context.Context) (systemCheckWeb.SystemCheckResponse, error)
}

// systemCheckService implements SystemCheckService interface
type systemCheckService struct {
	SystemCheckRepository repository.SystemCheckRepository
}

func NewSystemCheckService(systemCheckRepository repository.SystemCheckRepository) SystemCheckService {
	return &systemCheckService{
		SystemCheckRepository: systemCheckRepository,
	}
}

func (s *systemCheckService) ReadinessCheck(ctx context.Context) (systemCheckWeb.SystemCheckResponse, error) {
	checks := make(map[string]string)
	isReady := true
	status := "UP"

	// Call check Postgres repository
	if err := s.SystemCheckRepository.CheckPostgresql(ctx); err != nil {
		slog.Error("postgres health check failed", "err", err)

		checks["postgres"] = "DOWN"
		isReady = false
	} else {
		checks["postgres"] = "UP"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	// Call check Redis repository
	if err := s.SystemCheckRepository.CheckRedis(ctx, rdb); err != nil {
		slog.Error("redis health check failed", "err", err)

		checks["redis"] = "DOWN"
		isReady = false
	} else {
		checks["redis"] = "UP"
	}

	if !isReady {
		status = "DOWN"
	}

	return systemCheckWeb.SystemCheckResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}, nil
}
