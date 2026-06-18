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

export type MT5HeartbeatDto = {
  bridgeId: string;
  terminal: string;
  accountLogin: string;
  server: string;
  status: string;
  lastError?: string;
  sentAt: string;
};

export type MT5TickDto = {
  symbol: string;
  bid: number;
  ask: number;
  last: number;
  volume: number;
  time: string;
};

export type MT5StatusDto = {
  heartbeat: MT5HeartbeatDto;
  latestTick: MT5TickDto;
};

export type MT5AccountSnapshotDto = {
  accountLogin: string;
  currency: string;
  balance: number;
  equity: number;
  margin: number;
  freeMargin: number;
  marginLevel: number;
  time: string;
};

export type MT5PositionDto = {
  accountLogin: string;
  ticket: string;
  symbol: string;
  side: string;
  volume: number;
  openPrice: number;
  stopLoss: number;
  takeProfit: number;
  profit: number;
  openedAt: string;
  snapshotTime: string;
};

export async function getMT5Status(): Promise<MT5StatusDto> {
  const response = await fetch(`${API_BASE_URL}/api/mt5/status`, {
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error(`Failed to load MT5 status: ${response.status}`);
  }
  return response.json();
}

export async function getLatestMT5Account(accountLogin: string): Promise<MT5AccountSnapshotDto> {
  const query = new URLSearchParams({ accountLogin });
  const response = await fetch(`${API_BASE_URL}/api/mt5/account/latest?${query.toString()}`, {
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error(`Failed to load MT5 account: ${response.status}`);
  }
  return response.json();
}

export async function getLatestMT5Positions(accountLogin: string): Promise<MT5PositionDto[]> {
  const query = new URLSearchParams({ accountLogin });
  const response = await fetch(`${API_BASE_URL}/api/mt5/positions/latest?${query.toString()}`, {
    cache: "no-store",
  });
  if (!response.ok) {
    throw new Error(`Failed to load MT5 positions: ${response.status}`);
  }
  return response.json();
}
