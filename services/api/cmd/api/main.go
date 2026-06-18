package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	httpadapter "github.com/sutad-p/oh-my-trading/services/api/internal/adapters/http"
	"github.com/sutad-p/oh-my-trading/services/api/internal/adapters/mockdata"
	"github.com/sutad-p/oh-my-trading/services/api/internal/adapters/postgres"
	"github.com/sutad-p/oh-my-trading/services/api/internal/application/marketdata"
	appmt5 "github.com/sutad-p/oh-my-trading/services/api/internal/application/mt5"
	"github.com/sutad-p/oh-my-trading/services/api/internal/platform/config"
	"github.com/sutad-p/oh-my-trading/services/api/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Environment)

	var handler http.Handler
	if cfg.APIMockMode {
		log.Info("api mock mode enabled", slog.String("mode", "mock"))
		handler = httpadapter.NewRouter(
			httpadapter.WithSymbolService(mockdata.NewSymbolService(uuid.NewString)),
			httpadapter.WithCandleService(mockdata.NewCandleService()),
			httpadapter.WithIndicatorService(mockdata.NewIndicatorService()),
		)
	} else {
		db, err := sql.Open("pgx", cfg.DatabaseURL)
		if err != nil {
			log.Error("open database", slog.String("error", err.Error()))
			os.Exit(1)
		}
		defer db.Close()

		if err := postgres.RunMigrations(context.Background(), db, filepath.Join("migrations")); err != nil {
			log.Error("run database migrations", slog.String("error", err.Error()))
			os.Exit(1)
		}

		symbolRepo := postgres.NewSymbolRepository(db)
		candleRepo := postgres.NewCandleRepository(db)
		mt5Repo := postgres.NewMT5Repository(db)

		handler = httpadapter.NewRouter(
			httpadapter.WithSymbolService(marketdata.NewSymbolService(symbolRepo, uuid.NewString)),
			httpadapter.WithCandleService(marketdata.NewCandleService(candleRepo)),
			httpadapter.WithIndicatorService(marketdata.NewIndicatorService(candleRepo, time.Now)),
			httpadapter.WithMT5Service(appmt5.NewService(mt5Repo, symbolRepo, candleRepo)),
		)
	}

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: handler,
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
