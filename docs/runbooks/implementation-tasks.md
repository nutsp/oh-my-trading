# AI Trading Agent Dashboard Tasks

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

## Milestone Plan

### Phase 0: Project Setup

- Repo structure.
- Docker Compose.
- Go API health endpoint.
- Next.js shell.
- Postgres migrations.
- CI test workflow.

### Phase 1: Market Data

- Symbol model.
- Candle ingestion.
- Candle API.
- Chart page.
- Redis latest-candle cache.

### Phase 2: Indicators And Structure

- EMA, RSI, MACD, ATR.
- Indicator persistence.
- Swing detection.
- Basic BOS, CHoCH, FVG detection.

### Phase 3: AI Signals

- Agent workflow engine.
- AI synthesis node.
- Signal persistence.
- Signal UI.
- Agent run logs.

### Phase 4: Risk And Journal

- Risk settings.
- Position sizing.
- Trade validation.
- Journal CRUD.
- Trade detail page.

### Phase 5: Backtesting

- Strategy config.
- Historical replay.
- Backtest result persistence.
- Equity curve UI.

### Phase 6: Monitoring And Notifications

- Telegram/Discord notifications.
- Agent dashboard.
- Error visibility.
- Metrics.

### Phase 7: Advanced

- Paper trading.
- Broker integration.
- Auto-trading.
- Portfolio optimization.
- AI self-review.

## Development Tasks

## Current MVP Pivot

The project priority has changed to **XAUUSD only, MT5 read-only integration, and paper signals first**.

Before continuing generic roadmap tasks such as market structure detection, agent workflow, broad backtesting, or multi-market support, follow:

- [MT5 XAUUSD MVP](../product/mt5-xauusd-mvp.md)
- [MT5 Integration Architecture](../architecture/mt5-integration.md)
- [MT5 XAUUSD MVP Tasks](./mt5-xauusd-mvp-tasks.md)

Generic tasks below remain useful but are lower priority until the MT5 MVP is working.

### Task 1: Repository Scaffold

**Files:**

- Create: `apps/web/`
- Create: `services/api/`
- Create: `deployments/`
- Create: `docs/architecture/`
- Create: `docs/product/`
- Create: `docs/runbooks/`

- [x] Create base directory structure.
- [x] Add root README with project purpose, local setup, and architecture summary.
- [x] Commit with message `chore: add project scaffold`.

### Task 2: Local Infrastructure

**Files:**

- Create: `deployments/docker-compose.yml`
- Create: `deployments/.env.example`

- [x] Add PostgreSQL/TimescaleDB service.
- [x] Add Redis service.
- [x] Add RabbitMQ service.
- [x] Add health checks for all services.
- [x] Run `docker compose -f deployments/docker-compose.yml up -d`.
  - Verification note: `docker compose config` passed; stack is healthy.
  - Retry note: default PostgreSQL host port changed to `15432` because local port `5432` was already allocated.
- [x] Commit with message `chore: add local infrastructure`.

### Task 3: Go API Foundation

**Files:**

- Create: `services/api/cmd/api/main.go`
- Create: `services/api/internal/platform/config/`
- Create: `services/api/internal/platform/logger/`
- Create: `services/api/internal/adapters/http/`

- [x] Initialize Go module.
- [x] Add config loader.
- [x] Add structured logger.
- [x] Add `GET /api/health`.
- [x] Add graceful shutdown.
- [x] Add API smoke test.
- [x] Commit with message `feat(api): add service foundation`.

### Task 4: Database Migrations

**Files:**

- Create: `services/api/migrations/`
- Create: `services/api/internal/adapters/postgres/`

- [x] Add migration tool.
- [x] Add symbols migration.
- [x] Add candles migration.
- [x] Enable TimescaleDB extension.
- [x] Convert candles to hypertable.
- [x] Add migration test.
- [x] Commit with message `feat(db): add symbols and candles schema`.

### Task 5: Symbol Management

**Files:**

- Create: `services/api/internal/domain/market/`
- Create: `services/api/internal/application/marketdata/`
- Modify: `services/api/internal/adapters/http/`
- Modify: `services/api/internal/adapters/postgres/`

- [x] Add `Symbol` domain entity.
- [x] Add symbol repository port.
- [x] Add PostgreSQL repository adapter.
- [x] Add `GET /api/symbols`.
- [x] Add `POST /api/symbols`.
- [x] Add unit and integration tests.
- [x] Commit with message `feat(market): add symbol management`.

### Task 6: Candle Storage And Query API

**Files:**

- Modify: `services/api/internal/domain/market/`
- Modify: `services/api/internal/application/marketdata/`
- Modify: `services/api/internal/adapters/http/`
- Modify: `services/api/internal/adapters/postgres/`

- [x] Add `Candle` domain entity.
- [x] Add candle repository port.
- [x] Add bulk upsert.
- [x] Add time-range query.
- [x] Add `GET /api/candles`.
- [x] Add integration tests.
- [x] Commit with message `feat(market): add candle storage and query api`.

### Task 7: Market Data Sync Worker

**Files:**

- Create: `services/api/cmd/worker/main.go`
- Create: `services/api/internal/adapters/marketdata/`
- Create: `services/api/internal/adapters/rabbitmq/`

- [x] Define market data provider interface.
- [x] Implement first provider adapter.
- [x] Add candle sync job.
- [x] Publish sync requests through RabbitMQ.
- [x] Store fetched candles in TimescaleDB.
- [x] Add fake-provider worker test.
- [x] Commit with message `feat(worker): add market data sync`.

### Task 8: Next.js Dashboard Shell

**Files:**

- Create: `apps/web/app/`
- Create: `apps/web/components/`
- Create: `apps/web/features/`
- Create: `apps/web/lib/`

- [x] Initialize Next.js app.
- [x] Add dashboard layout.
- [x] Add navigation for Dashboard, Markets, Signals, Risk, Journal, Backtests, Agent, Settings.
- [x] Add API client wrapper.
- [x] Add basic visual system.
- [x] Add smoke test.
- [x] Commit with message `feat(web): add dashboard shell`.

### Task 9: Market Chart Page

**Files:**

- Create: `apps/web/app/markets/[symbol]/page.tsx`
- Create: `apps/web/features/markets/`

- [x] Install TradingView Lightweight Charts.
- [x] Add symbol route.
- [x] Fetch candles from API.
- [x] Render candlestick chart.
- [x] Add timeframe selector.
- [x] Add loading, empty, and error states.
- [x] Add Playwright smoke test.
- [ ] Commit with message `feat(web): add market chart page`.

### Task 10: Indicators

**Files:**

- Create: `services/api/internal/domain/market/indicators/`
- Modify: `services/api/internal/application/marketdata/`
- Modify: `services/api/internal/adapters/http/`

- [x] Implement EMA.
- [x] Implement RSI.
- [x] Implement MACD.
- [x] Implement ATR.
- [x] Add golden fixture tests.
- [x] Add `GET /api/indicators`.
- [x] Add chart overlays in frontend.
- [ ] Commit with message `feat(analysis): add core indicators`.

### Task 11: Market Structure Detection

**Files:**

- Create: `services/api/internal/domain/market/structure/`
- Modify: `services/api/migrations/`
- Modify: `services/api/internal/adapters/postgres/`

- [ ] Add swing high/low detection.
- [ ] Add basic BOS detection.
- [ ] Add basic CHoCH detection.
- [ ] Add basic FVG detection.
- [ ] Persist detected structures.
- [ ] Render structures on chart.
- [ ] Commit with message `feat(analysis): add market structure detection`.

### Task 12: Agent Workflow Engine

**Files:**

- Create: `services/api/internal/domain/agent/`
- Create: `services/api/internal/application/analysis/`
- Create: `services/api/internal/adapters/ai/`
- Modify: `services/api/migrations/`

- [ ] Add workflow state model.
- [ ] Add node interface.
- [ ] Add deterministic nodes for market context, indicators, structure, and top-down bias.
- [ ] Add AI synthesis provider interface.
- [ ] Persist agent runs.
- [ ] Add workflow tests with mocked AI.
- [ ] Commit with message `feat(agent): add analysis workflow engine`.

### Task 13: Signal Generation

**Files:**

- Create: `services/api/internal/domain/signal/`
- Modify: `services/api/internal/application/analysis/`
- Modify: `services/api/internal/adapters/http/`
- Modify: `apps/web/app/signals/`

- [ ] Add signal domain entity.
- [ ] Add signal repository.
- [ ] Add signal scoring.
- [ ] Add `POST /api/analysis/run`.
- [ ] Add `GET /api/signals`.
- [ ] Add signal list page.
- [ ] Add signal detail page.
- [ ] Commit with message `feat(signals): add ai signal workflow`.

### Task 14: Risk Engine

**Files:**

- Create: `services/api/internal/domain/risk/`
- Create: `services/api/internal/application/risk/`
- Modify: `services/api/internal/adapters/http/`
- Create: `apps/web/app/risk/page.tsx`

- [ ] Add risk settings model.
- [ ] Add position sizing calculation.
- [ ] Add reward-to-risk validation.
- [ ] Add daily loss validation.
- [ ] Add ATR stop validation.
- [ ] Add `POST /api/risk/calculate-position`.
- [ ] Add `POST /api/risk/validate-trade`.
- [ ] Add risk calculator UI.
- [ ] Commit with message `feat(risk): add risk engine`.

### Task 15: Trade Journal

**Files:**

- Create: `services/api/internal/domain/journal/`
- Create: `services/api/internal/application/journal/`
- Modify: `services/api/migrations/`
- Modify: `services/api/internal/adapters/http/`
- Create: `apps/web/app/journal/`

- [ ] Add trades table migration.
- [ ] Add trade notes table migration.
- [ ] Add trade lifecycle rules.
- [ ] Add journal CRUD APIs.
- [ ] Add journal list page.
- [ ] Add trade detail page.
- [ ] Add post-trade review fields.
- [ ] Commit with message `feat(journal): add trade journal`.

### Task 16: Backtesting

**Files:**

- Create: `services/api/internal/domain/backtest/`
- Create: `services/api/internal/application/backtest/`
- Modify: `services/api/migrations/`
- Modify: `services/api/internal/adapters/http/`
- Create: `apps/web/app/backtests/`

- [ ] Add backtest config model.
- [ ] Add sequential candle replay engine.
- [ ] Add no-lookahead tests.
- [ ] Add risk model simulation.
- [ ] Add metrics calculation.
- [ ] Persist backtest runs.
- [ ] Add backtest list and detail UI.
- [ ] Commit with message `feat(backtest): add strategy replay engine`.

### Task 17: Notifications

**Files:**

- Create: `services/api/internal/adapters/notification/`
- Modify: `services/api/internal/application/analysis/`
- Modify: `services/api/internal/application/backtest/`

- [ ] Add notification provider interface.
- [ ] Add Telegram adapter.
- [ ] Add Discord webhook adapter.
- [ ] Notify on signal creation.
- [ ] Notify on backtest completion.
- [ ] Notify on agent run failure.
- [ ] Commit with message `feat(notification): add signal and agent alerts`.

### Task 18: Agent Monitoring

**Files:**

- Modify: `services/api/internal/adapters/http/`
- Create: `apps/web/app/agent/page.tsx`
- Create: `apps/web/features/agent/`

- [ ] Add `GET /api/agent-runs`.
- [ ] Add `GET /api/agent-runs/{id}`.
- [ ] Add agent runs table.
- [ ] Add run detail view with input, output, errors, and timing.
- [ ] Add failed-run highlighting.
- [ ] Commit with message `feat(agent): add monitoring dashboard`.

### Task 19: CI And Test Coverage

**Files:**

- Create: `.github/workflows/ci.yml`
- Modify: backend and frontend test config files.

- [ ] Run Go unit tests in CI.
- [ ] Run Go integration tests with Postgres service.
- [ ] Run frontend unit tests.
- [ ] Run Playwright smoke tests.
- [ ] Add lint checks.
- [ ] Commit with message `ci: add test workflow`.

### Task 20: Deployment Documentation

**Files:**

- Create: `docs/runbooks/local-development.md`
- Create: `docs/runbooks/deployment.md`
- Create: `docs/runbooks/backup-restore.md`
- Create: `docs/architecture/system-overview.md`

- [ ] Document local setup.
- [ ] Document environment variables.
- [ ] Document deployment steps.
- [ ] Document backup and restore.
- [ ] Document architecture boundaries.
- [ ] Commit with message `docs: add deployment and architecture runbooks`.
