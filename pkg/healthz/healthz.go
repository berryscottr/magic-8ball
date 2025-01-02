package healthz

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// Start initializes and starts the health check HTTP server
func (h *HealthCheckServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	h.server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Info().Msg("health check endpoint running on :8080")
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("failed to start health check server")
		}
	}()
}

// Close gracefully shuts down the health check HTTP server
func (h *HealthCheckServer) Close() {
	log.Info().Msg("shutting down health check server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		log.Err(err).Msg("error shutting down health check server")
	} else {
		log.Info().Msg("health check server stopped successfully")
	}
}
