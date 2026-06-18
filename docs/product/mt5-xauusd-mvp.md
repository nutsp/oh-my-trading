# MT5 XAUUSD MVP

## Goal

Build the first usable MVP around **XAUUSD only**, using MT5 as the market/account data source, with **read-only ingestion** and **paper signals** before any real execution.

## MVP Boundary

### Included

- XAUUSD only.
- MT5 terminal as the source of market/account truth.
- Python bridge reads from MT5 and sends normalized data to Go.
- Go service persists MT5 candles, ticks, account snapshots, positions, and bridge heartbeats.
- Dashboard shows MT5 bridge status, XAUUSD chart, latest tick/spread, account snapshot, open positions, and paper signals.
- AI/paper signal flow can propose trade ideas but cannot execute orders.
- Risk validation is required before a paper signal can become a paper trade idea.

### Excluded

- Real-money execution.
- MT5 order placement.
- EA WebSocket command channel.
- Multi-symbol expansion.
- Broker abstraction across multiple platforms.
- Portfolio optimization.
- Agent self-improvement loops.

## Recommended MVP Flow

```text
MT5 Terminal
  <-> Python Bridge
      -> Go REST ingest
Go Service
  -> PostgreSQL/TimescaleDB
  -> XAUUSD paper signal engine
Next.js Dashboard
  -> MT5 bridge status
  -> XAUUSD chart
  -> account/position snapshot
  -> paper signals
```

## Why Python Bridge First

Use Python first because it is the fastest reliable bridge to MT5 terminal data.

Benefits:

- Easier MT5 integration through the `MetaTrader5` Python package.
- Faster debugging than MQL5 WebSocket work.
- Keeps Go as the source of truth for storage, risk, signal, and dashboard APIs.
- Lets the project prove data quality before execution risk is introduced.

Tradeoffs:

- Python bridge must run on the same machine/VPS as MT5 terminal.
- It is not ideal for low-latency order execution.
- Needs heartbeat and restart monitoring.

## Later EA/WebSocket Path

After the read-only MVP is stable:

```text
Go creates approved command
  -> Python bridge or MT5 EA polls command
  -> MT5 executes in demo/paper mode
  -> execution result returns to Go
```

Only after that should we consider:

```text
MT5 EA <-> WebSocket <-> Go Service
```

Use EA/WebSocket when we need lower-latency event streaming or bidirectional command handling. It should not be the first integration path.

## Paper Signal Rules

- Every signal is paper-only by default.
- Every signal must include entry, stop-loss, take-profit, confidence, rationale, invalidation, and risk result.
- Signals must be tied to the latest MT5 account snapshot where possible.
- No signal can become an executable command in this MVP.
- A paper signal can be marked as accepted, rejected, expired, or reviewed.

## Success Criteria

- Python bridge connects to MT5 and reports healthy heartbeat.
- Go receives and stores XAUUSD candles and ticks.
- Dashboard shows latest XAUUSD candles from MT5.
- Dashboard shows latest bid, ask, spread, account equity, balance, margin, and open positions.
- Paper signal endpoint can create a non-executable XAUUSD signal.
- All ingest endpoints are idempotent or safe to retry.

