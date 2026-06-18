CREATE TABLE IF NOT EXISTS symbols (
  id uuid PRIMARY KEY,
  code text UNIQUE NOT NULL,
  market text NOT NULL,
  base_asset text,
  quote_asset text,
  enabled boolean NOT NULL DEFAULT true
);

