# MT5 XAUUSD MVP Implementation Tasks

> **For agentic workers:** Implement this plan task-by-task using TDD. Keep this MVP read-only until paper signal review is stable.

## Goal

Integrate MT5 for XAUUSD only through a Python bridge and Go ingest APIs, then support paper signals without real order execution.

## Task 1: MT5 Database Schema

**Files:**

- Create: `services/api/migrations/000004_create_mt5_tables.sql`
- Create: `services/api/internal/adapters/postgres/mt5_repository.go`
- Create: `services/api/internal/adapters/postgres/mt5_repository_test.go`

- [x] Write failing repository integration test for heartbeat, tick, account snapshot, and position snapshot persistence.
- [x] Add migration for `mt5_bridge_heartbeats`.
- [x] Add migration for `mt5_ticks`.
- [x] Add migration for `mt5_account_snapshots`.
- [x] Add migration for `mt5_position_snapshots`.
- [x] Implement PostgreSQL repository.
- [x] Run `go test ./internal/adapters/postgres`.
- [x] Commit with message `feat(mt5): add read-only mt5 schema`.

## Task 2: MT5 Domain And Application Service

**Files:**

- Create: `services/api/internal/domain/mt5/`
- Create: `services/api/internal/application/mt5/`
- Test: `services/api/internal/application/mt5/*_test.go`

- [x] Write failing tests for ingesting heartbeat, ticks, candles, account snapshots, and positions.
- [x] Add domain models for `Heartbeat`, `Tick`, `AccountSnapshot`, and `PositionSnapshot`.
- [x] Add application service that validates `XAUUSD` only.
- [x] Add candle ingest path that upserts MT5 candles into existing `candles`.
- [x] Reject non-XAUUSD payloads during MVP.
- [x] Run `go test ./internal/application/mt5`.
- [x] Commit with message `feat(mt5): add xauusd ingest service`.

## Task 3: MT5 HTTP Ingest API

**Files:**

- Create: `services/api/internal/adapters/http/mt5_handler.go`
- Test: `services/api/internal/adapters/http/mt5_handler_test.go`
- Modify: `services/api/internal/adapters/http/router.go`
- Modify: `services/api/cmd/api/main.go`

- [x] Write failing HTTP contract tests for MT5 ingest endpoints.
- [x] Add `POST /api/mt5/heartbeat`.
- [x] Add `POST /api/mt5/ticks`.
- [x] Add `POST /api/mt5/candles`.
- [x] Add `POST /api/mt5/account-snapshot`.
- [x] Add `POST /api/mt5/positions`.
- [x] Add `GET /api/mt5/status`.
- [x] Add latest account and positions read endpoints.
- [x] Run `go test ./internal/adapters/http ./...`.
- [x] Commit with message `feat(mt5): add ingest api`.

## Task 4: Python MT5 Bridge Skeleton

**Files:**

- Create: `bridges/mt5-python/`
- Create: `bridges/mt5-python/pyproject.toml`
- Create: `bridges/mt5-python/README.md`
- Create: `bridges/mt5-python/src/mt5_bridge/`
- Test: `bridges/mt5-python/tests/`

- [x] Add Python project skeleton.
- [x] Add config for Go API URL, symbol `XAUUSD`, timeframes, and polling intervals.
- [x] Add payload builders with unit tests.
- [x] Add dry-run mode that prints payloads without connecting to MT5.
- [x] Add `MetaTrader5` adapter boundary but keep it replaceable in tests.
- [ ] Commit with message `feat(mt5): add python bridge skeleton`.

## Task 5: Bridge Status Dashboard

**Files:**

- Create: `apps/web/app/mt5/page.tsx`
- Create: `apps/web/features/mt5/`
- Modify: `apps/web/components/app-shell.tsx`
- Modify: `apps/web/lib/api-client.ts`

- [ ] Add API client methods for MT5 status, account, and positions.
- [ ] Add MT5 navigation item.
- [ ] Add bridge health panel.
- [ ] Add latest tick/spread panel.
- [ ] Add account snapshot panel.
- [ ] Add open positions table.
- [ ] Add smoke test for MT5 page route.
- [ ] Run `npm test`, `npm run typecheck`, and `npm run build`.
- [ ] Commit with message `feat(web): add mt5 bridge dashboard`.

## Task 6: Paper Signal Foundation

**Files:**

- Create: `services/api/internal/domain/signal/`
- Create: `services/api/internal/application/signals/`
- Create: `services/api/internal/adapters/http/signals_handler.go`
- Create: `services/api/migrations/000005_create_paper_signals.sql`

- [ ] Add paper signal schema.
- [ ] Add signal domain model with status lifecycle.
- [ ] Add `POST /api/paper-signals`.
- [ ] Add `GET /api/paper-signals`.
- [ ] Add `PATCH /api/paper-signals/{id}/status`.
- [ ] Ensure no execution command is created.
- [ ] Run `go test ./...`.
- [ ] Commit with message `feat(signals): add xauusd paper signals`.
