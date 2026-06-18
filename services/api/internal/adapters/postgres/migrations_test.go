package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestRunMigrationsCreatesInitialSchema(t *testing.T) {
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

	migrationsDir := filepath.Join("..", "..", "..", "migrations")
	if err := RunMigrations(ctx, db, migrationsDir); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	if err := RunMigrations(ctx, db, migrationsDir); err != nil {
		t.Fatalf("run migrations again: %v", err)
	}

	assertTableExists(t, ctx, db, "symbols")
	assertTableExists(t, ctx, db, "candles")
	assertTimescaleExtensionExists(t, ctx, db)
	assertCandlesIsHypertable(t, ctx, db)
}

func resetMigrationState(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()

	statements := []string{
		"DROP TABLE IF EXISTS candles",
		"DROP TABLE IF EXISTS symbols",
		"DROP TABLE IF EXISTS schema_migrations",
	}
	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("reset migration state with %q: %v", statement, err)
		}
	}
}

func assertTableExists(t *testing.T, ctx context.Context, db *sql.DB, table string) {
	t.Helper()

	var exists bool
	err := db.QueryRowContext(ctx, "SELECT to_regclass($1) IS NOT NULL", "public."+table).Scan(&exists)
	if err != nil {
		t.Fatalf("check table %s: %v", table, err)
	}
	if !exists {
		t.Fatalf("expected table %s to exist", table)
	}
}

func assertTimescaleExtensionExists(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()

	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'timescaledb')").Scan(&exists)
	if err != nil {
		t.Fatalf("check timescaledb extension: %v", err)
	}
	if !exists {
		t.Fatal("expected timescaledb extension to exist")
	}
}

func assertCandlesIsHypertable(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()

	var exists bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM timescaledb_information.hypertables
			WHERE hypertable_schema = 'public'
			  AND hypertable_name = 'candles'
		)
	`).Scan(&exists)
	if err != nil {
		t.Fatalf("check candles hypertable: %v", err)
	}
	if !exists {
		t.Fatal("expected candles to be a TimescaleDB hypertable")
	}
}
