from datetime import UTC, datetime
import unittest

from mt5_bridge.payloads import (
    build_account_snapshot,
    build_candles,
    build_heartbeat,
    build_positions,
    build_ticks,
)


class PayloadBuilderTest(unittest.TestCase):
    def test_builds_heartbeat_payload(self) -> None:
        payload = build_heartbeat(
            bridge_id="local-mt5",
            terminal="MetaTrader 5",
            account_login="12345678",
            server="Broker-Demo",
            status="healthy",
            sent_at=datetime(2026, 6, 18, 10, 0, tzinfo=UTC),
        )

        self.assertEqual(payload["bridgeId"], "local-mt5")
        self.assertEqual(payload["sentAt"], "2026-06-18T10:00:00Z")

    def test_builds_tick_payload(self) -> None:
        payload = build_ticks(
            [
                {
                    "symbol": "XAUUSD",
                    "bid": 2325.42,
                    "ask": 2325.62,
                    "last": 2325.52,
                    "volume": 12,
                    "time": 1_782_986_400,
                }
            ]
        )

        self.assertEqual(payload["ticks"][0]["symbol"], "XAUUSD")
        self.assertEqual(payload["ticks"][0]["ask"], 2325.62)

    def test_builds_candle_account_and_positions_payloads(self) -> None:
        now = datetime(2026, 6, 18, 10, 0, tzinfo=UTC)

        candles = build_candles(
            symbol="XAUUSD",
            timeframe="1m",
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
        )
        account = build_account_snapshot(
            {
                "account_login": "12345678",
                "currency": "USD",
                "balance": 10000,
                "equity": 10080,
                "time": now,
            }
        )
        positions = build_positions(
            account_login="12345678",
            time=now,
            positions=[
                {
                    "ticket": 987654321,
                    "symbol": "XAUUSD",
                    "side": "buy",
                    "volume": 0.1,
                    "open_price": 2320.1,
                    "profit": 55,
                    "opened_at": now,
                }
            ],
        )

        self.assertEqual(candles["candles"][0]["close"], 2325.5)
        self.assertEqual(account["freeMargin"], 0)
        self.assertEqual(positions["positions"][0]["ticket"], "987654321")


if __name__ == "__main__":
    unittest.main()
