# Oh My Trading

Personal AI Trading Agent Dashboard for market analysis, AI signals, risk management, trade journaling, backtesting, and agent monitoring.

## Purpose

Oh My Trading is designed as a personal command center for an AI-assisted trading workflow across XAUUSD, BTC, Forex, and Crypto markets.

The system starts as decision support:

- Analyze markets across multiple timeframes.
- Generate explainable AI trade signals.
- Validate trade ideas against risk rules.
- Journal planned and executed trades.
- Backtest strategy configurations.
- Monitor agent runs, failures, and signal quality.

Auto-trading and broker execution are future features. They should only be added after analysis, risk validation, journaling, backtesting, monitoring, paper trading, audit logs, and a kill switch are in place.

## Architecture

The project uses a clean/hexagonal modular monolith.

```text
apps/web
  Next.js dashboard

services/api
  Go API and worker

packages/shared
  Shared API contracts and schemas

deployments
  Local and production deployment assets

docs
  Product, architecture, runbooks, and agent memory
```

Backend boundaries:

- `services/api/internal/domain`: business rules and domain models.
- `services/api/internal/application`: use-case orchestration.
- `services/api/internal/adapters`: HTTP, database, cache, queue, AI, market data, and notification integrations.
- `services/api/internal/platform`: config, logging, metrics, auth, and runtime utilities.

## Tech Stack

- Backend: Go
- Frontend: Next.js/React
- Database: PostgreSQL with TimescaleDB
- Cache: Redis
- Queue: RabbitMQ
- Charts: TradingView Lightweight Charts
- Deployment: Docker Compose first

## Documentation

Start here:

- [Agent Guide](./AGENTS.md)
- [AI Agent Index](./docs/AI_AGENT_INDEX.md)
- [Agent Memory](./docs/AGENT_MEMORY.md)
- [Product Plan](./docs/product/ai-trading-agent-dashboard-plan.md)
- [System Overview](./docs/architecture/system-overview.md)
- [Database Design](./docs/architecture/database.md)
- [API Spec](./packages/shared/openapi/api-spec.md)
- [Implementation Tasks](./docs/runbooks/implementation-tasks.md)

## Local Setup

The runnable stack is not implemented yet. Follow the implementation tasks in order:

1. Complete repository scaffold.
2. Add Docker Compose for TimescaleDB, Redis, and RabbitMQ.
3. Add the Go API health endpoint.
4. Add the Next.js dashboard shell.
5. Add migrations and market data flows.

Once Task 2 is implemented, local infrastructure should run with:

```bash
docker compose -f deployments/docker-compose.yml up -d
```

You can also run the API in mock mode (no database dependency) for frontend display work:

```bash
cd services/api
OMT_API_MOCK_MODE=true go run ./cmd/api
```

Mock mode serves deterministic sample data for:

- `GET /api/symbols`
- `GET /api/candles`
- `GET /api/indicators`

## Implementation Status

Current phase: planning and scaffold.

Next source of truth: [Implementation Tasks](./docs/runbooks/implementation-tasks.md).

