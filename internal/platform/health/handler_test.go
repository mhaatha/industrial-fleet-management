package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeSystemCheckService struct{}

func (f *fakeSystemCheckService) ReadinessCheck(ctx context.Context) (map[string]string, bool) {
	return map[string]string{}, false
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
