CREATE TABLE IF NOT EXISTS mt5_bridge_heartbeats (
  bridge_id text PRIMARY KEY,
  terminal text NOT NULL,
  account_login text NOT NULL,
  server text NOT NULL,
  status text NOT NULL,
  last_error text,
  sent_at timestamptz NOT NULL,
  received_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS mt5_ticks (
  id bigserial,
  symbol text NOT NULL,
  bid numeric NOT NULL,
  ask numeric NOT NULL,
  last numeric,
  volume numeric,
  ts timestamptz NOT NULL,
  received_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS mt5_ticks_symbol_ts_bid_ask_idx
  ON mt5_ticks (symbol, ts, bid, ask);

CREATE INDEX IF NOT EXISTS mt5_ticks_symbol_ts_desc_idx
  ON mt5_ticks (symbol, ts DESC);

SELECT create_hypertable('mt5_ticks', 'ts', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS mt5_account_snapshots (
  id bigserial,
  account_login text NOT NULL,
  currency text NOT NULL,
  balance numeric NOT NULL,
  equity numeric NOT NULL,
  margin numeric NOT NULL,
  free_margin numeric NOT NULL,
  margin_level numeric,
  ts timestamptz NOT NULL,
  received_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS mt5_account_snapshots_login_ts_idx
  ON mt5_account_snapshots (account_login, ts);

CREATE INDEX IF NOT EXISTS mt5_account_snapshots_login_ts_desc_idx
  ON mt5_account_snapshots (account_login, ts DESC);

SELECT create_hypertable('mt5_account_snapshots', 'ts', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS mt5_position_snapshots (
  id bigserial,
  account_login text NOT NULL,
  ticket text NOT NULL,
  symbol text NOT NULL,
  side text NOT NULL,
  volume numeric NOT NULL,
  open_price numeric NOT NULL,
  stop_loss numeric,
  take_profit numeric,
  profit numeric NOT NULL,
  opened_at timestamptz NOT NULL,
  snapshot_ts timestamptz NOT NULL,
  received_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS mt5_position_snapshots_login_ticket_snapshot_idx
  ON mt5_position_snapshots (account_login, ticket, snapshot_ts);

CREATE INDEX IF NOT EXISTS mt5_position_snapshots_login_snapshot_desc_idx
  ON mt5_position_snapshots (account_login, snapshot_ts DESC);

SELECT create_hypertable('mt5_position_snapshots', 'snapshot_ts', if_not_exists => TRUE);
