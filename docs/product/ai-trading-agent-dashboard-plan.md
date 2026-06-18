# AI Trading Agent Dashboard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a personal AI trading command center for market analysis, AI signal generation, risk control, trade journaling, backtesting, and agent monitoring across XAUUSD, BTC, Forex, and Crypto.

**Architecture:** Use a modular monolith with clean/hexagonal architecture. Domain logic stays isolated from HTTP, database, market data providers, AI providers, queues, cache, and notification adapters. Split into separate services only after the MVP proves useful and operational complexity is justified.

**Tech Stack:** Go, Next.js, PostgreSQL with TimescaleDB, Redis, RabbitMQ, TradingView Lightweight Charts, Docker Compose, custom LangGraph-style agent workflow.

---

## 1. Product Vision

The dashboard should act as the operating system for a personal trading AI agent.

It should answer:

- What is the market doing across multiple timeframes?
- What setups does the AI agent see?
- Is the trade valid under the risk rules?
- What is the risk before entry?
- What happened after execution?
- Is the agent behaving correctly?
- Which strategies perform best historically?

The product should prioritize decision support before automation. Auto-trading should come later after analysis, journaling, risk checks, and backtesting are trustworthy.

## 2. Core Modules

### Market Data

- Symbols: `XAUUSD`, `BTCUSD`, major Forex pairs, selected crypto pairs.
- Candles: `1m`, `5m`, `15m`, `1h`, `4h`, `1d`.
- Indicators: EMA, RSI, MACD, ATR.
- Market structure: swing highs/lows, BOS, CHoCH, order blocks, fair value gaps.

### AI Analysis

- Multi-timeframe top-down analysis.
- Strategy confluence scoring.
- Signal generation.
- Agent reasoning trace.
- Confidence score and invalidation conditions.

### Risk Management

- Account risk limits.
- Position sizing.
- Stop-loss validation.
- ATR-based risk checks.
- Max daily loss.
- Max open exposure.
- Reward-to-risk validation.

### Trade Journal

- Planned trades.
- Executed trades.
- Screenshots and notes.
- AI rationale.
- Post-trade review.
- Mistake tags.
- PnL and performance analytics.

### Backtesting

- Strategy definitions.
- Historical candle replay.
- Signal simulation.
- Risk model simulation.
- Metrics: win rate, expectancy, drawdown, profit factor.

### Agent Monitoring

- Agent runs.
- Prompt/input snapshots.
- Tool calls.
- Errors.
- Latency.
- Signal outcomes.

## 3. MVP Scope

The MVP should avoid live execution. Build a reliable analysis and journaling system first.

### MVP Includes

- Single-user personal dashboard with local admin auth.
- Symbol and timeframe watchlist.
- Candle ingestion from one provider.
- TimescaleDB candle storage.
- TradingView Lightweight Charts.
- EMA, RSI, MACD, ATR calculations.
- Manual AI analysis trigger.
- AI-generated signal object.
- Risk calculator.
- Trade journal CRUD.
- Agent run logs.
- Basic backtest engine for one strategy.
- Telegram, Discord, or email notification for generated signals.

### MVP Excludes

- Auto-trading.
- Broker execution.
- Multi-user roles.
- Portfolio optimization.
- AI self-improvement loop.
- Full ICT/SMC automation beyond basic structure, FVG, and order block detection.

## 4. Full Feature Roadmap

### Phase 1: Foundation

- Backend skeleton.
- Database migrations.
- Symbol and timeframe model.
- Candle ingestion.
- Chart UI.

### Phase 2: Analysis

- Indicators.
- Market structure detection.
- AI analysis workflow.
- Signal dashboard.

### Phase 3: Risk And Journal

- Risk rules.
- Position sizing.
- Trade planning.
- Trade journal.
- Post-trade review.

### Phase 4: Backtesting

- Historical replay.
- Strategy configuration.
- Backtest reports.
- Equity curve.

### Phase 5: Agent Operations

- Agent run inspector.
- Prompt and version tracking.
- Tool-call logs.
- Signal performance tracking.

### Phase 6: Advanced Trading

- Broker integration.
- Paper trading.
- Semi-automated execution.
- Full auto-trading with kill switch.

### Phase 7: Intelligence Layer

- AI self-review.
- Strategy comparison.
- Portfolio optimization.
- Agent memory.
- Market regime detection.

## 5. System Architecture

Use a modular monolith first.

```text
Next.js Dashboard
   |
REST/WebSocket API
   |
Go Backend
   |
Domain Services
   |-- Market Data
   |-- Indicators
   |-- Signal Engine
   |-- Risk Engine
   |-- Journal
   |-- Backtesting
   |-- Agent Runtime
   |
Adapters
   |-- PostgreSQL/TimescaleDB
   |-- Redis
   |-- RabbitMQ
   |-- Market Data Provider
   |-- AI Provider
   |-- Notification Provider
```

### Technical Decisions

- **Go backend:** Strong concurrency, reliable services, simple deployment.
- **Next.js frontend:** Dashboard routing, React ecosystem, and chart integration.
- **TimescaleDB:** Best fit for candle and indicator time-series data.
- **Redis:** Cache latest candles, indicators, and agent states.
- **RabbitMQ:** Simpler than Kafka for a personal workflow queue.
- **Kafka later:** Useful only if event volume or stream replay requirements grow significantly.
- **TradingView Lightweight Charts:** Fast and clean trading UX.
- **Custom agent workflow:** Easier to own and debug than adopting LangGraph concepts wholesale.

## 6. Database Design

Use PostgreSQL with the TimescaleDB extension.

### Core Tables

```sql
CREATE TABLE symbols (
  id uuid PRIMARY KEY,
  code text UNIQUE NOT NULL,
  market text NOT NULL,
  base_asset text,
  quote_asset text,
  enabled boolean NOT NULL DEFAULT true
);

CREATE TABLE candles (
  symbol_id uuid NOT NULL,
  timeframe text NOT NULL,
  ts timestamptz NOT NULL,
  open numeric NOT NULL,
  high numeric NOT NULL,
  low numeric NOT NULL,
  close numeric NOT NULL,
  volume numeric,
  PRIMARY KEY (symbol_id, timeframe, ts)
);
```

Convert `candles` to a Timescale hypertable.

```sql
SELECT create_hypertable('candles', 'ts');
```

```sql
CREATE TABLE indicator_values (
  id uuid PRIMARY KEY,
  symbol_id uuid NOT NULL,
  timeframe text NOT NULL,
  ts timestamptz NOT NULL,
  indicator text NOT NULL,
  params jsonb NOT NULL,
  value jsonb NOT NULL
);

CREATE TABLE market_structures (
  id uuid PRIMARY KEY,
  symbol_id uuid NOT NULL,
  timeframe text NOT NULL,
  ts timestamptz NOT NULL,
  type text NOT NULL,
  direction text,
  price_low numeric,
  price_high numeric,
  metadata jsonb NOT NULL DEFAULT '{}'
);

CREATE TABLE signals (
  id uuid PRIMARY KEY,
  symbol_id uuid NOT NULL,
  timeframe text NOT NULL,
  direction text NOT NULL,
  entry numeric,
  stop_loss numeric,
  take_profit numeric,
  confidence numeric NOT NULL,
  status text NOT NULL,
  reason text NOT NULL,
  invalidation text,
  created_at timestamptz NOT NULL
);

CREATE TABLE trades (
  id uuid PRIMARY KEY,
  signal_id uuid,
  symbol_id uuid NOT NULL,
  direction text NOT NULL,
  status text NOT NULL,
  planned_entry numeric,
  executed_entry numeric,
  stop_loss numeric,
  take_profit numeric,
  position_size numeric,
  risk_amount numeric,
  pnl numeric,
  opened_at timestamptz,
  closed_at timestamptz
);

CREATE TABLE trade_notes (
  id uuid PRIMARY KEY,
  trade_id uuid NOT NULL,
  note_type text NOT NULL,
  content text NOT NULL,
  tags text[],
  created_at timestamptz NOT NULL
);

CREATE TABLE agent_runs (
  id uuid PRIMARY KEY,
  workflow text NOT NULL,
  symbol_id uuid,
  timeframe text,
  status text NOT NULL,
  input jsonb NOT NULL,
  output jsonb,
  error text,
  started_at timestamptz NOT NULL,
  finished_at timestamptz
);

CREATE TABLE backtests (
  id uuid PRIMARY KEY,
  name text NOT NULL,
  strategy text NOT NULL,
  config jsonb NOT NULL,
  symbol_id uuid NOT NULL,
  timeframe text NOT NULL,
  from_ts timestamptz NOT NULL,
  to_ts timestamptz NOT NULL,
  metrics jsonb,
  created_at timestamptz NOT NULL
);
```

## 7. API Design

Use REST for CRUD/query operations and WebSocket for live updates.

### Market Data

- `GET /api/symbols`
- `POST /api/symbols`
- `GET /api/candles?symbol=XAUUSD&timeframe=1h&from=&to=`
- `POST /api/market-data/sync`
- `GET /api/indicators?symbol=BTCUSD&timeframe=4h`

### Analysis And Signals

- `POST /api/analysis/run`
- `GET /api/analysis/runs`
- `GET /api/signals`
- `GET /api/signals/{id}`
- `PATCH /api/signals/{id}/status`

Example request:

```json
{
  "symbol": "XAUUSD",
  "timeframes": ["1d", "4h", "1h", "15m"],
  "strategyProfile": "top_down_ema_rsi_macd_atr_smc"
}
```

Example signal response:

```json
{
  "symbol": "XAUUSD",
  "direction": "long",
  "confidence": 0.72,
  "entry": 2325.5,
  "stopLoss": 2314.2,
  "takeProfit": 2350.0,
  "reason": "Daily bullish structure, 4H pullback into EMA zone, RSI recovery, ATR supports stop distance.",
  "invalidation": "Close below 2314.2 or bearish CHoCH on 1H."
}
```

### Risk

- `POST /api/risk/calculate-position`
- `POST /api/risk/validate-trade`
- `GET /api/risk/settings`
- `PUT /api/risk/settings`

### Journal

- `GET /api/trades`
- `POST /api/trades`
- `GET /api/trades/{id}`
- `PATCH /api/trades/{id}`
- `POST /api/trades/{id}/notes`

### Backtesting

- `POST /api/backtests`
- `GET /api/backtests`
- `GET /api/backtests/{id}`
- `GET /api/backtests/{id}/trades`

### Monitoring

- `GET /api/agent-runs`
- `GET /api/agent-runs/{id}`
- `GET /api/health`
- `GET /api/metrics`

### WebSocket

- `/ws/market`
- `/ws/signals`
- `/ws/agent-runs`

## 8. Agent Workflow Design

Use an explicit state machine.

```text
Start
  -> Load Market Context
  -> Compute Indicators
  -> Detect Structure
  -> Build Multi-Timeframe Bias
  -> Generate Candidate Setup
  -> Validate Risk
  -> Score Confidence
  -> Produce Signal
  -> Store Agent Run
  -> Notify
End
```

### Agent State

```json
{
  "symbol": "BTCUSD",
  "timeframes": ["1d", "4h", "1h", "15m"],
  "marketContext": {},
  "indicators": {},
  "structure": {},
  "bias": "bullish",
  "candidateSetup": {},
  "riskValidation": {},
  "signal": {},
  "reasoning": []
}
```

### Recommended Agent Nodes

- `MarketContextNode`
- `IndicatorNode`
- `StructureNode`
- `TopDownBiasNode`
- `SetupDetectionNode`
- `RiskValidationNode`
- `SignalScoringNode`
- `NarrativeNode`
- `PersistenceNode`
- `NotificationNode`

Keep every node deterministic where possible. Use AI mostly for synthesis, explanation, scenario analysis, and self-review, not raw indicator calculation.

## 9. Frontend Page Structure

Use Next.js App Router.

```text
/dashboard
  Overview, watchlist, latest signals, agent status

/markets/[symbol]
  TradingView chart, indicators, structures, timeframe selector

/signals
  Signal list, filters, status, confidence, risk validity

/signals/[id]
  Full reasoning, chart context, risk result, journal conversion

/risk
  Risk settings, calculator, exposure overview

/journal
  Trade list, calendar, PnL summary, mistake tags

/journal/[tradeId]
  Trade detail, screenshots, notes, AI review

/backtests
  Backtest runs, create backtest, compare results

/backtests/[id]
  Equity curve, trade list, metrics, chart replay

/agent
  Agent runs, workflow status, errors, latency

/settings
  Symbols, providers, notifications, risk config
```

## 10. Backend Service Structure

Recommended Go layout:

```text
cmd/
  api/
    main.go
  worker/
    main.go

internal/
  domain/
    market/
    signal/
    risk/
    journal/
    backtest/
    agent/

  application/
    marketdata/
    analysis/
    risk/
    journal/
    backtest/
    monitoring/

  adapters/
    http/
    postgres/
    timescale/
    redis/
    rabbitmq/
    marketdata/
    ai/
    notification/

  platform/
    config/
    logger/
    metrics/
    auth/
    clock/

migrations/
web/
  app/
  components/
  features/
  lib/
  styles/
```

### Boundary Rule

- `domain` knows no database, HTTP, Redis, AI provider, or queue.
- `application` orchestrates use cases.
- `adapters` implement ports/interfaces.
- `cmd` wires dependencies.

## 11. Risk Management Rules

### MVP Rules

- Max risk per trade: configurable, default `1%`.
- Max daily loss: configurable, default `3%`.
- Max open trades: default `3`.
- Minimum reward-to-risk: default `1.5`.
- Stop loss is required.
- Position size must be calculated from account equity, stop distance, and pip/tick value.
- Trade is blocked if ATR stop is too tight or too wide.
- Trade is blocked if correlated exposure is too high.

Example validation result:

```json
{
  "valid": false,
  "reasons": [
    "Reward-to-risk 1.1 is below minimum 1.5",
    "Daily loss limit would be exceeded"
  ],
  "suggestedPositionSize": 0.04,
  "riskAmount": 100
}
```

### Advanced Rules

- News-event lockout.
- Session-based rules.
- Volatility regime adjustment.
- Consecutive loss cooldown.
- AI confidence threshold by strategy.
- Separate rules per asset class.

## 12. Backtesting Design

### MVP Backtester

- Input: symbol, timeframe, date range, strategy config.
- Load historical candles.
- Walk candles sequentially.
- Calculate indicators without lookahead.
- Generate simulated signals.
- Apply risk model.
- Track trades and equity.

### Metrics

- Total return.
- Win rate.
- Profit factor.
- Max drawdown.
- Average R multiple.
- Expectancy.
- Sharpe-like ratio later.
- Long/short breakdown.
- Session/timeframe performance.

### Important Rules

- No future candle access.
- Store strategy version.
- Store config snapshot.
- Store fees, spread, and slippage assumptions.
- Make every backtest reproducible.

## 13. Trade Journal Design

A trade can start from either a manual idea or an AI signal.

### Trade Lifecycle

```text
planned -> active -> closed -> reviewed
```

### Journal Fields

- Symbol.
- Direction.
- Strategy tag.
- Setup type.
- Entry, stop, take profit.
- Planned risk.
- Actual position size.
- Entry/exit screenshots.
- AI thesis.
- Human notes.
- Mistake tags.
- Emotional state.
- Outcome in R.
- Post-trade review.

### AI Review

After close, the agent reviews:

- Was the setup valid?
- Was entry too early or too late?
- Was stop placement logical?
- Did price respect the original thesis?
- What should be improved?

## 14. Notification Design

Start simple.

### Channels

- Telegram bot.
- Discord webhook.
- Email later.
- In-app notification center.

### Events

- New signal generated.
- Risk validation failed.
- Backtest completed.
- Agent run failed.
- Daily loss limit hit.
- Trade review pending.

### Notification Payload

```json
{
  "event": "signal.created",
  "symbol": "XAUUSD",
  "direction": "long",
  "confidence": 0.72,
  "riskValid": true,
  "message": "XAUUSD long setup detected on 1H with 4H bullish bias."
}
```

## 15. Security Considerations

For personal use, keep security simple but disciplined.

- Local admin login or single-user auth.
- Store secrets in environment variables.
- Encrypt broker/API keys before database storage.
- Never expose auto-trading endpoints publicly.
- Use HTTPS in production.
- Rate-limit API endpoints.
- Validate all order/risk inputs server-side.
- Add audit logs for signal generation, trade changes, and future execution.
- Separate read-only market data keys from trading keys.
- Add a global kill switch before broker integration.

## 16. Deployment Plan

### MVP Deployment

Use Docker Compose on a VPS or local machine.

Services:

- `api`
- `worker`
- `web`
- `postgres-timescale`
- `redis`
- `rabbitmq`
- `prometheus` optional
- `grafana` optional

### Production-ish Setup

- Nginx or Caddy reverse proxy.
- HTTPS certificates.
- Nightly database backup.
- Structured logs.
- Health checks.
- Separate `.env.production`.
- Basic uptime monitor.

### Later

- Kubernetes only if the system becomes multi-service or high-availability.
- Object storage for screenshots and reports.
- Separate read replica for analytics.

## 17. Milestone Plan By Phase

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

## 18. Recommended Folder Structure

```text
oh-my-trading/
  apps/
    web/
      app/
      components/
      features/
      lib/
      tests/

  services/
    api/
      cmd/api/
      cmd/worker/
      internal/domain/
      internal/application/
      internal/adapters/
      internal/platform/
      migrations/
      tests/

  packages/
    shared/
      schemas/
      openapi/

  deployments/
    docker-compose.yml
    nginx/
    grafana/

  docs/
    architecture/
    product/
    runbooks/
    superpowers/
      plans/
```

## 19. Development Tasks For Codex

### Task 1: Repository Scaffold

**Files:**

- Create: `apps/web/`
- Create: `services/api/`
- Create: `deployments/`
- Create: `docs/architecture/`
- Create: `docs/product/`
- Create: `docs/runbooks/`

- [ ] Create base directory structure.
- [ ] Add root README with project purpose, local setup, and architecture summary.
- [ ] Commit with message `chore: add project scaffold`.

### Task 2: Local Infrastructure

**Files:**

- Create: `deployments/docker-compose.yml`
- Create: `deployments/.env.example`

- [ ] Add PostgreSQL/TimescaleDB service.
- [ ] Add Redis service.
- [ ] Add RabbitMQ service.
- [ ] Add health checks for all services.
- [ ] Run `docker compose -f deployments/docker-compose.yml up -d`.
- [ ] Commit with message `chore: add local infrastructure`.

### Task 3: Go API Foundation

**Files:**

- Create: `services/api/cmd/api/main.go`
- Create: `services/api/internal/platform/config/`
- Create: `services/api/internal/platform/logger/`
- Create: `services/api/internal/adapters/http/`

- [ ] Initialize Go module.
- [ ] Add config loader.
- [ ] Add structured logger.
- [ ] Add `GET /api/health`.
- [ ] Add graceful shutdown.
- [ ] Add API smoke test.
- [ ] Commit with message `feat(api): add service foundation`.

### Task 4: Database Migrations

**Files:**

- Create: `services/api/migrations/`
- Create: `services/api/internal/adapters/postgres/`

- [ ] Add migration tool.
- [ ] Add symbols migration.
- [ ] Add candles migration.
- [ ] Enable TimescaleDB extension.
- [ ] Convert candles to hypertable.
- [ ] Add migration test.
- [ ] Commit with message `feat(db): add symbols and candles schema`.

### Task 5: Symbol Management

**Files:**

- Create: `services/api/internal/domain/market/`
- Create: `services/api/internal/application/marketdata/`
- Modify: `services/api/internal/adapters/http/`
- Modify: `services/api/internal/adapters/postgres/`

- [ ] Add `Symbol` domain entity.
- [ ] Add symbol repository port.
- [ ] Add PostgreSQL repository adapter.
- [ ] Add `GET /api/symbols`.
- [ ] Add `POST /api/symbols`.
- [ ] Add unit and integration tests.
- [ ] Commit with message `feat(market): add symbol management`.

### Task 6: Candle Storage And Query API

**Files:**

- Modify: `services/api/internal/domain/market/`
- Modify: `services/api/internal/application/marketdata/`
- Modify: `services/api/internal/adapters/http/`
- Modify: `services/api/internal/adapters/postgres/`

- [ ] Add `Candle` domain entity.
- [ ] Add candle repository port.
- [ ] Add bulk upsert.
- [ ] Add time-range query.
- [ ] Add `GET /api/candles`.
- [ ] Add integration tests.
- [ ] Commit with message `feat(market): add candle storage and query api`.

### Task 7: Market Data Sync Worker

**Files:**

- Create: `services/api/cmd/worker/main.go`
- Create: `services/api/internal/adapters/marketdata/`
- Create: `services/api/internal/adapters/rabbitmq/`

- [ ] Define market data provider interface.
- [ ] Implement first provider adapter.
- [ ] Add candle sync job.
- [ ] Publish sync requests through RabbitMQ.
- [ ] Store fetched candles in TimescaleDB.
- [ ] Add fake-provider worker test.
- [ ] Commit with message `feat(worker): add market data sync`.

### Task 8: Next.js Dashboard Shell

**Files:**

- Create: `apps/web/app/`
- Create: `apps/web/components/`
- Create: `apps/web/features/`
- Create: `apps/web/lib/`

- [ ] Initialize Next.js app.
- [ ] Add dashboard layout.
- [ ] Add navigation for Dashboard, Markets, Signals, Risk, Journal, Backtests, Agent, Settings.
- [ ] Add API client wrapper.
- [ ] Add basic visual system.
- [ ] Add smoke test.
- [ ] Commit with message `feat(web): add dashboard shell`.

### Task 9: Market Chart Page

**Files:**

- Create: `apps/web/app/markets/[symbol]/page.tsx`
- Create: `apps/web/features/markets/`

- [ ] Install TradingView Lightweight Charts.
- [ ] Add symbol route.
- [ ] Fetch candles from API.
- [ ] Render candlestick chart.
- [ ] Add timeframe selector.
- [ ] Add loading, empty, and error states.
- [ ] Add Playwright smoke test.
- [ ] Commit with message `feat(web): add market chart page`.

### Task 10: Indicators

**Files:**

- Create: `services/api/internal/domain/market/indicators/`
- Modify: `services/api/internal/application/marketdata/`
- Modify: `services/api/internal/adapters/http/`

- [ ] Implement EMA.
- [ ] Implement RSI.
- [ ] Implement MACD.
- [ ] Implement ATR.
- [ ] Add golden fixture tests.
- [ ] Add `GET /api/indicators`.
- [ ] Add chart overlays in frontend.
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

## 20. Future Features

### Auto-Trading

- Broker adapter interface.
- Paper trading first.
- Dry-run mode.
- Order preview.
- Manual approval mode.
- Kill switch.
- Max slippage guard.
- Execution audit log.

### Broker Integration

- MetaTrader bridge.
- Binance/Bybit crypto APIs.
- OANDA/Interactive Brokers later.
- Separate execution service from analysis service.

### Portfolio Optimization

- Exposure by asset class.
- Correlation matrix.
- Volatility-adjusted sizing.
- Risk parity experiments.
- Strategy allocation by recent expectancy.

### AI Self-Review

- Review every signal after N candles.
- Compare thesis vs actual outcome.
- Tag failure modes.
- Track agent confidence calibration.
- Generate weekly strategy review.
- Recommend rule changes, but require human approval.

### Advanced Intelligence

- News sentiment.
- Economic calendar awareness.
- Market regime classifier.
- Session behavior analysis.
- Liquidity sweep detection.
- Multi-agent debate: bull case, bear case, risk officer.

## 21. Test Strategy

### Backend

- Unit tests for indicators, risk rules, backtest engine, and structure detection.
- Repository integration tests with test Postgres.
- API tests for core endpoints.
- Worker tests with fake queues/providers.
- Agent workflow tests with mocked AI output.

### Frontend

- Component tests for risk calculator, signal cards, and journal forms.
- Playwright tests for key flows:
  - View chart.
  - Run analysis.
  - Inspect signal.
  - Create trade journal entry.
  - Run backtest.

### Data Correctness

- Golden test fixtures for indicators.
- Backtest no-lookahead tests.
- Risk calculation snapshot tests.
- Migration tests.

## 22. Practical Build Recommendation

Start with a modular monolith:

- One Go API.
- One Go worker.
- One Next.js app.
- One TimescaleDB database.
- Redis for cache.
- RabbitMQ for async jobs.

Do not start with microservices. The domain complexity is already high; the first architectural goal should be correctness, observability, and clean boundaries. Once the system proves useful, split market ingestion, agent runtime, and execution into separate services.
