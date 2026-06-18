"use client";

import { useEffect, useRef } from "react";
import {
  CandlestickSeries,
  ColorType,
  LineSeries,
  createChart,
  type IChartApi,
  type UTCTimestamp,
} from "lightweight-charts";
import type { CandleDto, IndicatorSeriesDto } from "@/lib/api-client";

type CandlestickChartProps = {
  candles: CandleDto[];
  indicators?: IndicatorSeriesDto;
};

export function CandlestickChart({ candles, indicators }: CandlestickChartProps) {
  const containerRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!containerRef.current) {
      return;
    }

    const chart = createChart(containerRef.current, {
      autoSize: true,
      layout: {
        background: { type: ColorType.Solid, color: "#ffffff" },
        textColor: "#626a63",
      },
      grid: {
        vertLines: { color: "#eef2ec" },
        horzLines: { color: "#eef2ec" },
      },
      rightPriceScale: {
        borderColor: "#d9ded4",
      },
      timeScale: {
        borderColor: "#d9ded4",
      },
    });

    const candleSeries = chart.addSeries(CandlestickSeries, {
      upColor: "#117a55",
      downColor: "#b33b3b",
      wickUpColor: "#117a55",
      wickDownColor: "#b33b3b",
      borderVisible: false,
    });

    candleSeries.setData(
      candles.map((candle) => ({
        time: Math.floor(new Date(candle.timestamp).getTime() / 1000) as UTCTimestamp,
        open: candle.open,
        high: candle.high,
        low: candle.low,
        close: candle.close,
      })),
    );

    if (indicators) {
      const ema20 = chart.addSeries(LineSeries, {
        color: "#1f6feb",
        lineWidth: 2,
        priceLineVisible: false,
      });
      ema20.setData(
        indicators.ema20.map((point) => ({
          time: Math.floor(new Date(point.timestamp).getTime() / 1000) as UTCTimestamp,
          value: point.value,
        })),
      );

      const ema50 = chart.addSeries(LineSeries, {
        color: "#a46411",
        lineWidth: 2,
        priceLineVisible: false,
      });
      ema50.setData(
        indicators.ema50.map((point) => ({
          time: Math.floor(new Date(point.timestamp).getTime() / 1000) as UTCTimestamp,
          value: point.value,
        })),
      );
    }

    chart.timeScale().fitContent();

    const resizeObserver = new ResizeObserver(() => {
      chart.timeScale().fitContent();
    });
    resizeObserver.observe(containerRef.current);

    return () => {
      resizeObserver.disconnect();
      destroyChart(chart);
    };
  }, [candles]);

  return <div ref={containerRef} className="market-chart-canvas" />;
}

function destroyChart(chart: IChartApi) {
  chart.remove();
}
