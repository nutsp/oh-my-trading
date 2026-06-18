# AI Trading Agent Dashboard Architecture

## Product Vision

Build a personal AI trading command center for market analysis, AI signal generation, risk control, trade journaling, backtesting, and agent monitoring across XAUUSD, BTC, Forex, and Crypto.

The system should answer:

- What is the market doing across multiple timeframes?
- What setups does the AI agent see?
- Is the trade valid under the risk rules?
- What is the risk before entry?
- What happened after execution?
- Is the agent behaving correctly?
- Which strategies perform best historically?

The product should prioritize decision support before automation. Auto-trading should come later after analysis, journaling, risk checks, and backtesting are trustworthy.

## Core Modules

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

## MVP Scope

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

## System Architecture

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

## Technical Decisions

- **Go backend:** Strong concurrency, reliable services, simple deployment.
- **Next.js frontend:** Dashboard routing, React ecosystem, and chart integration.
- **TimescaleDB:** Best fit for candle and indicator time-series data.
- **Redis:** Cache latest candles, indicators, and agent states.
- **RabbitMQ:** Simpler than Kafka for a personal workflow queue.
- **Kafka later:** Useful only if event volume or stream replay requirements grow significantly.
- **TradingView Lightweight Charts:** Fast and clean trading UX.
- **Custom agent workflow:** Easier to own and debug than adopting LangGraph concepts wholesale.

## Agent Workflow Design

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

## Frontend Page Structure

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

## Backend Service Structure

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

## Risk Management Rules

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

## Backtesting Design

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

## Trade Journal Design

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

## Notification Design

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

## Security Considerations

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

## Deployment Plan

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

## Recommended Folder Structure

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

## Future Features

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

## Test Strategy

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

