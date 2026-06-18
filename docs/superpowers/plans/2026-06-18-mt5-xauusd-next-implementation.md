# MT5 XAUUSD Next Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Turn the current MT5 XAUUSD read-only foundation into a reliable local-to-demo workflow with real MT5 data, visible data quality, paper signal review, and safe handoff toward future automation.

**Architecture:** Keep the Go service as the source of truth for persistence, validation, paper signal lifecycle, and dashboard APIs. Keep the Python bridge as a replaceable read-only adapter around the local MT5 terminal. Keep all order execution out of scope until paper signal review, risk gates, monitoring, audit logs, and explicit approval workflows are complete.

**Tech Stack:** Go, PostgreSQL/TimescaleDB, Next.js/React, Python 3.11, optional `MetaTrader5` Python package, TradingView Lightweight Charts, Docker Compose.

---

## Current Baseline

Completed:

- MT5 ingest schema and repositories exist.
- `XAUUSD` is seeded by migration.
- Go API has MT5 ingest endpoints and paper signal endpoints.
- Python bridge supports `--dry-run` and `--post-sample`.
- Dashboard has `/mt5`.
- Local sample smoke flow works:
  - API health OK.
  - `/api/mt5/status` returns `waiting_for_bridge` before data.
  - `--post-sample` moves status to `connected`.
  - account, positions, and candles are readable.

Known next gap:

- Python bridge has not yet been verified against a real MT5 terminal/demo account.
- MT5 dashboard has no live refresh.
- Paper signals exist in API but do not yet have a review UI.
- Data quality, stale-data, and operational diagnostics are still thin.

---

## File Structure Plan

### Bridge Files

- `bridges/mt5-python/src/mt5_bridge/__main__.py`
  - CLI flags and run modes.
- `bridges/mt5-python/src/mt5_bridge/config.py`
  - Environment-driven bridge config.
- `bridges/mt5-python/src/mt5_bridge/mt5_adapter.py`
  - MT5 terminal boundary.
- `bridges/mt5-python/src/mt5_bridge/runner.py`
  - Poll loop, single-run, sample posting.
- `bridges/mt5-python/src/mt5_bridge/client.py`
  - Go API client.
- `bridges/mt5-python/tests/`
  - Unit tests for config, payloads, runner, adapter mapping.

### Go API Files

- `services/api/internal/domain/mt5/types.go`
  - MT5 domain contracts.
- `services/api/internal/application/mt5/service.go`
  - XAUUSD-only validation and ingest orchestration.
- `services/api/internal/adapters/postgres/mt5_repository.go`
  - MT5 persistence.
- `services/api/internal/adapters/http/mt5_handler.go`
  - MT5 ingest/read endpoints.
- `services/api/internal/domain/signal/types.go`
  - Paper signal lifecycle.
- `services/api/internal/application/signals/service.go`
  - Paper signal use cases.
- `services/api/internal/adapters/http/signals_handler.go`
  - Paper signal endpoints.

### Web Files

- `apps/web/app/mt5/page.tsx`
  - MT5 dashboard route.
- `apps/web/features/mt5/mt5-dashboard.tsx`
  - Bridge/account/positions UI.
- `apps/web/app/signals/page.tsx`
  - New paper signal review route.
- `apps/web/features/signals/`
  - Signal list, status actions, and create form.
- `apps/web/lib/api-client.ts`
  - MT5 and paper signal API client functions.

### Docs

- `docs/runbooks/mt5-terminal-setup.md`
  - Human setup guide for MT5 demo terminal and bridge.
- `docs/runbooks/mt5-xauusd-next-tasks.md`
  - Lightweight checklist derived from this plan.
- `docs/AGENT_MEMORY.md`
  - Durable state updates after each milestone.

---

## Phase 1: Real MT5 Demo Bridge

### Task 1: Add Bridge Poll Loop

**Files:**

- Modify: `bridges/mt5-python/src/mt5_bridge/__main__.py`
- Modify: `bridges/mt5-python/src/mt5_bridge/runner.py`
- Test: `bridges/mt5-python/tests/test_runner.py`

- [ ] **Step 1: Write failing test for repeated polling**

Add a test that uses a fake adapter and fake client. It should call a new `run_loop(..., iterations=2)` helper and assert that heartbeat/tick/account/positions are posted twice while candles are posted for configured timeframes.

Expected test shape:

```python
def test_run_loop_posts_multiple_iterations(self) -> None:
    adapter = FakeAdapter()
    client = FakeClient()

    run_loop(adapter, BridgeConfig(timeframes=("1m",)), client, iterations=2, sleep_seconds=0)

    heartbeat_posts = [path for path, _payload in client.posts if path == "/api/mt5/heartbeat"]
    self.assertEqual(len(heartbeat_posts), 2)
```

- [ ] **Step 2: Run failing test**

Run:

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m unittest tests/test_runner.py
```

Expected: fail because `run_loop` is not defined.

- [ ] **Step 3: Implement minimal `run_loop`**

Add `run_loop(adapter, config, client, iterations=None, sleep_seconds=None)` in `runner.py`.

Rules:

- `iterations=None` means run forever.
- `sleep_seconds=None` uses `config.poll_seconds`.
- call `run_once` each iteration.
- sleep after each iteration except the final bounded iteration.

- [ ] **Step 4: Add CLI flag**

Modify `__main__.py`:

- keep `--once`
- keep `--dry-run`
- keep `--post-sample`
- add `--loop`
- default error becomes `choose --dry-run, --post-sample, --once, or --loop`

- [ ] **Step 5: Verify**

Run:

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m unittest discover -s tests
```

Expected: all bridge tests pass.

- [ ] **Step 6: Commit**

```bash
git add bridges/mt5-python
git commit -m "feat(mt5): add bridge poll loop"
```

### Task 2: Add Real MT5 Setup Runbook

**Files:**

- Create: `docs/runbooks/mt5-terminal-setup.md`
- Modify: `docs/README.md`
- Modify: `docs/AI_AGENT_INDEX.md`

- [ ] **Step 1: Create setup doc**

Include:

- install MT5 terminal
- open demo account
- enable Algo Trading only if needed later, but clarify current bridge is read-only
- install bridge dependencies
- run Go API
- run `--post-sample`
- run `--once`
- run `--loop`
- troubleshooting for missing `MetaTrader5` package, terminal not initialized, unavailable symbol, no account info

- [ ] **Step 2: Link setup doc**

Add links in:

- `docs/README.md`
- `docs/AI_AGENT_INDEX.md`

- [ ] **Step 3: Verify docs links manually**

Run:

```bash
rg "mt5-terminal-setup" docs
```

Expected: setup doc is linked from docs index files.

- [ ] **Step 4: Commit**

```bash
git add docs
git commit -m "docs: add mt5 terminal setup runbook"
```

### Task 3: Add Bridge Runtime Diagnostics

**Files:**

- Modify: `bridges/mt5-python/src/mt5_bridge/runner.py`
- Modify: `bridges/mt5-python/src/mt5_bridge/__main__.py`
- Test: `bridges/mt5-python/tests/test_runner.py`

- [ ] **Step 1: Write failing test for error heartbeat**

Expected behavior:

- If `run_once` catches adapter/client error in loop mode, bridge posts `/api/mt5/heartbeat` with:
  - `status: "error"`
  - `lastError` containing the exception message

- [ ] **Step 2: Run failing test**

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m unittest tests/test_runner.py
```

Expected: fail because error heartbeat is not sent.

- [ ] **Step 3: Implement error heartbeat**

Add helper in `runner.py`:

```python
def post_error_heartbeat(config: BridgeConfig, client: GoAPIClient, message: str) -> None:
    ...
```

Use `build_heartbeat(status="error", last_error=message)`.

- [ ] **Step 4: Verify**

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m unittest discover -s tests
```

Expected: pass.

- [ ] **Step 5: Commit**

```bash
git add bridges/mt5-python
git commit -m "feat(mt5): report bridge runtime errors"
```

---

## Phase 2: API Data Quality And Operational Readiness

### Task 4: Add MT5 Status Staleness

**Files:**

- Modify: `services/api/internal/adapters/http/mt5_handler.go`
- Modify: `services/api/internal/adapters/http/mt5_handler_test.go`
- Modify: `packages/shared/openapi/api-spec.md`

- [ ] **Step 1: Write failing HTTP test**

Add a test where latest heartbeat is older than 30 seconds. Expected response:

```json
{
  "state": "stale",
  "heartbeat": {
    "status": "healthy"
  }
}
```

- [ ] **Step 2: Run failing test**

```bash
cd services/api
go test ./internal/adapters/http
```

Expected: fail because status remains `connected`.

- [ ] **Step 3: Implement staleness**

In `mt5_handler.go`:

- add `const mt5HeartbeatStaleAfter = 30 * time.Second`
- compute state:
  - `waiting_for_bridge`
  - `connected`
  - `stale`
- use `time.Since(heartbeat.SentAt)`

- [ ] **Step 4: Update API spec**

Document `state: "stale"` and stale threshold.

- [ ] **Step 5: Verify**

```bash
cd services/api
go test ./...
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add services/api packages/shared/openapi/api-spec.md
git commit -m "feat(mt5): add bridge staleness state"
```

### Task 5: Add Latest Tick Endpoint

**Files:**

- Modify: `services/api/internal/adapters/http/mt5_handler.go`
- Modify: `services/api/internal/adapters/http/mt5_handler_test.go`
- Modify: `services/api/internal/adapters/http/router.go`
- Modify: `packages/shared/openapi/api-spec.md`
- Modify: `apps/web/lib/api-client.ts`

- [ ] **Step 1: Write failing HTTP test**

Endpoint:

```http
GET /api/mt5/tick/latest?symbol=XAUUSD
```

Expected: returns `mt5TickResponse`.

- [ ] **Step 2: Run failing test**

```bash
cd services/api
go test ./internal/adapters/http
```

Expected: 404 or missing route.

- [ ] **Step 3: Implement route**

Add handler:

```go
func mt5LatestTickHandler(service mt5Service) http.HandlerFunc
```

Add route in `router.go`.

- [ ] **Step 4: Add client function**

Add in `apps/web/lib/api-client.ts`:

```ts
export async function getLatestMT5Tick(symbol = "XAUUSD"): Promise<MT5TickDto>
```

- [ ] **Step 5: Verify**

```bash
cd services/api && go test ./...
cd ../../apps/web && npm test && npm run typecheck
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add services/api apps/web packages/shared/openapi/api-spec.md
git commit -m "feat(mt5): add latest tick endpoint"
```

### Task 6: Add Data Quality Summary

**Files:**

- Create: `services/api/internal/domain/mt5/quality.go`
- Create: `services/api/internal/application/mt5/quality.go`
- Create: `services/api/internal/application/mt5/quality_test.go`
- Modify: `services/api/internal/adapters/http/mt5_handler.go`
- Modify: `apps/web/features/mt5/mt5-dashboard.tsx`

- [ ] **Step 1: Write failing application test**

Given heartbeat time, tick time, and spread, service should return:

- `healthy` when heartbeat <= 30 seconds and spread > 0
- `stale` when heartbeat > 30 seconds
- `invalid_spread` when ask <= bid

- [ ] **Step 2: Run failing test**

```bash
cd services/api
go test ./internal/application/mt5
```

Expected: fail because quality summary does not exist.

- [ ] **Step 3: Implement minimal quality summary**

Add domain type:

```go
type QualityState string
const (
  QualityHealthy QualityState = "healthy"
  QualityStale QualityState = "stale"
  QualityInvalidSpread QualityState = "invalid_spread"
)
```

- [ ] **Step 4: Add summary to status response**

Extend `/api/mt5/status`:

```json
{
  "quality": {
    "state": "healthy",
    "spread": 0.2,
    "heartbeatAgeSeconds": 3
  }
}
```

- [ ] **Step 5: Show quality panel on dashboard**

Add a compact panel or metric row in `/mt5`.

- [ ] **Step 6: Verify**

```bash
cd services/api && go test ./...
cd ../../apps/web && npm test && npm run typecheck && npm run build
```

Expected: pass.

- [ ] **Step 7: Commit**

```bash
git add services/api apps/web
git commit -m "feat(mt5): add data quality summary"
```

---

## Phase 3: MT5 Dashboard Usability

### Task 7: Add Dashboard Auto Refresh

**Files:**

- Modify: `apps/web/app/mt5/page.tsx`
- Create: `apps/web/features/mt5/mt5-refresh-shell.tsx`
- Modify: `apps/web/features/mt5/mt5-dashboard.tsx`
- Test: `apps/web/tests/smoke.test.mjs`

- [ ] **Step 1: Write smoke test**

Assert:

- route imports `MT5RefreshShell`
- dashboard has refresh interval label or client refresh control

- [ ] **Step 2: Run failing test**

```bash
cd apps/web
npm test
```

Expected: fail because refresh shell does not exist.

- [ ] **Step 3: Implement client refresh shell**

Use a small client component:

```tsx
"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

export function MT5RefreshShell({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  useEffect(() => {
    const id = window.setInterval(() => router.refresh(), 5000);
    return () => window.clearInterval(id);
  }, [router]);
  return <>{children}</>;
}
```

- [ ] **Step 4: Wrap page**

In `apps/web/app/mt5/page.tsx`, wrap `<MT5Dashboard />`.

- [ ] **Step 5: Verify**

```bash
cd apps/web
npm test && npm run typecheck && npm run build
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add apps/web
git commit -m "feat(web): auto refresh mt5 dashboard"
```

### Task 8: Add Paper Signal Review Page

**Files:**

- Create: `apps/web/app/signals/page.tsx`
- Create: `apps/web/features/signals/paper-signals-page.tsx`
- Modify: `apps/web/lib/api-client.ts`
- Modify: `apps/web/tests/smoke.test.mjs`

- [ ] **Step 1: Write smoke test**

Assert:

- route imports `PaperSignalsPage`
- API client contains `/api/paper-signals`
- page contains `Pending Review`, `Approve Paper`, and `Reject`

- [ ] **Step 2: Run failing test**

```bash
cd apps/web
npm test
```

Expected: fail because files do not exist.

- [ ] **Step 3: Add API client types**

Add:

```ts
export type PaperSignalDto = {
  id: string;
  symbol: string;
  timeframe: string;
  side: string;
  status: string;
  confidence: number;
  entryPrice: number;
  stopLoss: number;
  takeProfit: number;
  thesis: string;
  createdAt: string;
  updatedAt: string;
};
```

Add:

- `listPaperSignals()`
- `createPaperSignal(input)`
- `updatePaperSignalStatus(id, status)`

- [ ] **Step 4: Add server page**

Fetch and render paper signals. Show empty state if none.

- [ ] **Step 5: Add status action form**

Use server actions or standard form POST route only if project pattern exists. If not, keep read-only list first and add action buttons in next task.

- [ ] **Step 6: Verify**

```bash
cd apps/web
npm test && npm run typecheck && npm run build
```

Expected: pass.

- [ ] **Step 7: Commit**

```bash
git add apps/web
git commit -m "feat(web): add paper signal review page"
```

### Task 9: Add Create Paper Signal UI

**Files:**

- Modify: `apps/web/features/signals/paper-signals-page.tsx`
- Create: `apps/web/features/signals/create-paper-signal-form.tsx`
- Test: `apps/web/tests/smoke.test.mjs`

- [ ] **Step 1: Write smoke test**

Assert the form includes:

- `XAUUSD`
- `timeframe`
- `side`
- `confidence`
- `entryPrice`
- `stopLoss`
- `takeProfit`
- `thesis`

- [ ] **Step 2: Run failing test**

```bash
cd apps/web
npm test
```

Expected: fail because form is missing.

- [ ] **Step 3: Implement form**

Add a simple form with default symbol fixed to `XAUUSD`.

- [ ] **Step 4: Wire to API**

Use a server action or a minimal client submit handler. Prefer server action if Next.js version supports it cleanly in the existing app.

- [ ] **Step 5: Verify**

```bash
cd apps/web
npm test && npm run typecheck && npm run build
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add apps/web
git commit -m "feat(web): add paper signal creation form"
```

---

## Phase 4: Risk And Review Before Automation

### Task 10: Add Risk Preview For Paper Signals

**Files:**

- Create: `services/api/internal/domain/risk/position_size.go`
- Create: `services/api/internal/domain/risk/position_size_test.go`
- Modify: `services/api/internal/domain/signal/types.go`
- Modify: `services/api/internal/application/signals/service.go`
- Modify: `services/api/migrations/000007_add_signal_risk_fields.sql`

- [ ] **Step 1: Write failing risk test**

Input:

- account equity `10000`
- risk percent `1`
- entry `2325`
- stop `2310`

Expected:

- risk amount `100`
- stop distance `15`
- position size computed deterministically

- [ ] **Step 2: Run failing test**

```bash
cd services/api
go test ./internal/domain/risk
```

Expected: fail because package does not exist.

- [ ] **Step 3: Implement risk calculator**

Keep it pure and broker-agnostic. Do not create execution commands.

- [ ] **Step 4: Store risk preview on paper signal**

Add fields:

- `risk_amount`
- `risk_percent`
- `stop_distance`
- `position_size_estimate`

- [ ] **Step 5: Verify**

```bash
cd services/api
go test ./...
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add services/api
git commit -m "feat(risk): add paper signal risk preview"
```

### Task 11: Add Signal Audit Trail

**Files:**

- Create: `services/api/migrations/000008_create_signal_events.sql`
- Create: `services/api/internal/domain/signal/event.go`
- Modify: `services/api/internal/adapters/postgres/paper_signal_repository.go`
- Modify: `services/api/internal/adapters/http/signals_handler.go`

- [ ] **Step 1: Write failing repository test**

When status changes from `pending_review` to `approved_paper`, repository writes an event:

```text
paper_signal.status_changed
```

- [ ] **Step 2: Run failing test**

```bash
cd services/api
go test ./internal/adapters/postgres
```

Expected: fail because event table/repository code does not exist.

- [ ] **Step 3: Implement migration and repository write**

Signal events are append-only.

- [ ] **Step 4: Add API read endpoint**

```http
GET /api/paper-signals/{id}/events
```

- [ ] **Step 5: Verify**

```bash
cd services/api
go test ./...
```

Expected: pass.

- [ ] **Step 6: Commit**

```bash
git add services/api
git commit -m "feat(signals): add paper signal audit trail"
```

### Task 12: Add Manual Approval Gate Document

**Files:**

- Create: `docs/architecture/auto-trading-readiness.md`
- Modify: `docs/AI_AGENT_INDEX.md`
- Modify: `docs/AGENT_MEMORY.md`

- [ ] **Step 1: Document readiness checklist**

Required before any real order execution:

- demo-only execution adapter
- risk preview
- signal audit trail
- kill switch
- max daily loss
- max open risk
- broker error handling
- explicit user approval
- replayable paper-trade journal

- [ ] **Step 2: Verify links**

```bash
rg "auto-trading-readiness" docs
```

- [ ] **Step 3: Commit**

```bash
git add docs
git commit -m "docs: add auto trading readiness checklist"
```

---

## Phase 5: Paper Trading And Journal

### Task 13: Add Paper Trade Journal Schema

**Files:**

- Create: `services/api/migrations/000009_create_paper_trades.sql`
- Create: `services/api/internal/domain/journal/`
- Create: `services/api/internal/application/journal/`
- Create: `services/api/internal/adapters/http/journal_handler.go`

- [ ] **Step 1: Write failing repository/application tests**

Paper trade should link to:

- `paper_signal_id`
- symbol
- side
- planned entry/SL/TP
- actual entry/exit optional
- outcome
- notes

- [ ] **Step 2: Run failing tests**

```bash
cd services/api
go test ./internal/application/journal ./internal/adapters/postgres
```

- [ ] **Step 3: Implement minimal schema/service/API**

Endpoints:

- `POST /api/paper-trades`
- `GET /api/paper-trades`
- `PATCH /api/paper-trades/{id}`

- [ ] **Step 4: Verify**

```bash
cd services/api
go test ./...
```

- [ ] **Step 5: Commit**

```bash
git add services/api
git commit -m "feat(journal): add paper trade journal"
```

### Task 14: Add Journal UI

**Files:**

- Create: `apps/web/app/journal/page.tsx`
- Create: `apps/web/features/journal/`
- Modify: `apps/web/lib/api-client.ts`
- Modify: `apps/web/tests/smoke.test.mjs`

- [ ] **Step 1: Write smoke test**

Assert:

- route imports journal page component
- page includes `Paper Trades`
- API client includes `/api/paper-trades`

- [ ] **Step 2: Implement journal page**

Show:

- open paper trades
- closed paper trades
- P/L summary
- notes

- [ ] **Step 3: Verify**

```bash
cd apps/web
npm test && npm run typecheck && npm run build
```

- [ ] **Step 4: Commit**

```bash
git add apps/web
git commit -m "feat(web): add paper trade journal"
```

---

## Phase 6: Backtesting And Agent Review

### Task 15: Add Simple XAUUSD Backtest Runner

**Files:**

- Create: `services/api/internal/domain/backtest/`
- Create: `services/api/internal/application/backtests/`
- Create: `services/api/internal/adapters/http/backtests_handler.go`
- Create: `services/api/migrations/000010_create_backtest_runs.sql`

- [ ] **Step 1: Write failing domain test**

Backtest runner consumes historical candles and a simple signal rule. It must not look ahead.

- [ ] **Step 2: Implement deterministic baseline**

Start with a simple EMA cross or close-vs-EMA baseline. Keep ICT/SMC out until the pipeline is trustworthy.

- [ ] **Step 3: Add API endpoints**

- `POST /api/backtests`
- `GET /api/backtests`
- `GET /api/backtests/{id}`

- [ ] **Step 4: Verify**

```bash
cd services/api
go test ./...
```

- [ ] **Step 5: Commit**

```bash
git add services/api
git commit -m "feat(backtests): add xauusd baseline runner"
```

### Task 16: Add Agent Self-Review Record

**Files:**

- Create: `services/api/migrations/000011_create_agent_reviews.sql`
- Create: `services/api/internal/domain/agent/`
- Create: `services/api/internal/application/agent/`
- Create: `services/api/internal/adapters/http/agent_handler.go`

- [ ] **Step 1: Write failing service test**

Agent review record must include:

- signal id
- market context summary
- risk summary
- reasons for approval/rejection
- model/provider metadata optional
- created time

- [ ] **Step 2: Implement storage and API**

Endpoints:

- `POST /api/agent/reviews`
- `GET /api/agent/reviews?signalId=...`

- [ ] **Step 3: Verify**

```bash
cd services/api
go test ./...
```

- [ ] **Step 4: Commit**

```bash
git add services/api
git commit -m "feat(agent): add signal self review records"
```

---

## Final Verification Before Real MT5 Demo Session

Run:

```bash
cd services/api
go test ./...
```

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m unittest discover -s tests
```

```bash
cd apps/web
npm test && npm run typecheck && npm run build
```

Local smoke:

```bash
cd services/api
go run ./cmd/api
```

In another terminal:

```bash
cd bridges/mt5-python
PYTHONPATH=src python3 -m mt5_bridge --post-sample
```

Check:

```bash
curl -sS http://localhost:8080/api/mt5/status
curl -sS 'http://localhost:8080/api/mt5/account/latest?accountLogin=12345678'
curl -sS 'http://localhost:8080/api/mt5/positions/latest?accountLogin=12345678'
```

Expected:

- status state is `connected`
- account returns latest snapshot
- positions returns latest snapshot
- `/mt5` dashboard renders without crashing

---

## Commit Strategy

Commit after every task. Keep commits small and reversible.

Recommended sequence:

1. `feat(mt5): add bridge poll loop`
2. `docs: add mt5 terminal setup runbook`
3. `feat(mt5): report bridge runtime errors`
4. `feat(mt5): add bridge staleness state`
5. `feat(mt5): add latest tick endpoint`
6. `feat(mt5): add data quality summary`
7. `feat(web): auto refresh mt5 dashboard`
8. `feat(web): add paper signal review page`
9. `feat(web): add paper signal creation form`
10. `feat(risk): add paper signal risk preview`
11. `feat(signals): add paper signal audit trail`
12. `docs: add auto trading readiness checklist`
13. `feat(journal): add paper trade journal`
14. `feat(web): add paper trade journal`
15. `feat(backtests): add xauusd baseline runner`
16. `feat(agent): add signal self review records`

---

## Do Not Build Yet

Do not implement:

- MT5 order placement.
- EA WebSocket execution.
- real-money broker integration.
- automatic trade execution.
- portfolio optimization.

Those come after demo-mode paper signal review, audit trail, risk gates, and kill switch are proven.
