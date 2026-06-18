from __future__ import annotations

from dataclasses import dataclass
from os import environ
from typing import Mapping


@dataclass(frozen=True)
class BridgeConfig:
    api_url: str = "http://localhost:8080"
    bridge_id: str = "local-mt5"
    symbol: str = "XAUUSD"
    timeframes: tuple[str, ...] = ("1m", "5m", "15m", "1h")
    poll_seconds: float = 2.0

    @classmethod
    def from_env(cls, env: Mapping[str, str] | None = None) -> "BridgeConfig":
        values = environ if env is None else env
        timeframes = tuple(
            item.strip()
            for item in values.get("OMT_MT5_TIMEFRAMES", "1m,5m,15m,1h").split(",")
            if item.strip()
        )
        return cls(
            api_url=values.get("OMT_API_URL", cls.api_url).rstrip("/"),
            bridge_id=values.get("OMT_MT5_BRIDGE_ID", cls.bridge_id),
            symbol=values.get("OMT_MT5_SYMBOL", cls.symbol),
            timeframes=timeframes or cls.timeframes,
            poll_seconds=float(values.get("OMT_MT5_POLL_SECONDS", cls.poll_seconds)),
        )
