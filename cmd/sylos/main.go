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

	"codeberg.org/Sylos/Sylos/internal/browser"
	sylosui "codeberg.org/Sylos/Sylos/internal/ui"
	"codeberg.org/Sylos/Sylos-API/pkg/app"
	"codeberg.org/Sylos/Sylos-API/pkg/config"
)

const (
	healthPollAttempts = 30
	healthPollInterval = 200 * time.Millisecond
)

func main() {
	noBrowser := flag.Bool("no-browser", false, "do not open the web UI in a browser")
	apiOnly := flag.Bool("api-only", false, "serve the API only; do not serve the embedded web UI")
	configPath := flag.String("config", "config.yaml", "path to config file")
	port := flag.Int("port", 0, "HTTP port override (default from config)")
	flag.Parse()

	if *configPath != "" {
		if err := os.Setenv("SYLOS_CONFIG_PATH", *configPath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to set config path: %v\n", err)
			os.Exit(1)
		}
	}

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

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx, app.Options{
			Config:        cfg,
			StaticHandler: staticHandler,
		})
	}()

	if !*apiOnly && !*noBrowser {
		if err := waitForHealthy(cfg.HTTP.Port); err != nil {
			fmt.Fprintf(os.Stderr, "server health check failed: %v\n", err)
			stop()
			os.Exit(1)
		}

		url := fmt.Sprintf("http://127.0.0.1:%d/", cfg.HTTP.Port)
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
