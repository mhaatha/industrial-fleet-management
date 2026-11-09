package health

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type fakeSystemCheckRepository struct {
	pgErr, redisErr error
}

func (f *fakeSystemCheckRepository) CheckPostgresql(ctx context.Context) error { return f.pgErr }

func (f *fakeSystemCheckRepository) CheckRedis(ctx context.Context) error { return f.redisErr }

func TestReadinessCheck(t *testing.T) {
	boom := errors.New("boom")

	tests := []struct {
		name     string
		fakeRepo fakeSystemCheckRepository
		checks   map[string]string
		isReady  bool
	}{
		{
			name:     "all services are running",
			fakeRepo: fakeSystemCheckRepository{},
			checks:   map[string]string{"postgres": "UP", "redis": "UP"},
			isReady:  true,
		},
		{
			name:     "postgresql is down",
			fakeRepo: fakeSystemCheckRepository{pgErr: boom},
			checks:   map[string]string{"postgres": "DOWN", "redis": "UP"},
			isReady:  false,
		},
		{
			name:     "redis is down",
			fakeRepo: fakeSystemCheckRepository{redisErr: boom},
			checks:   map[string]string{"postgres": "UP", "redis": "DOWN"},
			isReady:  false,
		},
		{
			name:     "all services are down",
			fakeRepo: fakeSystemCheckRepository{pgErr: boom, redisErr: boom},
			checks:   map[string]string{"postgres": "DOWN", "redis": "DOWN"},
			isReady:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSystemCheckService(&tt.fakeRepo)

			checks, isReady := service.ReadinessCheck(context.Background())

			if !reflect.DeepEqual(checks, tt.checks) {
				t.Errorf("expected checks to be equal, but got %v want %v", checks, tt.checks)
			}

			if isReady != tt.isReady {
				t.Errorf("expected isReady to be equal, but got %v, want %v", isReady, tt.isReady)
			}
		})
	}
}
