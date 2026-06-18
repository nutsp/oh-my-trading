package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	amqp "github.com/rabbitmq/amqp091-go"
	marketdataadapter "github.com/sutad-p/oh-my-trading/services/api/internal/adapters/marketdata"
	"github.com/sutad-p/oh-my-trading/services/api/internal/adapters/postgres"
	app "github.com/sutad-p/oh-my-trading/services/api/internal/application/marketdata"
	"github.com/sutad-p/oh-my-trading/services/api/internal/platform/config"
	"github.com/sutad-p/oh-my-trading/services/api/internal/platform/logger"
)

const syncQueueName = "market-data-sync"

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Environment)

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

	rabbitURL := envString("OMT_RABBITMQ_URL", "amqp://omt:omt_local_password@localhost:5672/")
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Error("connect rabbitmq", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Error("open rabbitmq channel", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer channel.Close()

	if _, err := channel.QueueDeclare(syncQueueName, true, false, false, false, nil); err != nil {
		log.Error("declare sync queue", slog.String("error", err.Error()))
		os.Exit(1)
	}

	deliveries, err := channel.Consume(syncQueueName, "oh-my-trading-worker", false, false, false, false, nil)
	if err != nil {
		log.Error("consume sync queue", slog.String("error", err.Error()))
		os.Exit(1)
	}

	syncer := app.NewSyncService(marketdataadapter.NewSyntheticProvider(), postgres.NewCandleRepository(db))
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	log.Info("market data worker started", slog.String("queue", syncQueueName))
	for {
		select {
		case <-shutdown:
			log.Info("market data worker stopped")
			return
		case delivery := <-deliveries:
			var request app.SyncRequest
			if err := json.Unmarshal(delivery.Body, &request); err != nil {
				log.Error("decode sync request", slog.String("error", err.Error()))
				_ = delivery.Nack(false, false)
				continue
			}
			if err := syncer.SyncCandles(context.Background(), request); err != nil {
				log.Error("sync candles", slog.String("error", err.Error()))
				_ = delivery.Nack(false, true)
				continue
			}
			_ = delivery.Ack(false)
		}
	}
}

func envString(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
