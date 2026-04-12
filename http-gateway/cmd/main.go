package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/chishkin-afk/posted/http-gateway/docs"
	"github.com/chishkin-afk/posted/http-gateway/internal/app"
)

// @title Posted HTTP Gateway API
// @version 1.0.0
// @description API gateway for authentication and posts management services.
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type JWT token.

func main() {
	app, cleanup, err := app.New()
	if err != nil {
		slog.Error("failed to setup app", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cleanup()

	go func() {
		if err := app.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server",
				slog.String("error", err.Error()),
			)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
}
