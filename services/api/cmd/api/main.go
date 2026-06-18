package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpadapter "github.com/sutad-p/oh-my-trading/services/api/internal/adapters/http"
	"github.com/sutad-p/oh-my-trading/services/api/internal/platform/config"
	"github.com/sutad-p/oh-my-trading/services/api/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Environment)

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpadapter.NewRouter(),
	}

	errs := make(chan error, 1)
	go func() {
		log.Info("api server starting", slog.String("addr", cfg.HTTPAddr))
		errs <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errs:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("api server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case sig := <-shutdown:
		log.Info("api server shutting down", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("api server shutdown failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	log.Info("api server stopped")
}
