const metrics = [
  ["Open Signals", "3"],
  ["Risk Used", "0.8%"],
  ["Agent Runs", "12"],
  ["Backtests", "4"],
] as const;

const signals = [
  ["XAUUSD", "Long", "0.72", "Watching"],
  ["BTCUSD", "Short", "0.64", "Risk Check"],
  ["EURUSD", "Long", "0.58", "New"],
] as const;

export function DashboardOverview() {
  return (
    <>
      <div className="page-header">
        <div>
          <h1>Trading Agent Dashboard</h1>
          <p>Top-down market context, AI signals, risk state, and agent operations.</p>
        </div>
        <div className="toolbar">
          <button className="button secondary" type="button">
            Sync candles
          </button>
          <button className="button" type="button">
            Run analysis
          </button>
        </div>
      </div>

      <section className="metrics-grid" aria-label="Dashboard metrics">
        {metrics.map(([label, value]) => (
          <div className="metric" key={label}>
            <span>{label}</span>
            <strong>{value}</strong>
          </div>
        ))}
      </section>

      <section className="workspace-grid">
        <div className="panel">
          <div className="panel-header">
            <h2>XAUUSD 1H</h2>
            <span>TradingView chart slot</span>
          </div>
          <div className="chart-placeholder" role="img" aria-label="Candlestick chart placeholder">
            <div className="chart-line" />
          </div>
        </div>

        <div className="panel">
          <div className="panel-header">
            <h2>Latest Signals</h2>
            <span>AI reviewed</span>
          </div>
          <table className="table">
            <thead>
              <tr>
                <th>Symbol</th>
                <th>Bias</th>
                <th>Conf.</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {signals.map(([symbol, bias, confidence, status]) => (
                <tr key={symbol}>
                  <td>{symbol}</td>
                  <td>{bias}</td>
                  <td>{confidence}</td>
                  <td>
                    <span className={`status ${status === "Risk Check" ? "amber" : "green"}`}>
                      {status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </>
  );
}

