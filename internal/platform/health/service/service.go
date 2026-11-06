package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/mhaatha/industrial-fleet-management/internal/platform/health/repository"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/web"
)

// SystemCheckService is a contract interface for system checks.
type SystemCheckService interface {
	ReadinessCheck(ctx context.Context) web.SystemCheckResponse
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

func (s *systemCheckService) ReadinessCheck(ctx context.Context) web.SystemCheckResponse {
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

	// Call check Redis repository
	if err := s.SystemCheckRepository.CheckRedis(ctx); err != nil {
		slog.Error("redis health check failed", "err", err)

		checks["redis"] = "DOWN"
		isReady = false
	} else {
		checks["redis"] = "UP"
	}

	if !isReady {
		status = "DOWN"
	}

	return web.SystemCheckResponse{
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}
}
