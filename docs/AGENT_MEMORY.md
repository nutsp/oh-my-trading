# Agent Memory

This file stores durable project memory for AI agents working on Oh My Trading.

Update this file when an important decision, constraint, milestone, or handoff state changes. Keep entries concise and factual.

## Current State

- Project is in planning and scaffold phase.
- Repository structure has been created according to the recommended folder structure.
- Repository structure has been cleaned so project documents live in their intended folders.
- Core documentation has been split into product, architecture, database, API, and implementation task files.
- Root `AGENTS.md` defines project-specific AI agent behavior.
- Root `README.md` has been created as the project entrypoint.
- Go API foundation has been implemented with config loading, structured logging, graceful shutdown, and `GET /api/health`.
- Initial database migration runner and TimescaleDB schema have been implemented.
- Symbol management has been implemented across domain, application service, PostgreSQL repository, and HTTP API.
- Task 1 has been committed as `fd4f633` with message `chore: add project scaffold`.
- Task 2 local infrastructure has been implemented and verified healthy with Docker Compose.

## Product Memory

- The system is a personal AI Trading Agent Dashboard.
- Target markets are XAUUSD, BTC, Forex, and Crypto.
- Trading style is top-down analysis with EMA, RSI, MACD, ATR, and ICT/SMC concepts.
- The product should support market analysis, AI signals, risk management, trade journal, backtesting, and agent monitoring.
- The system should prioritize decision support before automation.
- Auto-trading is a future feature and should not be implemented before risk, journal, backtesting, monitoring, paper trading, audit logs, and a kill switch exist.

## Technical Memory

- Backend preference is Go.
- Frontend preference is Next.js/React.
- Database choice is PostgreSQL with TimescaleDB.
- Cache choice is Redis.
- Queue choice is RabbitMQ for MVP.
- Kafka is reserved for later if event volume or stream replay requirements justify it.
- Chart library is TradingView Lightweight Charts.
- Agent workflow should be custom and LangGraph-style: explicit state, deterministic nodes where possible, AI used mainly for synthesis and review.
- Architecture should be clean/hexagonal and start as a modular monolith.

## Architecture Decisions

- Keep domain logic in `services/api/internal/domain`.
- Keep use-case orchestration in `services/api/internal/application`.
- Keep infrastructure adapters in `services/api/internal/adapters`.
- Keep config, logging, metrics, auth, and clock utilities in `services/api/internal/platform`.
- Keep API contracts under `packages/shared/openapi`.
- Keep product docs under `docs/product`.
- Keep architecture docs under `docs/architecture`.
- Keep implementation runbooks under `docs/runbooks`.

## Safety Memory

- This project may eventually influence financial decisions.
- All trade and signal logic must make risk visible and auditable.
- Server-side risk validation must not be bypassed.
- Backtesting must avoid lookahead bias.
- Broker/API secrets must never be stored as plaintext.
- Real-money execution requires explicit user approval and prerequisite safety controls.

## Documentation Memory

Primary documents:

- `AGENTS.md`
- `docs/AI_AGENT_INDEX.md`
- `docs/product/ai-trading-agent-dashboard-plan.md`
- `docs/architecture/system-overview.md`
- `docs/architecture/database.md`
- `packages/shared/openapi/api-spec.md`
- `docs/runbooks/implementation-tasks.md`

## Implementation Handoff

Next recommended task:

- Commit `Task 5: Symbol Management`, then start `Task 6: Candle Storage And Query API` from `docs/runbooks/implementation-tasks.md`.

Current scaffold exists:

- `apps/web`
- `services/api`
- `packages/shared`
- `deployments`
- `docs`

Before implementing code:

- Confirm whether to initialize Go first, Next.js first, or Docker Compose first.
- If no preference is given, follow `docs/runbooks/implementation-tasks.md` sequentially.

## Decision Log

### 2026-06-18

- Created full implementation plan for AI Trading Agent Dashboard.
- Split plan into focused docs:
  - Product plan
  - System overview
  - Database design
  - API spec
  - Implementation tasks
- Created recommended folder structure.
- Created root `AGENTS.md` for repository-wide agent behavior.
- Created this memory file for durable agent context.
- Cleaned project structure by removing unnecessary `.gitkeep` files from folders that now contain real documentation.
- Completed Task 1 scaffold files and root README content.
- Committed Task 1 as `fd4f633`.
- Added Task 2 Docker Compose and `.env.example`.
- Changed local PostgreSQL host port to `15432` because port `5432` was already allocated.
- Verified Task 2 with `docker compose -f deployments/docker-compose.yml up -d`; TimescaleDB, Redis, and RabbitMQ were healthy.
- Implemented Task 3 Go API foundation; `go test ./...` passed in `services/api`.
- Committed Task 3 as `58ae57d`.
- Implemented Task 4 migration runner and initial schema for TimescaleDB extension, `symbols`, and `candles` hypertable.
- Verified Task 4 with `go test ./...` in `services/api`.
- Committed Task 4 as `e48c044`.
- Implemented Task 5 symbol management.
- Verified Task 5 with `go test ./...` and smoke tested `POST /api/symbols` plus `GET /api/symbols`.
