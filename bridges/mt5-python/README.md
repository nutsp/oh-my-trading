# MT5 Python Bridge

Read-only bridge for the XAUUSD MVP.

The bridge is intentionally small:

- Connect to the local MT5 terminal through the optional `MetaTrader5` package.
- Pull XAUUSD ticks, candles, account state, and open positions.
- Normalize payloads.
- Send them to the Go API MT5 ingest endpoints.
- Support dry-run mode without connecting to MT5.

## Dry Run

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m mt5_bridge --dry-run
```

## Post Sample Payloads

Use this after the Go API is running locally. It posts sample heartbeat, tick, candle, account, and position payloads without connecting to MT5.

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m mt5_bridge --post-sample
```

## Tests

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m unittest discover -s tests
```

## Environment

```bash
OMT_API_URL=http://localhost:8080
OMT_MT5_BRIDGE_ID=local-mt5
OMT_MT5_SYMBOL=XAUUSD
OMT_MT5_TIMEFRAMES=1m,5m,15m,1h
OMT_MT5_POLL_SECONDS=2
```

Only `XAUUSD` is supported in this MVP. The bridge must not place orders.
