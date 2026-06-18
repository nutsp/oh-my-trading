# MT5 Integration Architecture

## Direction

The first MT5 integration should be read-only and XAUUSD-only.

Primary path:

```text
MT5 Terminal
  <-> Python Bridge
      -> Go REST API
      -> PostgreSQL/TimescaleDB
      -> Dashboard
```

Secondary future path:

```text
MT5 EA
  <-> WebSocket
      <-> Go Service
```

Do not start with EA WebSocket unless the read-only bridge has already proven data quality and operational stability.

## Components

### MT5 Terminal

Owns:

- Broker connection.
- Live XAUUSD market data.
- Account state.
- Open positions.

### Python Bridge

Owns:

- Connecting to local MT5 terminal.
- Pulling XAUUSD ticks and rates.
- Normalizing payloads.
- Sending data to Go ingest endpoints.
- Sending heartbeat.
- Logging bridge errors.

The bridge should not own trading strategy, risk rules, signal decisions, or persistence.

### Go Service

Owns:

- Data contracts.
- Persistence.
- Idempotency.
- Risk rules.
- Paper signal lifecycle.
- Dashboard APIs.
- Future command queue.

### Dashboard

Owns:

- MT5 bridge health view.
- XAUUSD chart.
- Latest tick/spread.
- Account and position snapshots.
- Paper signal review.

## Data Contracts

### Heartbeat

```json
{
  "bridgeId": "local-mt5",
  "terminal": "MetaTrader 5",
  "accountLogin": "12345678",
  "server": "Broker-Demo",
  "status": "healthy",
  "sentAt": "2026-06-18T11:00:00Z"
}
```

### Tick

```json
{
  "symbol": "XAUUSD",
  "bid": 2325.42,
  "ask": 2325.62,
  "last": 2325.52,
  "volume": 12,
  "time": "2026-06-18T11:00:00Z"
}
```

### Candle Batch

```json
{
  "symbol": "XAUUSD",
  "timeframe": "1m",
  "source": "mt5-python-bridge",
  "candles": [
    {
      "timestamp": "2026-06-18T11:00:00Z",
      "open": 2320.1,
      "high": 2328.3,
      "low": 2318.6,
      "close": 2325.5,
      "volume": 12345
    }
  ]
}
```

### Account Snapshot

```json
{
  "accountLogin": "12345678",
  "currency": "USD",
  "balance": 10000,
  "equity": 10080,
  "margin": 400,
  "freeMargin": 9680,
  "marginLevel": 2520,
  "time": "2026-06-18T11:00:00Z"
}
```

### Position Snapshot

```json
{
  "accountLogin": "12345678",
  "positions": [
    {
      "ticket": "987654321",
      "symbol": "XAUUSD",
      "side": "buy",
      "volume": 0.1,
      "openPrice": 2320.1,
      "stopLoss": 2310,
      "takeProfit": 2340,
      "profit": 55,
      "openedAt": "2026-06-18T10:00:00Z"
    }
  ],
  "time": "2026-06-18T11:00:00Z"
}
```

## Go API Boundary

Recommended first endpoints:

- `POST /api/mt5/heartbeat`
- `POST /api/mt5/ticks`
- `POST /api/mt5/candles`
- `POST /api/mt5/account-snapshot`
- `POST /api/mt5/positions`
- `GET /api/mt5/status`
- `GET /api/mt5/account/latest`
- `GET /api/mt5/positions/latest`

Keep these separate from generic market data endpoints. The MT5 adapter owns MT5 protocol translation; domain/application code should speak trading language, not MQL/Python details.

## Database Additions

Recommended tables:

- `mt5_bridge_heartbeats`
- `mt5_ticks`
- `mt5_account_snapshots`
- `mt5_position_snapshots`
- `paper_signals`

Keep existing `candles` table for OHLCV. MT5 candle ingest should upsert into existing `candles` after resolving `XAUUSD` symbol ID.

## Operational Rules

- Bridge heartbeat older than 30 seconds is stale.
- Tick ingest should tolerate duplicates.
- Candle ingest should be idempotent.
- Account snapshots are append-only.
- Position snapshots are append-only until a normalized latest view is needed.
- No endpoint in this phase should execute an order.

