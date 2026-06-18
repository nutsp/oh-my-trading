# AI Trading Agent Dashboard Database Design

## Database Choice

Use PostgreSQL with the TimescaleDB extension.

Reasons:

- Trading candles and indicator values are time-series data.
- TimescaleDB keeps PostgreSQL ergonomics while adding hypertables, compression, retention policies, and time-window query performance.
- PostgreSQL remains a strong fit for journal, signal, agent run, and backtest metadata.

## Core Schema

### Symbols

Stores configured markets and instruments.

```sql
CREATE TABLE symbols (
  id uuid PRIMARY KEY,
  code text UNIQUE NOT NULL,
  market text NOT NULL,
  base_asset text,
  quote_asset text,
  enabled boolean NOT NULL DEFAULT true
);
```

Example values:

- `XAUUSD`, market `forex`
- `BTCUSD`, market `crypto`
- `EURUSD`, market `forex`

### Candles

Stores OHLCV market data per symbol and timeframe.

```sql
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

SELECT create_hypertable('candles', 'ts');
```

Recommended indexes:

```sql
CREATE INDEX candles_symbol_timeframe_ts_desc_idx
  ON candles (symbol_id, timeframe, ts DESC);
```

### Indicator Values

Stores computed indicators when persistence is useful for caching, backtest reproducibility, or UI overlays.

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
```

Example `value` for MACD:

```json
{
  "macd": 12.4,
  "signal": 10.1,
  "histogram": 2.3
}
```

### Market Structures

Stores detected technical structures.

```sql
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
```

Example `type` values:

- `swing_high`
- `swing_low`
- `bos`
- `choch`
- `fvg`
- `order_block`

### Signals

Stores AI-generated trade ideas.

```sql
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
```

Recommended `status` values:

- `new`
- `watching`
- `accepted`
- `rejected`
- `expired`
- `converted_to_trade`

### Trades

Stores planned, active, closed, and reviewed trades.

```sql
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
```

Trade lifecycle:

```text
planned -> active -> closed -> reviewed
```

### Trade Notes

Stores journal notes and review comments.

```sql
CREATE TABLE trade_notes (
  id uuid PRIMARY KEY,
  trade_id uuid NOT NULL,
  note_type text NOT NULL,
  content text NOT NULL,
  tags text[],
  created_at timestamptz NOT NULL
);
```

Example `note_type` values:

- `plan`
- `execution`
- `emotion`
- `mistake`
- `review`
- `ai_review`

### Agent Runs

Stores workflow executions for observability and review.

```sql
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
```

Recommended `status` values:

- `running`
- `succeeded`
- `failed`
- `cancelled`

### Backtests

Stores backtest run metadata and summary metrics.

```sql
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

Example `metrics`:

```json
{
  "totalReturn": 0.184,
  "winRate": 0.56,
  "profitFactor": 1.72,
  "maxDrawdown": 0.08,
  "expectancyR": 0.31,
  "tradeCount": 126
}
```

## Future Tables

Add these after MVP:

- `risk_settings`
- `notification_events`
- `strategy_versions`
- `backtest_trades`
- `broker_accounts`
- `orders`
- `positions`
- `audit_logs`
- `screenshots`

## Data Correctness Rules

- Candle inserts should be idempotent through `(symbol_id, timeframe, ts)` primary key.
- Backtests must store strategy config snapshots for reproducibility.
- Indicator calculations must avoid lookahead.
- Agent runs must store enough input/output to debug signal generation.
- Broker/API secrets must not be stored as plaintext.

