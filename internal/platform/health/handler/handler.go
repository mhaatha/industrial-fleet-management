package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/mhaatha/industrial-fleet-management/internal/platform/health/service"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/helper"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/web"
)

type SystemCheckHandler struct {
	SystemCheckService service.SystemCheckService
}

func NewSystemCheckHandler(systemCheckService service.SystemCheckService) *SystemCheckHandler {
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

	systemCheckResponse := h.SystemCheckService.ReadinessCheck(ctx)

	if systemCheckResponse.Checks["postgres"] == "DOWN" || systemCheckResponse.Checks["redis"] == "DOWN" {
		helper.WriteToResponseBody(w, http.StatusServiceUnavailable, web.WebFailedResponse{
			Message: "Some service are unavailable, please try again later.",
			Details: systemCheckResponse,
		})
		return
	}

	helper.WriteToResponseBody(w, http.StatusOK, web.WebSuccessResponse{
		Message: "OK",
		Data:    systemCheckResponse,
	})
}
