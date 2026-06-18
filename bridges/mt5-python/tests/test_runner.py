from io import StringIO
import unittest

from mt5_bridge.config import BridgeConfig
from mt5_bridge.runner import post_sample, print_dry_run


class RunnerTest(unittest.TestCase):
    def test_dry_run_prints_expected_endpoints(self) -> None:
        out = StringIO()

        print_dry_run(BridgeConfig(), out)

        text = out.getvalue()
        self.assertIn("POST /api/mt5/heartbeat", text)
        self.assertIn("POST /api/mt5/ticks", text)
        self.assertIn("POST /api/mt5/candles", text)
        self.assertIn("POST /api/mt5/account-snapshot", text)
        self.assertIn("POST /api/mt5/positions", text)
        self.assertIn('"symbol": "XAUUSD"', text)

    def test_post_sample_sends_expected_payloads(self) -> None:
        client = FakeClient()

        post_sample(BridgeConfig(), client)

        self.assertEqual(
            [path for path, _payload in client.posts],
            [
                "/api/mt5/heartbeat",
                "/api/mt5/ticks",
                "/api/mt5/candles",
                "/api/mt5/account-snapshot",
                "/api/mt5/positions",
            ],
        )
        self.assertEqual(client.posts[1][1]["ticks"][0]["symbol"], "XAUUSD")


class FakeClient:
    def __init__(self) -> None:
        self.posts = []

    def post_json(self, path: str, payload: dict) -> None:
        self.posts.append((path, payload))


if __name__ == "__main__":
    unittest.main()
