from __future__ import annotations

from datetime import UTC, datetime
from typing import Any, Protocol


class MT5Adapter(Protocol):
    def terminal_name(self) -> str:
        ...

    def account_info(self) -> dict[str, Any]:
        ...

    def latest_tick(self, symbol: str) -> dict[str, Any]:
        ...

    def rates(self, symbol: str, timeframe: str, count: int) -> list[dict[str, Any]]:
        ...

    def positions(self, symbol: str) -> list[dict[str, Any]]:
        ...


class MetaTrader5Adapter:
    def __init__(self) -> None:
        try:
            import MetaTrader5 as mt5
        except ImportError as exc:
            raise RuntimeError("Install the optional MetaTrader5 dependency to use the live bridge.") from exc
        self.mt5 = mt5
        if not self.mt5.initialize():
            raise RuntimeError(f"MT5 initialize failed: {self.mt5.last_error()}")

    def terminal_name(self) -> str:
        info = self.mt5.terminal_info()
        return getattr(info, "name", "MetaTrader 5")

    def account_info(self) -> dict[str, Any]:
        info = self.mt5.account_info()
        if info is None:
            raise RuntimeError(f"MT5 account_info failed: {self.mt5.last_error()}")
        return {
            "account_login": str(info.login),
            "currency": info.currency,
            "balance": info.balance,
            "equity": info.equity,
            "margin": info.margin,
            "free_margin": info.margin_free,
            "margin_level": info.margin_level,
            "server": info.server,
            "time": datetime.now(UTC),
        }

    def latest_tick(self, symbol: str) -> dict[str, Any]:
        tick = self.mt5.symbol_info_tick(symbol)
        if tick is None:
            raise RuntimeError(f"MT5 symbol_info_tick failed for {symbol}: {self.mt5.last_error()}")
        return {
            "symbol": symbol,
            "bid": tick.bid,
            "ask": tick.ask,
            "last": tick.last,
            "volume": tick.volume,
            "time": tick.time,
        }

    def rates(self, symbol: str, timeframe: str, count: int) -> list[dict[str, Any]]:
        mt5_timeframe = self._timeframe(timeframe)
        rates = self.mt5.copy_rates_from_pos(symbol, mt5_timeframe, 0, count)
        if rates is None:
            raise RuntimeError(f"MT5 copy_rates_from_pos failed for {symbol} {timeframe}: {self.mt5.last_error()}")
        return [
            {
                "timestamp": int(rate["time"]),
                "open": rate["open"],
                "high": rate["high"],
                "low": rate["low"],
                "close": rate["close"],
                "volume": rate["tick_volume"],
            }
            for rate in rates
        ]

    def positions(self, symbol: str) -> list[dict[str, Any]]:
        positions = self.mt5.positions_get(symbol=symbol)
        if positions is None:
            raise RuntimeError(f"MT5 positions_get failed for {symbol}: {self.mt5.last_error()}")
        return [
            {
                "ticket": position.ticket,
                "symbol": position.symbol,
                "side": "buy" if position.type == self.mt5.POSITION_TYPE_BUY else "sell",
                "volume": position.volume,
                "open_price": position.price_open,
                "stop_loss": position.sl,
                "take_profit": position.tp,
                "profit": position.profit,
                "opened_at": position.time,
            }
            for position in positions
        ]

    def shutdown(self) -> None:
        self.mt5.shutdown()

    def _timeframe(self, timeframe: str) -> int:
        mapping = {
            "1m": self.mt5.TIMEFRAME_M1,
            "5m": self.mt5.TIMEFRAME_M5,
            "15m": self.mt5.TIMEFRAME_M15,
            "1h": self.mt5.TIMEFRAME_H1,
            "4h": self.mt5.TIMEFRAME_H4,
            "1d": self.mt5.TIMEFRAME_D1,
        }
        try:
            return mapping[timeframe]
        except KeyError as exc:
            raise ValueError(f"unsupported timeframe: {timeframe}") from exc
