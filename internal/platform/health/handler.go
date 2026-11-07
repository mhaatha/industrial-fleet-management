package health

import (
	"context"
	"net/http"
	"time"

	"github.com/mhaatha/industrial-fleet-management/internal/platform/helper"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/web"
)

type SystemCheckHandler struct {
	SystemCheckService SystemCheckService
}

func NewSystemCheckHandler(systemCheckService SystemCheckService) *SystemCheckHandler {
	return &SystemCheckHandler{
		SystemCheckService: systemCheckService,
	}
}

func (h *SystemCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	helper.WriteToResponseBody(w, http.StatusOK, web.WebSuccessResponse{
		Message: "OK",
		Data: web.SystemCheckResponse{
			Status:    "UP",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *SystemCheckHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	checks, ready := h.SystemCheckService.ReadinessCheck(ctx)

	status, httpStatus, message := "UP", http.StatusOK, "OK"
	if !ready {
		status, httpStatus, message = "DOWN", http.StatusServiceUnavailable, "Some services are unavailable."

		helper.WriteToResponseBody(w, httpStatus, web.WebFailedResponse{
			Message: message,
			Details: web.SystemCheckResponse{
				Status:    status,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Checks:    checks,
			},
		})
	} else {
		helper.WriteToResponseBody(w, httpStatus, web.WebSuccessResponse{
			Message: message,
			Data: web.SystemCheckResponse{
				Status:    status,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Checks:    checks,
			},
		})
	}
}
