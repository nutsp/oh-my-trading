const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export type SymbolDto = {
  id: string;
  code: string;
  market: string;
  baseAsset: string;
  quoteAsset: string;
  enabled: boolean;
};

export async function listSymbols(): Promise<SymbolDto[]> {
  const response = await fetch(`${API_BASE_URL}/api/symbols`, {
    next: { revalidate: 15 },
  });
  if (!response.ok) {
    throw new Error(`Failed to load symbols: ${response.status}`);
  }
  return response.json();
}

export type CandleDto = {
  timestamp: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
};

export type ListCandlesInput = {
  symbol: string;
  timeframe: string;
  from: string;
  to: string;
};

export async function listCandles(input: ListCandlesInput): Promise<CandleDto[]> {
  const query = new URLSearchParams({
    symbol: input.symbol,
    timeframe: input.timeframe,
    from: input.from,
    to: input.to,
  });
  const response = await fetch(`${API_BASE_URL}/api/candles?${query.toString()}`, {
    next: { revalidate: 30 },
  });
  if (!response.ok) {
    throw new Error(`Failed to load candles: ${response.status}`);
  }
  return response.json();
}

export type IndicatorPointDto = {
  timestamp: string;
  value: number;
};

export type MACDPointDto = {
  timestamp: string;
  macd: number;
  signal: number;
  histogram: number;
};

export type IndicatorSeriesDto = {
  ema20: IndicatorPointDto[];
  ema50: IndicatorPointDto[];
  rsi14: IndicatorPointDto[];
  macd: MACDPointDto[];
  atr14: IndicatorPointDto[];
};

export type IndicatorResponseDto = {
  symbol: string;
  timeframe: string;
  series: IndicatorSeriesDto;
};

export async function listIndicators(symbol: string, timeframe: string): Promise<IndicatorResponseDto> {
  const query = new URLSearchParams({
    symbol,
    timeframe,
  });
  const response = await fetch(`${API_BASE_URL}/api/indicators?${query.toString()}`, {
    next: { revalidate: 30 },
  });
  if (!response.ok) {
    throw new Error(`Failed to load indicators: ${response.status}`);
  }
  return response.json();
}

