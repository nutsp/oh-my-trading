from __future__ import annotations

import json
from datetime import UTC, datetime
from typing import Any, TextIO

from .client import GoAPIClient
from .config import BridgeConfig
from .mt5_adapter import MT5Adapter
from .payloads import (
    build_account_snapshot,
    build_candles,
    build_heartbeat,
    build_positions,
    build_ticks,
)


def run_once(adapter: MT5Adapter, config: BridgeConfig, client: GoAPIClient) -> None:
    now = datetime.now(UTC)
    account = adapter.account_info()
    account_login = str(account["account_login"])
    server = str(account.get("server", ""))

    client.post_json(
        "/api/mt5/heartbeat",
        build_heartbeat(
            bridge_id=config.bridge_id,
            terminal=adapter.terminal_name(),
            account_login=account_login,
            server=server,
            status="healthy",
            sent_at=now,
        ),
    )
    client.post_json("/api/mt5/ticks", build_ticks([adapter.latest_tick(config.symbol)]))
    for timeframe in config.timeframes:
        client.post_json(
            "/api/mt5/candles",
            build_candles(
                symbol=config.symbol,
                timeframe=timeframe,
                source="mt5-python-bridge",
                candles=adapter.rates(config.symbol, timeframe, count=100),
            ),
        )
    client.post_json("/api/mt5/account-snapshot", build_account_snapshot(account))
    client.post_json(
        "/api/mt5/positions",
        build_positions(
            account_login=account_login,
            positions=adapter.positions(config.symbol),
            time=now,
        ),
    )


def print_dry_run(config: BridgeConfig, out: TextIO) -> None:
    now = datetime(2026, 6, 18, 10, 0, 0, tzinfo=UTC)
    payloads: dict[str, dict[str, Any]] = {
        "/api/mt5/heartbeat": build_heartbeat(
            bridge_id=config.bridge_id,
            terminal="MetaTrader 5",
            account_login="12345678",
            server="Broker-Demo",
            status="healthy",
            sent_at=now,
        ),
        "/api/mt5/ticks": build_ticks(
            [
                {
                    "symbol": config.symbol,
                    "bid": 2325.42,
                    "ask": 2325.62,
                    "last": 2325.52,
                    "volume": 12,
                    "time": now,
                }
            ]
        ),
        "/api/mt5/candles": build_candles(
            symbol=config.symbol,
            timeframe=config.timeframes[0],
            source="mt5-python-bridge",
            candles=[
                {
                    "timestamp": now,
                    "open": 2320.1,
                    "high": 2328.3,
                    "low": 2318.6,
                    "close": 2325.5,
                    "volume": 12345,
                }
            ],
        ),
        "/api/mt5/account-snapshot": build_account_snapshot(
            {
                "account_login": "12345678",
                "currency": "USD",
                "balance": 10000,
                "equity": 10080,
                "margin": 400,
                "free_margin": 9680,
                "margin_level": 2520,
                "time": now,
            }
        ),
        "/api/mt5/positions": build_positions(
            account_login="12345678",
            positions=[
                {
                    "ticket": "987654321",
                    "symbol": config.symbol,
                    "side": "buy",
                    "volume": 0.1,
                    "open_price": 2320.1,
                    "stop_loss": 2310,
                    "take_profit": 2340,
                    "profit": 55,
                    "opened_at": now,
                }
            ],
            time=now,
        ),
    }
    for path, payload in payloads.items():
        out.write(f"POST {path}\n")
        out.write(json.dumps(payload, indent=2, sort_keys=True))
        out.write("\n")
