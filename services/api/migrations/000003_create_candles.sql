CREATE TABLE IF NOT EXISTS candles (
  symbol_id uuid NOT NULL REFERENCES symbols(id),
  timeframe text NOT NULL,
  ts timestamptz NOT NULL,
  open numeric NOT NULL,
  high numeric NOT NULL,
  low numeric NOT NULL,
  close numeric NOT NULL,
  volume numeric,
  PRIMARY KEY (symbol_id, timeframe, ts)
);

CREATE INDEX IF NOT EXISTS candles_symbol_timeframe_ts_desc_idx
  ON candles (symbol_id, timeframe, ts DESC);

SELECT create_hypertable('candles', 'ts', if_not_exists => TRUE);

