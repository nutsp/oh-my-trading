INSERT INTO symbols (id, code, market, base_asset, quote_asset, enabled)
VALUES ('018f4f8a-0000-7000-9000-000000000001', 'XAUUSD', 'forex', 'XAU', 'USD', true)
ON CONFLICT (code) DO UPDATE SET
  market = EXCLUDED.market,
  base_asset = EXCLUDED.base_asset,
  quote_asset = EXCLUDED.quote_asset,
  enabled = EXCLUDED.enabled;
