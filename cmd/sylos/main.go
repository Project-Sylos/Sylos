package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codeberg.org/Sylos/Sylos-API/pkg/app"
	"codeberg.org/Sylos/Sylos-API/pkg/config"
	"codeberg.org/Sylos/Sylos/internal/browser"
	sylosui "codeberg.org/Sylos/Sylos/internal/ui"
)

const (
	healthPollAttempts = 15
	healthPollInterval = 200 * time.Millisecond
	startupTimeout     = 120 * time.Second
)

func main() {
	noBrowser := flag.Bool("no-browser", false, "do not open the web UI in a browser")
	apiOnly := flag.Bool("api-only", false, "serve the API only; do not serve the embedded web UI")
	port := flag.Int("port", 0, "HTTP port override (default 8086)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if *port > 0 {
		cfg.HTTP.Port = *port
	}

	var staticHandler http.Handler
	if !*apiOnly {
		staticHandler, err = sylosui.StaticHandler()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to initialize UI static handler: %v\n", err)
			os.Exit(1)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ready := make(chan struct{})
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx, app.Options{
			Config:        cfg,
			StaticHandler: staticHandler,
			Ready:         ready,
		})
	}()

	if !*apiOnly && !*noBrowser {
		if err := waitUntilReady(ready, errCh, startupTimeout); err != nil {
			fmt.Fprintf(os.Stderr, "server failed to start: %v\n", err)
			stop()
			os.Exit(1)
		}

		if err := waitForHealthy(cfg.HTTP.Port); err != nil {
			fmt.Fprintf(os.Stderr, "server health check failed: %v\n", err)
			stop()
			os.Exit(1)
		}

		url := fmt.Sprintf("http://localhost:%d/", cfg.HTTP.Port)
		if err := browser.OpenURL(url); err != nil {
			fmt.Fprintf(os.Stderr, "failed to open browser (visit %s manually): %v\n", url, err)
		}
	}

	if err := <-errCh; err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func waitForHealthy(port int) error {
	url := fmt.Sprintf("http://127.0.0.1:%d/health", port)
	client := &http.Client{Timeout: 2 * time.Second}

	for attempt := 0; attempt < healthPollAttempts; attempt++ {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(healthPollInterval)
	}

	return fmt.Errorf("timed out waiting for %s", url)
}

func waitUntilReady(ready <-chan struct{}, errCh <-chan error, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-ready:
		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("server exited during startup: %w", err)
		}
		return fmt.Errorf("server exited before becoming ready")
	case <-timer.C:
		return fmt.Errorf("timed out after %s waiting for HTTP listener", timeout)
	}
}
