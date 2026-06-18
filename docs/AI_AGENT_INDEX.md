# AI Agent Index

This index is the documentation entrypoint for AI coding agents. Agents should also read the root [Agent Guide](../AGENTS.md), which defines repository-wide working rules.

## Project

**Name:** Oh My Trading

**Goal:** Build a personal AI Trading Agent Dashboard for market analysis, AI signals, risk management, trade journaling, backtesting, and agent monitoring across XAUUSD, BTC, Forex, and Crypto.

**Primary stack:**

- Backend: Go
- Frontend: Next.js/React
- Database: PostgreSQL with TimescaleDB
- Cache: Redis
- Queue: RabbitMQ
- Charts: TradingView Lightweight Charts
- Architecture: Clean/hexagonal modular monolith

## Read Order

Read these documents in order:

1. [Agent Guide](../AGENTS.md)
2. [Agent Memory](./AGENT_MEMORY.md)
3. [MT5 XAUUSD MVP](./product/mt5-xauusd-mvp.md)
4. [MT5 Integration Architecture](./architecture/mt5-integration.md)
5. [Product Plan](./product/ai-trading-agent-dashboard-plan.md)
6. [System Overview](./architecture/system-overview.md)
7. [Database Design](./architecture/database.md)
8. [API Spec](../packages/shared/openapi/api-spec.md)
9. [MT5 XAUUSD MVP Tasks](./runbooks/mt5-xauusd-mvp-tasks.md)
10. [MT5 XAUUSD Next Tasks](./runbooks/mt5-xauusd-next-tasks.md)
11. [MT5 XAUUSD Next Implementation Plan](./superpowers/plans/2026-06-18-mt5-xauusd-next-implementation.md)
12. [Implementation Tasks](./runbooks/implementation-tasks.md)

## Repository Map

```text
apps/web/
  Next.js dashboard application.

services/api/
  Go backend API and worker code.

packages/shared/schemas/
  Shared schemas used by backend and frontend.

packages/shared/openapi/
  API contract documentation and future OpenAPI files.

deployments/
  Docker Compose, reverse proxy, Grafana, and deployment assets.

docs/architecture/
  System architecture, service boundaries, and database design.

docs/product/
  Product vision, roadmap, MVP scope, and future feature planning.

docs/runbooks/
  Implementation tasks, setup notes, deployment notes, and operations guides.
```

## Work Routing

Use this table to choose the right document before editing.

| Work type | Read first |
| --- | --- |
| Resuming prior work | [Agent Memory](./AGENT_MEMORY.md) |
| Current MVP scope | [MT5 XAUUSD MVP](./product/mt5-xauusd-mvp.md) |
| MT5 bridge or ingest work | [MT5 Integration Architecture](./architecture/mt5-integration.md) |
| Product or roadmap changes | [Product Plan](./product/ai-trading-agent-dashboard-plan.md) |
| Backend architecture changes | [System Overview](./architecture/system-overview.md) |
| Database schema or migrations | [Database Design](./architecture/database.md) |
| API handlers or frontend API client | [API Spec](../packages/shared/openapi/api-spec.md) |
| Completed MVP foundation tasks | [MT5 XAUUSD MVP Tasks](./runbooks/mt5-xauusd-mvp-tasks.md) |
| Current feature implementation order | [MT5 XAUUSD Next Tasks](./runbooks/mt5-xauusd-next-tasks.md) |
| Detailed next implementation plan | [MT5 XAUUSD Next Implementation Plan](./superpowers/plans/2026-06-18-mt5-xauusd-next-implementation.md) |
| Deployment or local infrastructure | [Implementation Tasks](./runbooks/implementation-tasks.md) |

## Current Build Strategy

Start with a modular monolith:

- One Go API service.
- One Go worker service.
- One Next.js web app.
- One PostgreSQL/TimescaleDB database.
- Redis for cache.
- RabbitMQ for async jobs.

Do not start with microservices. Split market ingestion, agent runtime, and execution services only after the MVP is working and operational pressure justifies it.

## Implementation Rule Of Thumb

Current focus is the [MT5 XAUUSD MVP](./product/mt5-xauusd-mvp.md). The MVP foundation tasks are complete; follow [MT5 XAUUSD Next Tasks](./runbooks/mt5-xauusd-next-tasks.md) and the [MT5 XAUUSD Next Implementation Plan](./superpowers/plans/2026-06-18-mt5-xauusd-next-implementation.md) before returning to the broader generic roadmap.

Each task should:

- Keep domain logic independent from adapters.
- Add focused tests for the touched behavior.
- Prefer small commits.
- Avoid broker execution or auto-trading until the analysis, risk, journal, and backtesting flows are stable.
- Update [Agent Memory](./AGENT_MEMORY.md) when the durable project state changes.
