package helper

import (
	"log/slog"
	"net/http"

	"github.com/mhaatha/industrial-fleet-management/internal/platform/apperror"
	"github.com/mhaatha/industrial-fleet-management/internal/platform/web"
)

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	appErr, ok := apperror.As(err)
	if !ok {
		// unknown error is wrapped as internal error
		appErr = apperror.InternalError(err)
	}

	// Log the original error for logging purposes
	if appErr.HTTPStatus >= 500 {
		slog.Error("request failed",
			"method", r.Method,
			"path", r.URL.Path,
			"code", appErr.Code,
			"error", appErr.Error(),
		)
	} else {
		slog.Warn("request rejected",
			"path", r.URL.Path,
			"code", appErr.Code,
			"message", appErr.Message,
		)
	}

	WriteToResponseBody(w, appErr.HTTPStatus, web.WebFailedResponse{Message: appErr.Message, Details: appErr.Details})
}
