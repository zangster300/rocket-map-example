package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"rocket-map-example/config"
	"rocket-map-example/web/resources"

	"github.com/starfederation/datastar-go/datastar"
)

func SetupRoutes(ctx context.Context, mux *http.ServeMux) (err error) {
	setupReloadRoutes(ctx, mux)

	mux.Handle("GET /static/", resources.Handler())

	if err := errors.Join(
		setupIndexRoute(ctx, mux),
		setupMapRoute(ctx, mux),
	); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	return nil
}

func setupReloadRoutes(ctx context.Context, mux *http.ServeMux) {
	reloadChan := make(chan struct{}, 1)
	var hotReloadOnce sync.Once
	mux.HandleFunc("GET /reload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		reload := func() { sse.ExecuteScript("window.location.reload()") }
		hotReloadOnce.Do(reload)
		select {
		case <-reloadChan:
			reload()
		case <-r.Context().Done():
		case <-ctx.Done():
		}
	})

	if config.Global.Environment == config.Dev {
		mux.HandleFunc("/hotreload", func(w http.ResponseWriter, r *http.Request) {
			select {
			case reloadChan <- struct{}{}:
			default:
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	}
}
