package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestSymbolRepositoryCreatesAndListsSymbols(t *testing.T) {
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

	repo := NewSymbolRepository(db)
	created, err := repo.CreateSymbol(ctx, market.Symbol{
		ID:         "018f4f8a-0000-7000-9000-000000000101",
		Code:       "BTCUSD",
		Market:     "forex",
		BaseAsset:  "BTC",
		QuoteAsset: "USD",
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("CreateSymbol returned error: %v", err)
	}
	if created.Code != "BTCUSD" {
		t.Fatalf("created.Code = %q, want BTCUSD", created.Code)
	}

	symbols, err := repo.ListSymbols(ctx)
	if err != nil {
		t.Fatalf("ListSymbols returned error: %v", err)
	}
	if len(symbols) != 2 {
		t.Fatalf("len(symbols) = %d, want 2", len(symbols))
	}
	if symbols[0].Code != "BTCUSD" || symbols[1].Code != "XAUUSD" {
		t.Fatalf("symbol order = %q, %q; want BTCUSD, XAUUSD", symbols[0].Code, symbols[1].Code)
	}
}
