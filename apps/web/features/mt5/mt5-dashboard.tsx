import {
  getLatestMT5Account,
  getLatestMT5Positions,
  getMT5Status,
  type MT5AccountSnapshotDto,
  type MT5PositionDto,
  type MT5StatusDto,
} from "@/lib/api-client";

type LoadResult<T> =
  | { data: T; error?: never }
  | { data?: never; error: string };

export async function MT5Dashboard() {
  const statusResult = await loadStatus();
  const accountLogin = statusResult.data?.heartbeat.accountLogin;
  const accountResult = accountLogin ? await loadAccount(accountLogin) : undefined;
  const positionsResult = accountLogin ? await loadPositions(accountLogin) : undefined;

  return (
    <div className="mt5-page">
      <div className="page-header">
        <div>
          <h1>MT5 Bridge</h1>
          <p>XAUUSD read-only ingest status, account state, and position snapshots.</p>
        </div>
      </div>

      <section className="metrics-grid" aria-label="MT5 bridge metrics">
        <div className="metric">
          <span>Bridge</span>
          <strong>{statusResult.data?.heartbeat.status ?? "Offline"}</strong>
        </div>
        <div className="metric">
          <span>Spread</span>
          <strong>{formatSpread(statusResult.data)}</strong>
        </div>
        <div className="metric">
          <span>Equity</span>
          <strong>{formatMoney(accountResult?.data?.equity, accountResult?.data?.currency)}</strong>
        </div>
        <div className="metric">
          <span>Positions</span>
          <strong>{positionsResult?.data?.length ?? 0}</strong>
        </div>
      </section>

      {statusResult.error ? (
        <div className="market-inline-warning" role="status">
          MT5 status unavailable: {statusResult.error}
        </div>
      ) : null}

      <section className="mt5-grid">
        <div className="panel">
          <div className="panel-header">
            <h2>Bridge Health</h2>
            <span>{statusResult.data?.heartbeat.bridgeId ?? "local-mt5"}</span>
          </div>
          {statusResult.data ? <BridgeHealth status={statusResult.data} /> : <EmptyState text="No heartbeat has been received yet." />}
        </div>

        <div className="panel">
          <div className="panel-header">
            <h2>Account Snapshot</h2>
            <span>{accountLogin ?? "waiting"}</span>
          </div>
          {accountResult?.data ? (
            <AccountSnapshot account={accountResult.data} />
          ) : (
            <EmptyState text={accountResult?.error ?? "Waiting for the first account snapshot."} />
          )}
        </div>
      </section>

      <section className="panel">
        <div className="panel-header">
          <h2>Open Positions</h2>
          <span>XAUUSD only</span>
        </div>
        {positionsResult?.data && positionsResult.data.length > 0 ? (
          <PositionsTable positions={positionsResult.data} />
        ) : (
          <EmptyState text={positionsResult?.error ?? "No open XAUUSD positions in the latest snapshot."} />
        )}
      </section>
    </div>
  );
}

function BridgeHealth({ status }: { status: MT5StatusDto }) {
  const spread = status.latestTick.ask - status.latestTick.bid;
  return (
    <div className="kv-list">
      <div>
        <span>Terminal</span>
        <strong>{status.heartbeat.terminal}</strong>
      </div>
      <div>
        <span>Server</span>
        <strong>{status.heartbeat.server}</strong>
      </div>
      <div>
        <span>Last heartbeat</span>
        <strong>{formatDateTime(status.heartbeat.sentAt)}</strong>
      </div>
      <div>
        <span>Latest tick</span>
        <strong>{formatDateTime(status.latestTick.time)}</strong>
      </div>
      <div>
        <span>Bid / Ask</span>
        <strong>
          {status.latestTick.bid.toFixed(2)} / {status.latestTick.ask.toFixed(2)}
        </strong>
      </div>
      <div>
        <span>Spread</span>
        <strong>{spread.toFixed(2)}</strong>
      </div>
    </div>
  );
}

function AccountSnapshot({ account }: { account: MT5AccountSnapshotDto }) {
  return (
    <div className="kv-list">
      <div>
        <span>Balance</span>
        <strong>{formatMoney(account.balance, account.currency)}</strong>
      </div>
      <div>
        <span>Equity</span>
        <strong>{formatMoney(account.equity, account.currency)}</strong>
      </div>
      <div>
        <span>Margin</span>
        <strong>{formatMoney(account.margin, account.currency)}</strong>
      </div>
      <div>
        <span>Free margin</span>
        <strong>{formatMoney(account.freeMargin, account.currency)}</strong>
      </div>
      <div>
        <span>Margin level</span>
        <strong>{account.marginLevel.toFixed(2)}%</strong>
      </div>
      <div>
        <span>Snapshot</span>
        <strong>{formatDateTime(account.time)}</strong>
      </div>
    </div>
  );
}

function PositionsTable({ positions }: { positions: MT5PositionDto[] }) {
  return (
    <table className="table">
      <thead>
        <tr>
          <th>Ticket</th>
          <th>Side</th>
          <th>Volume</th>
          <th>Open</th>
          <th>SL</th>
          <th>TP</th>
          <th>Profit</th>
          <th>Snapshot</th>
        </tr>
      </thead>
      <tbody>
        {positions.map((position) => (
          <tr key={`${position.ticket}-${position.snapshotTime}`}>
            <td>{position.ticket}</td>
            <td>
              <span className={`status ${position.side === "buy" ? "green" : "amber"}`}>{position.side}</span>
            </td>
            <td>{position.volume.toFixed(2)}</td>
            <td>{position.openPrice.toFixed(2)}</td>
            <td>{formatOptionalPrice(position.stopLoss)}</td>
            <td>{formatOptionalPrice(position.takeProfit)}</td>
            <td>{position.profit.toFixed(2)}</td>
            <td>{formatDateTime(position.snapshotTime)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function EmptyState({ text }: { text: string }) {
  return (
    <div className="market-state empty">
      <h3>Waiting for MT5 data</h3>
      <p>{text}</p>
    </div>
  );
}

async function loadStatus(): Promise<LoadResult<MT5StatusDto>> {
  try {
    return { data: await getMT5Status() };
  } catch (error) {
    return { error: error instanceof Error ? error.message : "Unexpected MT5 status error" };
  }
}

async function loadAccount(accountLogin: string): Promise<LoadResult<MT5AccountSnapshotDto>> {
  try {
    return { data: await getLatestMT5Account(accountLogin) };
  } catch (error) {
    return { error: error instanceof Error ? error.message : "Unexpected MT5 account error" };
  }
}

async function loadPositions(accountLogin: string): Promise<LoadResult<MT5PositionDto[]>> {
  try {
    return { data: await getLatestMT5Positions(accountLogin) };
  } catch (error) {
    return { error: error instanceof Error ? error.message : "Unexpected MT5 positions error" };
  }
}

function formatSpread(status?: MT5StatusDto): string {
  if (!status) {
    return "-";
  }
  return (status.latestTick.ask - status.latestTick.bid).toFixed(2);
}

function formatMoney(value?: number, currency = "USD"): string {
  if (value === undefined) {
    return "-";
  }
  return new Intl.NumberFormat("en-US", {
    currency,
    maximumFractionDigits: 2,
    style: "currency",
  }).format(value);
}

function formatDateTime(value: string): string {
  if (!value) {
    return "-";
  }
  return new Intl.DateTimeFormat("en-US", {
    dateStyle: "medium",
    timeStyle: "medium",
  }).format(new Date(value));
}

function formatOptionalPrice(value: number): string {
  if (value === 0) {
    return "-";
  }
  return value.toFixed(2);
}
