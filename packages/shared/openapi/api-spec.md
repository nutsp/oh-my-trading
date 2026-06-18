# AI Trading Agent Dashboard API Spec

## API Style

Use REST for CRUD/query operations and WebSocket for live updates.

Base path:

```text
/api
```

Use JSON for requests and responses.

## Market Data

### List Symbols

```http
GET /api/symbols
```

Response:

```json
[
  {
    "id": "018f4f8a-0000-7000-9000-000000000001",
    "code": "XAUUSD",
    "market": "forex",
    "baseAsset": "XAU",
    "quoteAsset": "USD",
    "enabled": true
  }
]
```

### Create Symbol

```http
POST /api/symbols
```

Request:

```json
{
  "code": "BTCUSD",
  "market": "crypto",
  "baseAsset": "BTC",
  "quoteAsset": "USD"
}
```

### Query Candles

```http
GET /api/candles?symbol=XAUUSD&timeframe=1h&from=2026-01-01T00:00:00Z&to=2026-01-31T00:00:00Z
```

Response:

```json
[
  {
    "timestamp": "2026-01-01T00:00:00Z",
    "open": 2320.1,
    "high": 2328.3,
    "low": 2318.6,
    "close": 2325.5,
    "volume": 12345
  }
]
```

### Trigger Market Data Sync

```http
POST /api/market-data/sync
```

Request:

```json
{
  "symbol": "XAUUSD",
  "timeframes": ["1h", "4h", "1d"],
  "from": "2026-01-01T00:00:00Z",
  "to": "2026-01-31T00:00:00Z"
}
```

### Query Indicators

```http
GET /api/indicators?symbol=BTCUSD&timeframe=4h
```

Response:

```json
{
  "symbol": "BTCUSD",
  "timeframe": "4h",
  "series": {
    "ema20": [
      {
        "timestamp": "2026-01-01T08:00:00Z",
        "value": 2308.2
      }
    ],
    "ema50": [
      {
        "timestamp": "2026-01-01T08:00:00Z",
        "value": 2287.4
      }
    ],
    "rsi14": [
      {
        "timestamp": "2026-01-01T08:00:00Z",
        "value": 58.1
      }
    ],
    "macd": [
      {
        "timestamp": "2026-01-01T08:00:00Z",
        "macd": 6.2,
        "signal": 5.7,
        "histogram": 0.5
      }
    ],
    "atr14": [
      {
        "timestamp": "2026-01-01T08:00:00Z",
        "value": 44.6
      }
    ]
  }
}
```

## MT5 Read-Only Ingest

The first MT5 MVP is XAUUSD-only and read-only. These endpoints must not execute orders.

### MT5 Heartbeat

```http
POST /api/mt5/heartbeat
```

Request:

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

### MT5 Tick Ingest

```http
POST /api/mt5/ticks
```

Request:

```json
{
  "ticks": [
    {
      "symbol": "XAUUSD",
      "bid": 2325.42,
      "ask": 2325.62,
      "last": 2325.52,
      "volume": 12,
      "time": "2026-06-18T11:00:00Z"
    }
  ]
}
```

### MT5 Candle Ingest

```http
POST /api/mt5/candles
```

Request:

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

### MT5 Account Snapshot

```http
POST /api/mt5/account-snapshot
```

Request:

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

### MT5 Positions Snapshot

```http
POST /api/mt5/positions
```

Request:

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

### MT5 Status

```http
GET /api/mt5/status
```

Response:

```json
{
  "state": "connected",
  "heartbeat": {
    "bridgeId": "local-mt5",
    "terminal": "MetaTrader 5",
    "accountLogin": "12345678",
    "server": "Broker-Demo",
    "status": "healthy",
    "sentAt": "2026-06-18T11:00:00Z"
  },
  "latestTick": {
    "symbol": "XAUUSD",
    "bid": 2325.42,
    "ask": 2325.62,
    "last": 2325.52,
    "volume": 12,
    "time": "2026-06-18T11:00:00Z"
  }
}
```

Before the first bridge payload arrives, the endpoint still returns `200 OK`:

```json
{
  "state": "waiting_for_bridge",
  "heartbeat": {
    "bridgeId": "local-mt5",
    "status": "disconnected",
    "sentAt": ""
  },
  "latestTick": {
    "symbol": "XAUUSD",
    "bid": 0,
    "ask": 0,
    "last": 0,
    "volume": 0,
    "time": ""
  }
}
```

### Latest MT5 Account

```http
GET /api/mt5/account/latest?accountLogin=12345678
```

### Latest MT5 Positions

```http
GET /api/mt5/positions/latest?accountLogin=12345678
```

## Analysis And Signals

### Run Analysis

```http
POST /api/analysis/run
```

Request:

```json
{
  "symbol": "XAUUSD",
  "timeframes": ["1d", "4h", "1h", "15m"],
  "strategyProfile": "top_down_ema_rsi_macd_atr_smc"
}
```

Response:

```json
{
  "agentRunId": "018f4f8a-0000-7000-9000-000000000101",
  "signalId": "018f4f8a-0000-7000-9000-000000000201",
  "status": "succeeded"
}
```

### List Analysis Runs

```http
GET /api/analysis/runs
```

### Create Paper Signal

```http
POST /api/paper-signals
```

Request:

```json
{
  "symbol": "XAUUSD",
  "timeframe": "1h",
  "side": "long",
  "confidence": 0.72,
  "entryPrice": 2325.5,
  "stopLoss": 2310,
  "takeProfit": 2340,
  "thesis": "Top-down bullish continuation with EMA alignment."
}
```

Response status: `201 Created`.

### List Paper Signals

```http
GET /api/paper-signals
```

### Update Paper Signal Status

```http
PATCH /api/paper-signals/{id}/status
```

Request:

```json
{
  "status": "approved_paper"
}
```

Supported status values after creation: `approved_paper`, `rejected`, `expired`.

### List Signals

```http
GET /api/signals?symbol=XAUUSD&status=new
```

Response:

```json
[
  {
    "id": "018f4f8a-0000-7000-9000-000000000201",
    "symbol": "XAUUSD",
    "timeframe": "1h",
    "direction": "long",
    "confidence": 0.72,
    "status": "new",
    "createdAt": "2026-01-01T12:00:00Z"
  }
]
```

### Get Signal

```http
GET /api/signals/{id}
```

Response:

```json
{
  "id": "018f4f8a-0000-7000-9000-000000000201",
  "symbol": "XAUUSD",
  "timeframe": "1h",
  "direction": "long",
  "confidence": 0.72,
  "entry": 2325.5,
  "stopLoss": 2314.2,
  "takeProfit": 2350.0,
  "status": "new",
  "reason": "Daily bullish structure, 4H pullback into EMA zone, RSI recovery, ATR supports stop distance.",
  "invalidation": "Close below 2314.2 or bearish CHoCH on 1H."
}
```

### Update Signal Status

```http
PATCH /api/signals/{id}/status
```

Request:

```json
{
  "status": "accepted"
}
```

## Risk

### Calculate Position

```http
POST /api/risk/calculate-position
```

Request:

```json
{
  "accountEquity": 10000,
  "riskPercent": 1,
  "entry": 2325.5,
  "stopLoss": 2314.2,
  "symbol": "XAUUSD"
}
```

Response:

```json
{
  "positionSize": 0.08,
  "riskAmount": 100,
  "stopDistance": 11.3
}
```

### Validate Trade

```http
POST /api/risk/validate-trade
```

Request:

```json
{
  "symbol": "XAUUSD",
  "direction": "long",
  "entry": 2325.5,
  "stopLoss": 2314.2,
  "takeProfit": 2350,
  "riskPercent": 1
}
```

Response:

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

### Get Risk Settings

```http
GET /api/risk/settings
```

### Update Risk Settings

```http
PUT /api/risk/settings
```

## Journal

### List Trades

```http
GET /api/trades?status=active&symbol=XAUUSD
```

### Create Trade

```http
POST /api/trades
```

Request:

```json
{
  "signalId": "018f4f8a-0000-7000-9000-000000000201",
  "symbol": "XAUUSD",
  "direction": "long",
  "plannedEntry": 2325.5,
  "stopLoss": 2314.2,
  "takeProfit": 2350,
  "positionSize": 0.08,
  "riskAmount": 100
}
```

### Get Trade

```http
GET /api/trades/{id}
```

### Update Trade

```http
PATCH /api/trades/{id}
```

### Add Trade Note

```http
POST /api/trades/{id}/notes
```

Request:

```json
{
  "noteType": "review",
  "content": "Entry followed the plan, but exit was early.",
  "tags": ["early-exit", "discipline"]
}
```

## Backtesting

### Create Backtest

```http
POST /api/backtests
```

Request:

```json
{
  "name": "XAUUSD 1H EMA RSI SMC",
  "strategy": "top_down_ema_rsi_macd_atr_smc",
  "symbol": "XAUUSD",
  "timeframe": "1h",
  "from": "2025-01-01T00:00:00Z",
  "to": "2025-12-31T23:59:59Z",
  "config": {
    "riskPercent": 1,
    "minRewardRisk": 1.5,
    "emaFast": 20,
    "emaSlow": 50,
    "rsiPeriod": 14,
    "atrPeriod": 14
  }
}
```

### List Backtests

```http
GET /api/backtests
```

### Get Backtest

```http
GET /api/backtests/{id}
```

### List Backtest Trades

```http
GET /api/backtests/{id}/trades
```

## Monitoring

### List Agent Runs

```http
GET /api/agent-runs
```

### Get Agent Run

```http
GET /api/agent-runs/{id}
```

### Health

```http
GET /api/health
```

Response:

```json
{
  "status": "ok"
}
```

### Metrics

```http
GET /api/metrics
```

## WebSocket Channels

### Market Updates

```text
/ws/market
```

Event:

```json
{
  "type": "candle.closed",
  "symbol": "XAUUSD",
  "timeframe": "1h",
  "candle": {
    "timestamp": "2026-01-01T12:00:00Z",
    "open": 2320.1,
    "high": 2328.3,
    "low": 2318.6,
    "close": 2325.5
  }
}
```

### Signal Updates

```text
/ws/signals
```

Event:

```json
{
  "type": "signal.created",
  "signalId": "018f4f8a-0000-7000-9000-000000000201",
  "symbol": "XAUUSD",
  "direction": "long",
  "confidence": 0.72
}
```

### Agent Run Updates

```text
/ws/agent-runs
```

Event:

```json
{
  "type": "agent_run.failed",
  "agentRunId": "018f4f8a-0000-7000-9000-000000000101",
  "workflow": "top_down_analysis",
  "error": "AI provider timeout"
}
```
