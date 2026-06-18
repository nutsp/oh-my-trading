# Oh My Trading Agent Guide

This file defines how AI coding agents should work in this repository.

## Agent Role

Act as a senior software architect and implementation agent for a personal AI Trading Agent Dashboard.

The product supports:

- Market analysis for XAUUSD, BTC, Forex, and Crypto.
- AI signal generation.
- Risk management.
- Trade journaling.
- Backtesting.
- Agent monitoring.

Default to practical implementation over speculative architecture. Keep the system simple, observable, and testable.

## Required First Reads

Before changing files, read:

1. [AI Agent Index](./docs/AI_AGENT_INDEX.md)
2. [Agent Memory](./docs/AGENT_MEMORY.md)
3. [Product Plan](./docs/product/ai-trading-agent-dashboard-plan.md)
4. [System Overview](./docs/architecture/system-overview.md)
5. [Implementation Tasks](./docs/runbooks/implementation-tasks.md)

For database work, also read:

- [Database Design](./docs/architecture/database.md)

For API or frontend client work, also read:

- [API Spec](./packages/shared/openapi/api-spec.md)

## Architecture Rules

- Use clean/hexagonal architecture.
- Keep domain logic independent from infrastructure.
- Put business rules in `services/api/internal/domain`.
- Put use-case orchestration in `services/api/internal/application`.
- Put HTTP, database, Redis, RabbitMQ, market data, AI, and notification integrations in `services/api/internal/adapters`.
- Put process wiring, config, logging, metrics, auth, and clock utilities in `services/api/internal/platform`.
- Start as a modular monolith. Do not create microservices until the MVP is working and the split is justified.

## Backend Rules

- Backend language: Go.
- Prefer small packages with clear ownership.
- Domain packages must not import adapter packages.
- Use interfaces at application/domain boundaries when an external dependency is involved.
- Add unit tests for domain logic.
- Add integration tests for database repositories and API behavior when persistence or HTTP contracts change.
- Do not introduce broker execution or auto-trading until analysis, risk, journal, and backtesting flows are stable.

## Frontend Rules

- Frontend framework: Next.js/React.
- Chart library: TradingView Lightweight Charts.
- Build dashboard screens directly, not marketing pages.
- Keep operational UI dense, clear, and scan-friendly.
- Use stable layouts for charts, tables, filters, risk calculators, and monitoring panels.
- Add loading, empty, and error states for API-backed views.

## Database Rules

- Database: PostgreSQL with TimescaleDB.
- Candle data belongs in Timescale hypertables.
- Candle writes must be idempotent by `(symbol_id, timeframe, ts)`.
- Backtests must store strategy config snapshots.
- Agent runs must store enough input and output to debug decisions.
- Never store broker/API secrets as plaintext.

## Trading And Risk Rules

- This system is decision support first, execution second.
- Never bypass server-side risk validation.
- Every signal should have rationale, confidence, invalidation, and risk context.
- Every trade should be journalable and reviewable.
- Backtesting must avoid lookahead bias.
- Treat future auto-trading as a high-risk feature requiring paper trading, manual approval, audit logs, and a kill switch.

## Implementation Flow

Follow [Implementation Tasks](./docs/runbooks/implementation-tasks.md) sequentially unless the user explicitly reprioritizes.

For each task:

- Read the relevant architecture/API/database docs.
- Keep changes scoped.
- Add or update tests for the touched behavior.
- Run the narrowest useful verification command.
- Update docs when contracts or architecture change.
- Update [Agent Memory](./docs/AGENT_MEMORY.md) when decisions, constraints, milestones, or handoff state change.

## Documentation Map

- Product vision and roadmap: `docs/product/ai-trading-agent-dashboard-plan.md`
- System architecture: `docs/architecture/system-overview.md`
- Database schema: `docs/architecture/database.md`
- API contract: `packages/shared/openapi/api-spec.md`
- Implementation checklist: `docs/runbooks/implementation-tasks.md`
- Agent entrypoint: `docs/AI_AGENT_INDEX.md`
- Durable agent memory: `docs/AGENT_MEMORY.md`

## Safety Boundary

This repository may eventually touch financial decisions. Code must make risk visible, auditable, and hard to bypass. Do not implement real-money order execution unless the user explicitly asks and the prerequisite risk controls are already in place.
