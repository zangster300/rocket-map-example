package routes

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"rocket-map-example/templating"
)

func setupIndexRoute(ctx context.Context, mux *http.ServeMux) error {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if err := templating.Index().Render(r.Context(), w); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("GET /stream", func(w http.ResponseWriter, r *http.Request) {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		// sse := datastar.NewSSE(w, r)

		for {
			select {
			case <-ctx.Done():
				slog.Debug("server closed connection")
				return

			case <-r.Context().Done():
				slog.Debug("client closed connection")
				return

			case <-ticker.C:
				// stream map marker events
			}
		}
	})

	return nil
}
