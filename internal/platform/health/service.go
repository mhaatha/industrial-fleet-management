package health

import (
	"context"
)

// SystemCheckService is a contract interface for system checks.
type SystemCheckService interface {
	ReadinessCheck(ctx context.Context) (map[string]string, bool)
}

// systemCheckService implements SystemCheckService interface
type systemCheckService struct {
	SystemCheckRepository SystemCheckRepository
}

func NewSystemCheckService(systemCheckRepository SystemCheckRepository) SystemCheckService {
	return &systemCheckService{
		SystemCheckRepository: systemCheckRepository,
	}
}

func (s *systemCheckService) ReadinessCheck(ctx context.Context) (map[string]string, bool) {
	checks := make(map[string]string)
	isReady := true

	// Call check Postgres repository
	if err := s.SystemCheckRepository.CheckPostgresql(ctx); err != nil {
		checks["postgres"] = "DOWN"
		isReady = false
	} else {
		checks["postgres"] = "UP"
	}

	// Call check Redis repository
	if err := s.SystemCheckRepository.CheckRedis(ctx); err != nil {
		checks["redis"] = "DOWN"
		isReady = false
	} else {
		checks["redis"] = "UP"
	}

	return checks, isReady
}
