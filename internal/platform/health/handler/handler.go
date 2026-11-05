package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/mhaatha/industrial-fleet-management/internal/platform/health/service"
	systemCheckWeb "github.com/mhaatha/industrial-fleet-management/internal/platform/health/web"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/helper"
	platformWeb "github.com/mhaatha/industrial-fleet-management/internal/platform/web"
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
	helper.WriteToResponseBody(w, http.StatusOK, platformWeb.WebSuccessResponse{
		Message: "OK",
		Data: systemCheckWeb.SystemCheckResponse{
			Status:    "UP",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *SystemCheckHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Error ignored because we don't implemented the error handling yet
	systemCheckResponse, _ := h.SystemCheckService.ReadinessCheck(ctx)

	// ⚠️ IMPLEMENT ERROR HANDLING IF ONE OF THE SERVICE ARE UNAVAILABLE ⚠️

	helper.WriteToResponseBody(w, http.StatusOK, platformWeb.WebSuccessResponse{
		Message: "OK",
		Data:    systemCheckResponse,
	})
}
