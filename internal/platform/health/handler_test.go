package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type fakeSystemCheckService struct {
	checks  map[string]string
	isReady bool
}

func (f *fakeSystemCheckService) ReadinessCheck(ctx context.Context) (map[string]string, bool) {
	return f.checks, f.isReady
}

func TestHandlerHealthCheck(t *testing.T) {
	handler := NewSystemCheckHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.HealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		Message string `json:"message"`
		Data    struct {
			Status    string            `json:"status"`
			Timestamp string            `json:"timestamp"`
			Checks    map[string]string `json:"checks"`
		} `json:"data"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Errorf("response body is not a valid JSON, err: %v", err)
	}

	if body.Message != "OK" {
		t.Errorf("got message %v want %v", body.Message, "OK")
	}

	if body.Data.Status != "UP" {
		t.Errorf("got status %v want %v", body.Data.Status, "UP")
	}

	if body.Data.Checks != nil {
		t.Errorf("body.Data.Checks should be nil, but got %v", body.Data.Checks)
	}

	if _, err := time.Parse(time.RFC3339, body.Data.Timestamp); err != nil {
		t.Errorf("timestamp is not RFC3339 format: %v", err)
	}
}

func TestHandlerReadinessCheck(t *testing.T) {
	tests := []struct {
		name        string
		fakeService fakeSystemCheckService
		statusCode  int
	}{
		{
			name: "all services are running",
			fakeService: fakeSystemCheckService{
				map[string]string{
					"postgres": "UP",
					"redis":    "UP",
				},
				true,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "postgresql is down",
			fakeService: fakeSystemCheckService{
				map[string]string{
					"postgres": "DOWN",
					"redis":    "UP",
				},
				false,
			},
			statusCode: http.StatusServiceUnavailable,
		},
		{
			name: "redis is down",
			fakeService: fakeSystemCheckService{
				map[string]string{
					"postgres": "UP",
					"redis":    "DOWN",
				},
				false,
			},
			statusCode: http.StatusServiceUnavailable,
		},
		{
			name: "all services are down",
			fakeService: fakeSystemCheckService{
				map[string]string{
					"postgres": "DOWN",
					"redis":    "DOWN",
				},
				false,
			},
			statusCode: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		var successBody struct {
			Message string `json:"message"`
			Data    struct {
				Status    string            `json:"status"`
				Timestamp string            `json:"timestamp"`
				Checks    map[string]string `json:"checks"`
			} `json:"data"`
		}

		var failedBody struct {
			Message string `json:"message"`
			Details struct {
				Status    string            `json:"status"`
				Timestamp string            `json:"timestamp"`
				Checks    map[string]string `json:"checks"`
			} `json:"details"`
		}

		t.Run(tt.name, func(t *testing.T) {
			handler := NewSystemCheckHandler(&tt.fakeService)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			handler.ReadinessCheck(rec, req)

			if rec.Code != tt.statusCode {
				t.Errorf("got %d, want %d", rec.Code, tt.statusCode)
			}

			switch tt.statusCode {
			case http.StatusOK:
				if err := json.NewDecoder(rec.Body).Decode(&successBody); err != nil {
					t.Errorf("response body is not a valid JSON, err: %v", err)
				}

				if successBody.Message != "OK" {
					t.Errorf("got message %v want %v", successBody.Message, "OK")
				}

				if successBody.Data.Status != "UP" {
					t.Errorf("got status %v want %v", successBody.Data.Status, "UP")
				}

				if !reflect.DeepEqual(successBody.Data.Checks, tt.fakeService.checks) {
					t.Errorf("expected checks to be equal, but got %v want %v", successBody.Data.Checks, tt.fakeService.checks)
				}

				if _, err := time.Parse(time.RFC3339, successBody.Data.Timestamp); err != nil {
					t.Errorf("timestamp is not RFC3339 format: %v", err)
				}
			case http.StatusServiceUnavailable:
				if err := json.NewDecoder(rec.Body).Decode(&failedBody); err != nil {
					t.Errorf("response body is not a valid JSON, err: %v", err)
				}

				if failedBody.Message != "Some services are unavailable." {
					t.Errorf("got message %v want %v", successBody.Message, "Some services are unavailable.")
				}

				if failedBody.Details.Status != "DOWN" {
					t.Errorf("got status %v want %v", failedBody.Details.Status, "DOWN")
				}

				if !reflect.DeepEqual(failedBody.Details.Checks, tt.fakeService.checks) {
					t.Errorf("expected checks to be equal, but got %v want %v", failedBody.Details.Checks, tt.fakeService.checks)
				}

				if _, err := time.Parse(time.RFC3339, failedBody.Details.Timestamp); err != nil {
					t.Errorf("timestamp is not RFC3339 format: %v", err)
				}
			default:
				t.Errorf("should't return %v status code", tt.statusCode)
			}
		})
	}
}
