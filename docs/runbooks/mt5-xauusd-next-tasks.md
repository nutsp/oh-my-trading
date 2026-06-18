# MT5 XAUUSD Next Tasks

> **For agentic workers:** Detailed execution plan lives at `docs/superpowers/plans/2026-06-18-mt5-xauusd-next-implementation.md`. Use that file when implementing. This file is the high-level checklist.

## Goal

Move from local sample MT5 smoke flow to real MT5 demo-readiness, then add paper signal review, risk preview, journal, backtesting, and agent self-review.

## Phase 1: Real MT5 Demo Bridge

- [ ] Add bridge poll loop.
- [ ] Add MT5 terminal setup runbook.
- [ ] Add bridge runtime diagnostics and error heartbeat.

## Phase 2: API Data Quality And Operational Readiness

- [ ] Add MT5 status staleness.
- [ ] Add latest tick endpoint.
- [ ] Add data quality summary.

## Phase 3: MT5 Dashboard Usability

- [ ] Add dashboard auto refresh.
- [ ] Add paper signal review page.
- [ ] Add create paper signal UI.

## Phase 4: Risk And Review Before Automation

- [ ] Add risk preview for paper signals.
- [ ] Add signal audit trail.
- [ ] Add manual approval gate document.

## Phase 5: Paper Trading And Journal

- [ ] Add paper trade journal schema/API.
- [ ] Add journal UI.

## Phase 6: Backtesting And Agent Review

- [ ] Add simple XAUUSD backtest runner.
- [ ] Add agent self-review records.

## Safety Boundary

Do not build MT5 order execution, real broker integration, or auto-trading until risk gates, audit trail, kill switch, journal, and explicit approval workflows are complete.
