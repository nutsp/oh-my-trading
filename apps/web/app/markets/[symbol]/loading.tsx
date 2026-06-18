export default function MarketPageLoading() {
  return (
    <div className="panel market-state loading" role="status" aria-live="polite">
      <h2>Loading market chart</h2>
      <p>Fetching latest candles from the API...</p>
    </div>
  );
}
