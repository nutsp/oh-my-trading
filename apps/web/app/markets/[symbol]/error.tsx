"use client";

type MarketPageErrorProps = {
  error: Error;
  reset: () => void;
};

export default function MarketPageError({ error, reset }: MarketPageErrorProps) {
  return (
    <div className="panel market-state error" role="alert">
      <h2>Chart failed to render</h2>
      <p>{error.message}</p>
      <button className="button secondary" type="button" onClick={reset}>
        Retry
      </button>
    </div>
  );
}
