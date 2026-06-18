package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestCandleRepositoryUpsertsAndListsCandles(t *testing.T) {
	databaseURL := os.Getenv("OMT_TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://omt:omt_local_password@localhost:15432/oh_my_trading?sslmode=disable"
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		t.Skipf("postgres integration database is unavailable: %v", err)
	}

	resetMigrationState(t, ctx, db)
	if err := RunMigrations(ctx, db, filepath.Join("..", "..", "..", "migrations")); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	symbolRepo := NewSymbolRepository(db)
	symbol, err := symbolRepo.CreateSymbol(ctx, market.Symbol{
		ID:         "018f4f8a-0000-7000-9000-000000000201",
		Code:       "XAUUSD",
		Market:     "forex",
		BaseAsset:  "XAU",
		QuoteAsset: "USD",
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("CreateSymbol returned error: %v", err)
	}

	repo := NewCandleRepository(db)
	first := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	second := first.Add(time.Hour)

	err = repo.UpsertCandles(ctx, []market.Candle{
		{
			SymbolID:  symbol.ID,
			Timeframe: "1h",
			Timestamp: first,
			Open:      2320.1,
			High:      2328.3,
			Low:       2318.6,
			Close:     2325.5,
			Volume:    12345,
		},
		{
			SymbolID:  symbol.ID,
			Timeframe: "1h",
			Timestamp: second,
			Open:      2325.5,
			High:      2330.0,
			Low:       2322.0,
			Close:     2328.0,
			Volume:    23456,
		},
	})
	if err != nil {
		t.Fatalf("UpsertCandles returned error: %v", err)
	}

	err = repo.UpsertCandles(ctx, []market.Candle{{
		SymbolID:  symbol.ID,
		Timeframe: "1h",
		Timestamp: first,
		Open:      2320.1,
		High:      2329.0,
		Low:       2318.6,
		Close:     2327.7,
		Volume:    99999,
	}})
	if err != nil {
		t.Fatalf("second UpsertCandles returned error: %v", err)
	}

	candles, err := repo.ListCandles(ctx, market.CandleQuery{
		SymbolCode: "XAUUSD",
		Timeframe:  "1h",
		From:       first.Add(-time.Minute),
		To:         second.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("ListCandles returned error: %v", err)
	}
	if len(candles) != 2 {
		t.Fatalf("len(candles) = %d, want 2", len(candles))
	}
	if !candles[0].Timestamp.Equal(first) {
		t.Fatalf("candles[0].Timestamp = %s, want %s", candles[0].Timestamp, first)
	}
	if candles[0].Close != 2327.7 {
		t.Fatalf("candles[0].Close = %f, want 2327.7", candles[0].Close)
	}
	if !candles[1].Timestamp.Equal(second) {
		t.Fatalf("candles[1].Timestamp = %s, want %s", candles[1].Timestamp, second)
	}
}
