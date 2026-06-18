import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

test("dashboard shell exposes primary navigation", async () => {
  const shell = await readFile(new URL("../components/app-shell.tsx", import.meta.url), "utf8");

  for (const label of [
    "Dashboard",
    "Markets",
    "MT5",
    "Signals",
    "Risk",
    "Journal",
    "Backtests",
    "Agent",
    "Settings",
  ]) {
    assert.match(shell, new RegExp(label));
  }
});

test("api client points to backend api by default", async () => {
  const client = await readFile(new URL("../lib/api-client.ts", import.meta.url), "utf8");

  assert.match(client, /NEXT_PUBLIC_API_BASE_URL/);
  assert.match(client, /http:\/\/localhost:8080/);
  assert.match(client, /\/api\/indicators/);
  assert.match(client, /\/api\/mt5\/status/);
});

test("markets route is wired with timeframe selector", async () => {
  const marketPage = await readFile(
    new URL("../features/markets/market-chart-page.tsx", import.meta.url),
    "utf8",
  );
  const route = await readFile(new URL("../app/markets/[symbol]/page.tsx", import.meta.url), "utf8");

  for (const timeframe of ["15m", "1h", "4h", "1d"]) {
    assert.match(marketPage, new RegExp(timeframe));
  }
  assert.match(route, /MarketChartPage/);
});

test("mt5 route is wired to bridge dashboard", async () => {
  const route = await readFile(new URL("../app/mt5/page.tsx", import.meta.url), "utf8");
  const dashboard = await readFile(new URL("../features/mt5/mt5-dashboard.tsx", import.meta.url), "utf8");

  assert.match(route, /MT5Dashboard/);
  assert.match(dashboard, /getMT5Status/);
  assert.match(dashboard, /Open Positions/);
});
