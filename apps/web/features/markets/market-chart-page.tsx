import Link from "next/link";
import { CandlestickChart } from "@/features/markets/candlestick-chart";
import { type CandleDto, type IndicatorSeriesDto, listCandles, listIndicators } from "@/lib/api-client";

const TIMEFRAMES = ["15m", "1h", "4h", "1d"] as const;
type Timeframe = (typeof TIMEFRAMES)[number];

type MarketChartPageProps = {
  symbol: string;
  timeframe: string;
};

export async function MarketChartPage({ symbol, timeframe }: MarketChartPageProps) {
  const selectedTimeframe = normalizeTimeframe(timeframe);
  const candlesResult = await loadCandles(symbol, selectedTimeframe);
  const indicatorsResult = await loadIndicators(symbol, selectedTimeframe);

  return (
    <div className="market-page">
      <div className="page-header">
        <div>
          <h1>{symbol} Market Chart</h1>
          <p>Live-ready candle view with timeframe controls for manual analysis.</p>
        </div>
      </div>

      <section className="panel">
        <div className="panel-header">
          <h2>{symbol}</h2>
          <div className="timeframe-switcher" role="tablist" aria-label="Timeframe selector">
            {TIMEFRAMES.map((item) => {
              const isActive = item === selectedTimeframe;
              const href = `/markets/${encodeURIComponent(symbol)}?timeframe=${item}`;
              return (
                <Link
                  key={item}
                  href={href}
                  className={`timeframe-item ${isActive ? "active" : ""}`}
                  aria-current={isActive ? "page" : undefined}
                >
                  {item}
                </Link>
              );
            })}
          </div>
        </div>

        <div className="market-chart-shell">
          {candlesResult.error ? (
            <div className="market-state error" role="alert">
              <h3>Unable to load candles</h3>
              <p>{candlesResult.error}</p>
            </div>
          ) : candlesResult.candles.length === 0 ? (
            <div className="market-state empty">
              <h3>No candles found</h3>
              <p>Try another timeframe or run a market sync for this symbol.</p>
            </div>
          ) : (
            <CandlestickChart candles={candlesResult.candles} indicators={indicatorsResult.series} />
          )}
        </div>
      </section>

      <section className="metrics-grid market-indicator-grid" aria-label="Market indicators">
        <div className="metric">
          <span>EMA 20</span>
          <strong>{formatLatest(indicatorsResult.series?.ema20)}</strong>
        </div>
        <div className="metric">
          <span>EMA 50</span>
          <strong>{formatLatest(indicatorsResult.series?.ema50)}</strong>
        </div>
        <div className="metric">
          <span>RSI 14</span>
          <strong>{formatLatest(indicatorsResult.series?.rsi14)}</strong>
        </div>
        <div className="metric">
          <span>ATR 14</span>
          <strong>{formatLatest(indicatorsResult.series?.atr14)}</strong>
        </div>
      </section>
      {indicatorsResult.error ? (
        <p className="market-inline-warning" role="status">
          Indicator overlays unavailable: {indicatorsResult.error}
        </p>
      ) : null}
    </div>
  );
}

function normalizeTimeframe(timeframe: string): Timeframe {
  if (TIMEFRAMES.includes(timeframe as Timeframe)) {
    return timeframe as Timeframe;
  }
  return "1h";
}

async function loadCandles(symbol: string, timeframe: Timeframe): Promise<{
  candles: CandleDto[];
  error?: string;
}> {
  const { from, to } = buildRange(timeframe);
  try {
    const candles = await listCandles({
      symbol,
      timeframe,
      from: from.toISOString(),
      to: to.toISOString(),
    });
    return { candles };
  } catch (error) {
    return {
      candles: [],
      error: error instanceof Error ? error.message : "Unexpected candle request error",
    };
  }
}

async function loadIndicators(symbol: string, timeframe: Timeframe): Promise<{
  series?: IndicatorSeriesDto;
  error?: string;
}> {
  try {
    const response = await listIndicators(symbol, timeframe);
    return { series: response.series };
  } catch (error) {
    return {
      error: error instanceof Error ? error.message : "Unexpected indicator request error",
    };
  }
}

function buildRange(timeframe: Timeframe): { from: Date; to: Date } {
  const to = new Date();
  const hoursByTimeframe: Record<Timeframe, number> = {
    "15m": 72,
    "1h": 24 * 14,
    "4h": 24 * 90,
    "1d": 24 * 365,
  };
  const from = new Date(to.getTime() - hoursByTimeframe[timeframe] * 60 * 60 * 1000);
  return { from, to };
}

function formatLatest(points?: { value: number }[]): string {
  if (!points || points.length === 0) {
    return "-";
  }
  const last = points[points.length - 1];
  return last.value.toFixed(2);
}
