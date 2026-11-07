package health

import "github.com/go-chi/chi/v5"

func SystemCheckRouter(r chi.Router, h *SystemCheckHandler) {
	r.Route("/health", func(r chi.Router) {
		r.Get("/", h.HealthCheck)
	})

	r.Route("/ready", func(r chi.Router) {
		r.Get("/", h.ReadinessCheck)
	})
}
