import unittest

from mt5_bridge.config import BridgeConfig


class BridgeConfigTest(unittest.TestCase):
    def test_from_env_uses_defaults(self) -> None:
        config = BridgeConfig.from_env({})

        self.assertEqual(config.api_url, "http://localhost:8080")
        self.assertEqual(config.bridge_id, "local-mt5")
        self.assertEqual(config.symbol, "XAUUSD")
        self.assertEqual(config.timeframes, ("1m", "5m", "15m", "1h"))
        self.assertEqual(config.poll_seconds, 2.0)

    def test_from_env_overrides_values(self) -> None:
        config = BridgeConfig.from_env(
            {
                "OMT_API_URL": "http://api.local/",
                "OMT_MT5_BRIDGE_ID": "desk-mt5",
                "OMT_MT5_SYMBOL": "XAUUSD",
                "OMT_MT5_TIMEFRAMES": "1m, 15m",
                "OMT_MT5_POLL_SECONDS": "5",
            }
        )

        self.assertEqual(config.api_url, "http://api.local")
        self.assertEqual(config.bridge_id, "desk-mt5")
        self.assertEqual(config.timeframes, ("1m", "15m"))
        self.assertEqual(config.poll_seconds, 5.0)


if __name__ == "__main__":
    unittest.main()
