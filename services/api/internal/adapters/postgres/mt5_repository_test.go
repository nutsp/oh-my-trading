package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestMT5RepositoryPersistsReadOnlySnapshots(t *testing.T) {
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

	repo := NewMT5Repository(db)
	heartbeatTime := time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC)
	err = repo.SaveHeartbeat(ctx, MT5Heartbeat{
		BridgeID:     "local-mt5",
		Terminal:     "MetaTrader 5",
		AccountLogin: "12345678",
		Server:       "Broker-Demo",
		Status:       "healthy",
		SentAt:       heartbeatTime,
	})
	if err != nil {
		t.Fatalf("SaveHeartbeat returned error: %v", err)
	}

	updatedHeartbeatTime := heartbeatTime.Add(30 * time.Second)
	err = repo.SaveHeartbeat(ctx, MT5Heartbeat{
		BridgeID:     "local-mt5",
		Terminal:     "MetaTrader 5",
		AccountLogin: "12345678",
		Server:       "Broker-Demo",
		Status:       "stale",
		LastError:    "poll timeout",
		SentAt:       updatedHeartbeatTime,
	})
	if err != nil {
		t.Fatalf("second SaveHeartbeat returned error: %v", err)
	}

	latestHeartbeat, err := repo.LatestHeartbeat(ctx, "local-mt5")
	if err != nil {
		t.Fatalf("LatestHeartbeat returned error: %v", err)
	}
	if latestHeartbeat.Status != "stale" {
		t.Fatalf("latestHeartbeat.Status = %q, want stale", latestHeartbeat.Status)
	}
	if latestHeartbeat.LastError != "poll timeout" {
		t.Fatalf("latestHeartbeat.LastError = %q, want poll timeout", latestHeartbeat.LastError)
	}
	if !latestHeartbeat.SentAt.Equal(updatedHeartbeatTime) {
		t.Fatalf("latestHeartbeat.SentAt = %s, want %s", latestHeartbeat.SentAt, updatedHeartbeatTime)
	}

	tickTime := time.Date(2026, 6, 18, 10, 1, 0, 0, time.UTC)
	ticks := []MT5Tick{
		{
			Symbol: "XAUUSD",
			Bid:    2325.42,
			Ask:    2325.62,
			Last:   2325.52,
			Volume: 12,
			Time:   tickTime,
		},
		{
			Symbol: "XAUUSD",
			Bid:    2325.50,
			Ask:    2325.72,
			Last:   2325.60,
			Volume: 14,
			Time:   tickTime.Add(time.Second),
		},
	}
	if err := repo.SaveTicks(ctx, ticks); err != nil {
		t.Fatalf("SaveTicks returned error: %v", err)
	}
	if err := repo.SaveTicks(ctx, ticks[:1]); err != nil {
		t.Fatalf("duplicate SaveTicks returned error: %v", err)
	}

	latestTick, err := repo.LatestTick(ctx, "XAUUSD")
	if err != nil {
		t.Fatalf("LatestTick returned error: %v", err)
	}
	if latestTick.Ask != 2325.72 {
		t.Fatalf("latestTick.Ask = %f, want 2325.72", latestTick.Ask)
	}
	assertRowCount(t, ctx, db, "mt5_ticks", 2)

	accountTime := time.Date(2026, 6, 18, 10, 2, 0, 0, time.UTC)
	err = repo.SaveAccountSnapshot(ctx, MT5AccountSnapshot{
		AccountLogin: "12345678",
		Currency:     "USD",
		Balance:      10000,
		Equity:       10080,
		Margin:       400,
		FreeMargin:   9680,
		MarginLevel:  2520,
		Time:         accountTime,
	})
	if err != nil {
		t.Fatalf("SaveAccountSnapshot returned error: %v", err)
	}

	latestAccount, err := repo.LatestAccountSnapshot(ctx, "12345678")
	if err != nil {
		t.Fatalf("LatestAccountSnapshot returned error: %v", err)
	}
	if latestAccount.Equity != 10080 {
		t.Fatalf("latestAccount.Equity = %f, want 10080", latestAccount.Equity)
	}

	snapshotTime := time.Date(2026, 6, 18, 10, 3, 0, 0, time.UTC)
	positions := []MT5PositionSnapshot{
		{
			AccountLogin: "12345678",
			Ticket:       "987654321",
			Symbol:       "XAUUSD",
			Side:         "buy",
			Volume:       0.1,
			OpenPrice:    2320.1,
			StopLoss:     2310,
			TakeProfit:   2340,
			Profit:       55,
			OpenedAt:     time.Date(2026, 6, 18, 9, 0, 0, 0, time.UTC),
			SnapshotTime: snapshotTime,
		},
	}
	if err := repo.SavePositionSnapshots(ctx, positions); err != nil {
		t.Fatalf("SavePositionSnapshots returned error: %v", err)
	}

	latestPositions, err := repo.LatestPositionSnapshots(ctx, "12345678")
	if err != nil {
		t.Fatalf("LatestPositionSnapshots returned error: %v", err)
	}
	if len(latestPositions) != 1 {
		t.Fatalf("len(latestPositions) = %d, want 1", len(latestPositions))
	}
	if latestPositions[0].Ticket != "987654321" {
		t.Fatalf("latestPositions[0].Ticket = %q, want 987654321", latestPositions[0].Ticket)
	}
}

func assertRowCount(t *testing.T, ctx context.Context, db *sql.DB, table string, want int) {
	t.Helper()

	var got int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM "+table).Scan(&got); err != nil {
		t.Fatalf("count rows in %s: %v", table, err)
	}
	if got != want {
		t.Fatalf("row count in %s = %d, want %d", table, got, want)
	}
}
