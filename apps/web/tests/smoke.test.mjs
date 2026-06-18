import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

test("dashboard shell exposes primary navigation", async () => {
  const shell = await readFile(new URL("../components/app-shell.tsx", import.meta.url), "utf8");

  for (const label of [
    "Dashboard",
    "Markets",
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
});

