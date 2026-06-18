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
3. [Product Plan](./product/ai-trading-agent-dashboard-plan.md)
4. [System Overview](./architecture/system-overview.md)
5. [Database Design](./architecture/database.md)
6. [API Spec](../packages/shared/openapi/api-spec.md)
7. [Implementation Tasks](./runbooks/implementation-tasks.md)

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
| Product or roadmap changes | [Product Plan](./product/ai-trading-agent-dashboard-plan.md) |
| Backend architecture changes | [System Overview](./architecture/system-overview.md) |
| Database schema or migrations | [Database Design](./architecture/database.md) |
| API handlers or frontend API client | [API Spec](../packages/shared/openapi/api-spec.md) |
| Feature implementation order | [Implementation Tasks](./runbooks/implementation-tasks.md) |
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

Follow the tasks in [Implementation Tasks](./runbooks/implementation-tasks.md) sequentially.

Each task should:

- Keep domain logic independent from adapters.
- Add focused tests for the touched behavior.
- Prefer small commits.
- Avoid broker execution or auto-trading until the analysis, risk, journal, and backtesting flows are stable.
- Update [Agent Memory](./AGENT_MEMORY.md) when the durable project state changes.
