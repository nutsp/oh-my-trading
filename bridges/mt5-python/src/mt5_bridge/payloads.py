from __future__ import annotations

from datetime import UTC, datetime
from typing import Any, Iterable, Mapping


def utc_iso(value: datetime) -> str:
    if value.tzinfo is None:
        value = value.replace(tzinfo=UTC)
    return value.astimezone(UTC).isoformat().replace("+00:00", "Z")


def build_heartbeat(
    *,
    bridge_id: str,
    terminal: str,
    account_login: str,
    server: str,
    status: str,
    sent_at: datetime,
    last_error: str = "",
) -> dict[str, Any]:
    return {
        "bridgeId": bridge_id,
        "terminal": terminal,
        "accountLogin": account_login,
        "server": server,
        "status": status,
        "lastError": last_error,
        "sentAt": utc_iso(sent_at),
    }


def build_ticks(ticks: Iterable[Mapping[str, Any]]) -> dict[str, Any]:
    return {
        "ticks": [
            {
                "symbol": str(tick["symbol"]),
                "bid": float(tick["bid"]),
                "ask": float(tick["ask"]),
                "last": float(tick.get("last", 0)),
                "volume": float(tick.get("volume", 0)),
                "time": utc_iso(as_datetime(tick["time"])),
            }
            for tick in ticks
        ]
    }


def build_candles(*, symbol: str, timeframe: str, source: str, candles: Iterable[Mapping[str, Any]]) -> dict[str, Any]:
    return {
        "symbol": symbol,
        "timeframe": timeframe,
        "source": source,
        "candles": [
            {
                "timestamp": utc_iso(as_datetime(candle["timestamp"])),
                "open": float(candle["open"]),
                "high": float(candle["high"]),
                "low": float(candle["low"]),
                "close": float(candle["close"]),
                "volume": float(candle.get("volume", 0)),
            }
            for candle in candles
        ],
    }


def build_account_snapshot(account: Mapping[str, Any]) -> dict[str, Any]:
    return {
        "accountLogin": str(account["account_login"]),
        "currency": str(account.get("currency", "USD")),
        "balance": float(account["balance"]),
        "equity": float(account["equity"]),
        "margin": float(account.get("margin", 0)),
        "freeMargin": float(account.get("free_margin", 0)),
        "marginLevel": float(account.get("margin_level", 0)),
        "time": utc_iso(as_datetime(account["time"])),
    }


def build_positions(*, account_login: str, positions: Iterable[Mapping[str, Any]], time: datetime) -> dict[str, Any]:
    return {
        "accountLogin": account_login,
        "time": utc_iso(time),
        "positions": [
            {
                "ticket": str(position["ticket"]),
                "symbol": str(position["symbol"]),
                "side": str(position["side"]),
                "volume": float(position["volume"]),
                "openPrice": float(position["open_price"]),
                "stopLoss": float(position.get("stop_loss", 0)),
                "takeProfit": float(position.get("take_profit", 0)),
                "profit": float(position.get("profit", 0)),
                "openedAt": utc_iso(as_datetime(position["opened_at"])),
            }
            for position in positions
        ],
    }


def as_datetime(value: Any) -> datetime:
    if isinstance(value, datetime):
        return value
    if isinstance(value, (int, float)):
        return datetime.fromtimestamp(value, UTC)
    if isinstance(value, str):
        normalized = value.replace("Z", "+00:00")
        return datetime.fromisoformat(normalized)
    raise TypeError(f"unsupported datetime value: {value!r}")
