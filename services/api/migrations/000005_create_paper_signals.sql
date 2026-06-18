CREATE TABLE IF NOT EXISTS paper_signals (
  id uuid PRIMARY KEY,
  symbol text NOT NULL,
  timeframe text NOT NULL,
  side text NOT NULL,
  status text NOT NULL,
  confidence numeric NOT NULL,
  entry_price numeric NOT NULL,
  stop_loss numeric NOT NULL,
  take_profit numeric NOT NULL,
  thesis text NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS paper_signals_symbol_created_desc_idx
  ON paper_signals (symbol, created_at DESC);

CREATE INDEX IF NOT EXISTS paper_signals_status_created_desc_idx
  ON paper_signals (status, created_at DESC);
