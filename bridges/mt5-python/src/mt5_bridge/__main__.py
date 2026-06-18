from __future__ import annotations

import argparse
import sys

from .client import GoAPIClient
from .config import BridgeConfig
from .mt5_adapter import MetaTrader5Adapter
from .runner import print_dry_run, run_once


def main() -> int:
    parser = argparse.ArgumentParser(description="Read-only MT5 bridge for the XAUUSD MVP.")
    parser.add_argument("--dry-run", action="store_true", help="print example payloads without connecting to MT5")
    parser.add_argument("--once", action="store_true", help="poll MT5 once and post to the Go API")
    args = parser.parse_args()

    config = BridgeConfig.from_env()
    if config.symbol != "XAUUSD":
        parser.error("this MVP only supports OMT_MT5_SYMBOL=XAUUSD")

    if args.dry_run:
        print_dry_run(config, sys.stdout)
        return 0

    if not args.once:
        parser.error("choose --dry-run or --once")

    adapter = MetaTrader5Adapter()
    try:
        run_once(adapter, config, GoAPIClient(config.api_url))
    finally:
        adapter.shutdown()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
